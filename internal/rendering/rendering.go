package rendering

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/sebps/template-engine/internal/utils"
)

type Loop struct {
	StartIndex int
	EndIndex   int
	Variable   string
	Values     []map[string]interface{}
	Block      string
	Joiner     string
	Offset     int
}

// Parse the loops blocks into a structure
func ParseLoops(
	structure string,
	variables map[string]interface{},
	leftLoopVariableDelimiter string,
	rightLoopVariableDelimiter string,
	leftLoopBlockDelimiter string,
	rightLoopBlockDelimiter string,
) []*Loop {
	// TODO:
	// 1- Handle the case in which loop.Values is a slice of primitive such as string, int, float, bool
	// 2- Handle infinite recursion
	var loops []*Loop

	// old "(?sm)(?P<loop>(?P<offset>^\\s*)%s(?P<variable>[a-zA-Z0-9\\-\\_]+)%s%s\n*(?P<block>[^%s]*)\n\\s*%s)",
	// test loopsRegexString = "(?sm)(?P<loop>(?P<offset>^\\s*)\\((?P<variable>[a-zA-Z0-9\\-\\_]+)\\)\\[\\[\n*(?P<block>[^\\]\\]]*)\n\\s*\\]\\])"

	variableWrapper := utils.GenerateWrapperRegexp(
		leftLoopVariableDelimiter,
		rightLoopVariableDelimiter,
		"variable",
		false,
	)

	joinerWrapper := utils.GenerateWrapperRegexp(
		leftLoopVariableDelimiter,
		rightLoopVariableDelimiter,
		"joiner",
		false,
	)

	blockWrapper := utils.GenerateWrapperRegexp(
		leftLoopBlockDelimiter,
		rightLoopBlockDelimiter,
		"block",
		true,
	)

	loopsRegexString := fmt.Sprintf(
		"(?sm)(?P<loop>(?P<offset>^\\s*)%s(%s)?%s)",
		variableWrapper,
		joinerWrapper,
		blockWrapper,
	)

	loopsRegexp := regexp.MustCompile(loopsRegexString)
	loopsGroupNames := loopsRegexp.SubexpNames()
	submatchIndexes := loopsRegexp.FindAllStringSubmatchIndex(structure, -1)

	for loopIdx, loopMatch := range loopsRegexp.FindAllStringSubmatch(structure, -1) {
		loop := &Loop{
			Joiner: "\n",
		}

		for loopGroupIdx, loopGroupContent := range loopMatch {
			name := loopsGroupNames[loopGroupIdx]
			if name == "loop" {
				loop.EndIndex = submatchIndexes[loopIdx][2*loopGroupIdx+1]
				loop.StartIndex = submatchIndexes[loopIdx][2*loopGroupIdx]
			} else if name == "offset" {
				loop.Offset = len(loopGroupContent)
			} else if name == "variable" {
				loop.Variable = loopGroupContent
			} else if name == "block" {
				loop.Block = loopGroupContent
			} else if name == "joiner" {
				loop.Joiner = loopGroupContent
			}
		}

		if variable := variables[loop.Variable]; variable != nil {
			loop.Values = make([]map[string]interface{}, 0)
			for i, e := range variable.([]interface{}) {
				eCast := e.(map[string]interface{})
				loop.Values = append(loop.Values, make(map[string]interface{}))
				for k, v := range eCast {
					loop.Values[i][k] = v
				}
			}
		}

		loops = append(loops, loop)
	}

	return loops
}

func CountLeadingWhitespaces(s string) int {
	spaces := 0
	runes := []rune(s)

	for _, r := range runes {
		if unicode.IsSpace(r) {
			spaces++
		} else {
			break
		}
	}

	return spaces
}

// Reindent a string to the required offset
func Reindent(content string, baseOffset int) string {
	var result string
	var contentOffset int
	var deltaOffset int

	contentOffset = CountLeadingWhitespaces(content)
	deltaOffset = contentOffset - baseOffset

	for cursor := 0; cursor < len(content); cursor++ {
		if cursor == 0 && "\n" != string(content[cursor]) {
			for i := 0; i < deltaOffset; i++ {
				cursor++
			}
		}

		result += string(content[cursor])

		if "\n" == string(content[cursor]) {
			for i := 0; i < deltaOffset; i++ {
				cursor++
			}
		}
	}

	return result
}

// Flattify the base structure transforming the loops block into flat content
func FlattifyStructure(structure string, loops []*Loop, leftDelimiter string, rightDelimiter string) string {
	var rendered []byte
	var cursor int

	for _, loop := range loops {
		for cursor < loop.StartIndex {
			rendered = append(rendered, structure[cursor])
			cursor++
		}

		var loopBlocks []string
		var loopRendered string

		for idx, value := range loop.Values {
			// TODO: handle the case in which loop.Values is a slice of primitive such as string, int, float, bool
			var mapping = make(map[string]string)

			for k := range value {
				mapping[k] = loop.Variable + "_" + k + "_" + fmt.Sprint(idx)
			}

			loopBlock := RenameVariables(loop.Block, mapping, leftDelimiter, rightDelimiter)
			indentedLoopBlock := Reindent(loopBlock, loop.Offset)
			loopBlockTrimmed := strings.TrimRight(indentedLoopBlock, "\n\r")
			loopBlockTrimmed = strings.TrimRight(loopBlockTrimmed, "\n")
			loopBlockTrimmed = strings.TrimRight(loopBlockTrimmed, "\\s")
			loopBlocks = append(loopBlocks, loopBlockTrimmed)
		}

		// loopRendered = strings.Join(loopBlocks, loop.Joiner)
		loopRendered = strings.Join(loopBlocks, loop.Joiner+"\n")
		loopTrimmed := strings.TrimRight(loopRendered, "\n\r")
		loopTrimmed = strings.TrimRight(loopTrimmed, "\n")
		loopTrimmed = strings.TrimRight(loopTrimmed, "\\s")

		rendered = append(rendered, []byte(loopTrimmed)...)

		// jump the cursor to the loop ending index character in the structure
		cursor = loop.EndIndex
	}

	for cursor < len(structure) {
		rendered = append(rendered, structure[cursor])
		cursor++
	}

	return string(rendered)
}

// Produce the flat variables map corresponding to the flattified structure
func FlattifyVariables(variables map[string]interface{}, loops []*Loop) map[string]interface{} {
	var flatVariables = make(map[string]interface{})

	for _, loop := range loops {
		for idx, value := range loop.Values {
			// TODO: handle the case in which loop.Values is a slice of primitive such as string, int, float, bool
			for k, v := range value {
				flatK := loop.Variable + "_" + k + "_" + fmt.Sprint(idx)
				flatVariables[flatK] = v
			}
		}
	}

	for k, v := range variables {
		flatVariables[k] = v
	}

	return flatVariables
}

// Rename the variable of a structure with a mapping of names
func RenameVariables(
	structure string,
	mapping map[string]string,
	leftDelimiter string,
	rightDelimiter string,
) string {
	var rendered = structure

	for k, v := range mapping {
		rendered = strings.ReplaceAll(
			rendered,
			leftDelimiter+k+rightDelimiter,
			leftDelimiter+v+rightDelimiter,
		)
	}

	return rendered
}

// Interpolate a structure with a map of variables
func Interpolate(
	structure string,
	variables map[string]interface{},
	leftDelimiter string,
	rightDelimiter string,
	panicIfNoMatch bool,
) string {
	var rendered = structure

	for k, v := range variables {
		rendered = strings.ReplaceAll(rendered, leftDelimiter+k+rightDelimiter, fmt.Sprintf("%v", v))
	}

	if panicIfNoMatch {
		variableWrapperRegexString := utils.GenerateWrapperRegexp(
			leftDelimiter,
			rightDelimiter,
			"variable",
			false,
		)
		variableRegexp := regexp.MustCompile(variableWrapperRegexString)
		variableGroupNames := variableRegexp.SubexpNames()

		for _, variableMatch := range variableRegexp.FindAllStringSubmatch(rendered, -1) {
			for variableGroupIdx, variableGroupContent := range variableMatch {
				name := variableGroupNames[variableGroupIdx]
				if name == "variable" {
					panic(fmt.Sprintf("variable : %q not found in data", variableGroupContent))
				}
			}
		}
	}

	return rendered
}

func Render(
	template string,
	variables map[string]interface{},
	leftDelimiter string,
	rightDelimiter string,
	leftLoopVariableDelimiter string,
	rightLoopVariableDelimiter string,
	leftLoopBlockDelimiter string,
	rightLoopBlockDelimiter string,
	panicIfNoMatch bool,
) string {
	if len(leftDelimiter) == 0 {
		leftDelimiter = "{{"
	}
	if len(rightDelimiter) == 0 {
		rightDelimiter = "}}"
	}
	if len(leftLoopVariableDelimiter) == 0 {
		leftDelimiter = "("
	}
	if len(rightLoopVariableDelimiter) == 0 {
		rightDelimiter = ")"
	}
	if len(leftLoopBlockDelimiter) == 0 {
		leftDelimiter = "["
	}
	if len(rightLoopBlockDelimiter) == 0 {
		rightDelimiter = "]"
	}

	var loops []*Loop
	var flatStructure string
	var flatVariables map[string]interface{}
	var rendered string

	loops = ParseLoops(
		template,
		variables,
		leftLoopVariableDelimiter,
		rightLoopVariableDelimiter,
		leftLoopBlockDelimiter,
		rightLoopBlockDelimiter,
	)

	flatStructure = FlattifyStructure(
		template,
		loops,
		leftDelimiter,
		rightDelimiter,
	)

	flatVariables = FlattifyVariables(
		variables,
		loops,
	)

	rendered = Interpolate(
		flatStructure,
		flatVariables,
		leftDelimiter,
		rightDelimiter,
		panicIfNoMatch,
	)

	return rendered
}
