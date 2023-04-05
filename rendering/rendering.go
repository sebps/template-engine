package rendering

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

type Loop struct {
	StartIndex int
	EndIndex   int
	Variable   string
	Values     []map[string]interface{}
	Block      string
	Offset     int
}

// Parse the loops blocks into a structure
func ParseLoops(structure string, args map[string]interface{}) []*Loop {
	// TODO:
	// 1- Handle the case in which loop.Values is a slice of primitive such as string, int, float, bool
	// 2- Handle infinite recursion
	var loops []*Loop
	loopsRegexp := regexp.MustCompile("(?sm)(?P<loop>(?P<offset>^\\s*)\\((?P<variable>[a-zA-Z0-9\\-\\_]+)\\)\\[\n*(?P<block>[^\\]]*)\n\\s*\\])")
	loopsGroupNames := loopsRegexp.SubexpNames()
	submatchIndexes := loopsRegexp.FindAllStringSubmatchIndex(structure, -1)

	for loopIdx, loopMatch := range loopsRegexp.FindAllStringSubmatch(structure, -1) {
		loop := &Loop{}

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
			}
		}

		if arg := args[loop.Variable]; arg != nil {
			loop.Values = make([]map[string]interface{}, 0)
			for i, e := range arg.([]interface{}) {
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

			for k, _ := range value {
				mapping[k] = loop.Variable + "_" + k + "_" + fmt.Sprint(idx)
			}

			loopBlock := RenameVariables(loop.Block, mapping, leftDelimiter, rightDelimiter)
			indentedLoopBlock := Reindent(loopBlock, loop.Offset)
			loopBlocks = append(loopBlocks, indentedLoopBlock)
		}

		loopRendered = strings.Join(loopBlocks, "\n")
		rendered = append(rendered, []byte(loopRendered)...)

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
func RenameVariables(structure string, mapping map[string]string, leftDelimiter string, rightDelimiter string) string {
	var rendered = structure

	for k, v := range mapping {
		rendered = strings.ReplaceAll(rendered, leftDelimiter+k+rightDelimiter, leftDelimiter+v+rightDelimiter)
	}

	return rendered
}

// Interpolate a structure with a map of variables
func Interpolate(structure string, variables map[string]interface{}, leftDelimiter string, rightDelimiter string) string {
	var rendered = structure

	for k, v := range variables {
		rendered = strings.ReplaceAll(rendered, leftDelimiter+k+rightDelimiter, fmt.Sprintf("%v", v))
	}

	return rendered
}

func Render(template string, variables map[string]interface{}, leftDelimiter string, rightDelimiter string) string {
	if len(leftDelimiter) == 0 {
		leftDelimiter = "{{"
	}
	if len(rightDelimiter) == 0 {
		rightDelimiter = "}}"
	}

	var loops []*Loop
	var flatStructure string
	var flatVariables map[string]interface{}
	var rendered string

	loops = ParseLoops(template, variables)
	flatStructure = FlattifyStructure(template, loops, leftDelimiter, rightDelimiter)
	flatVariables = FlattifyVariables(variables, loops)
	rendered = Interpolate(flatStructure, flatVariables, leftDelimiter, rightDelimiter)

	return rendered
}
