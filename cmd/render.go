/*
Copyright © 2022 Seb P sebpsdev@gmail.com
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sebps/template-engine/internal/parsing"
	"github.com/sebps/template-engine/internal/rendering"
	"github.com/sebps/template-engine/internal/utils"
	"github.com/spf13/cobra"
)

func renderAndWrite(
	template string,
	variablesSets []map[string]interface{},
	leftDelimiter string,
	rightDelimiter string,
	leftLoopVariableDelimiter string,
	rightLoopVariableDelimiter string,
	leftLoopBlockDelimiter string,
	rightLoopBlockDelimiter string,
	isMultipleOutput bool,
	multipleOutputFilenamePattern string,
	panicIfNoMatch bool,
	path string,
) {
	for i, variables := range variablesSets {
		currentPathOut := path

		if isMultipleOutput {
			currentPathDir := filepath.Dir(path)
			currentPathExtension := filepath.Ext(path)
			currentPathBase := filepath.Base(path)
			currentPathBase = strings.Replace(currentPathBase, currentPathExtension, "", 1)
			currentPathBase = strings.ReplaceAll(multipleOutputFilenamePattern, "{0}", currentPathBase)
			currentPathBase = strings.ReplaceAll(currentPathBase, "{i}", fmt.Sprint(strconv.Itoa(i)))
			currentPathBase, _, _ = rendering.Interpolate(currentPathBase, variables, "{", "}", true)
			currentPathBase = currentPathBase + currentPathExtension
			currentPathOut = filepath.Join(currentPathDir, currentPathBase)
		}

		rendered := rendering.Render(
			template,
			variables,
			leftDelimiter,
			rightDelimiter,
			leftLoopVariableDelimiter,
			rightLoopVariableDelimiter,
			leftLoopBlockDelimiter,
			rightLoopBlockDelimiter,
			panicIfNoMatch,
		)

		err := utils.WriteFileContent(currentPathOut, rendered)
		if err != nil {
			panic(err)
		}
	}
}

var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Render a single file or a full directory",
	Long:  "Render a single file or a full directory",
	Run: func(cmd *cobra.Command, args []string) {
		mode, _ := cmd.Flags().GetString("mode")
		in, _ := cmd.Flags().GetString("in")
		out, _ := cmd.Flags().GetString("out")
		data, _ := cmd.Flags().GetString("data")
		leftDelimiter, _ := cmd.Flags().GetString("left-delimiter")
		rightDelimiter, _ := cmd.Flags().GetString("right-delimiter")
		leftLoopVariableDelimiter, _ := cmd.Flags().GetString("left-loop-variable-delimiter")
		rightLoopVariableDelimiter, _ := cmd.Flags().GetString("right-loop-variable-delimiter")
		leftLoopBlockDelimiter, _ := cmd.Flags().GetString("left-loop-block-delimiter")
		rightLoopBlockDelimiter, _ := cmd.Flags().GetString("right-loop-block-delimiter")
		panicIfNoMatch, _ := cmd.Flags().GetBool("panic-if-no-match")
		keyColumn, _ := cmd.Flags().GetString("key-column")
		loopVariable, _ := cmd.Flags().GetString("wrapping-loop-variable")
		multipleOutput, _ := cmd.Flags().GetString("multiple-output")
		multipleOutputFilenamePattern, _ := cmd.Flags().GetString("multiple-output-filename-pattern")

		var isMultipleOutput bool
		if multipleOutput == "true" {
			isMultipleOutput = true
		} else {
			isMultipleOutput = false
		}

		var variablesSets []map[string]interface{}
		switch isMultipleOutput {
		case true:
			variablesSets = parsing.ParseMultiVariables(data, keyColumn, loopVariable)
		case false:
			variablesSets = make([]map[string]interface{}, 1)
			variablesSets[0] = parsing.ParseSingleVariables(data, keyColumn, loopVariable)
		}

		switch mode {
		case "dir", "DIR":
			filepath.Walk(in, func(pathIn string, info os.FileInfo, err error) error {
				if err != nil {
					panic(err)
				}
				if info.IsDir() {
					return nil
				}

				relativePathIn, _ := filepath.Rel(in, pathIn)
				pathOut := filepath.Join(out, relativePathIn)

				template, err := utils.ReadFileContent(pathIn)
				if err != nil {
					panic(err)
				}

				renderAndWrite(
					template,
					variablesSets,
					leftDelimiter,
					rightDelimiter,
					leftLoopVariableDelimiter,
					rightLoopVariableDelimiter,
					leftLoopBlockDelimiter,
					rightLoopBlockDelimiter,
					isMultipleOutput,
					multipleOutputFilenamePattern,
					panicIfNoMatch,
					pathOut,
				)

				return nil
			})
		case "file", "FILE":
			template, err := utils.ReadFileContent(in)
			if err != nil {
				panic(err)
			}

			renderAndWrite(
				template,
				variablesSets,
				leftDelimiter,
				rightDelimiter,
				leftLoopVariableDelimiter,
				rightLoopVariableDelimiter,
				leftLoopBlockDelimiter,
				rightLoopBlockDelimiter,
				isMultipleOutput,
				multipleOutputFilenamePattern,
				panicIfNoMatch,
				out,
			)
		}
	},
}

func init() {
	rootCmd.AddCommand(renderCmd)

	renderCmd.Flags().StringP("in", "i", "", "Input path ( file or dir )")
	renderCmd.Flags().StringP("out", "o", "", "Output path ( file or dir )")
	renderCmd.Flags().StringP("mode", "m", "file", "Parsing mode ( 'file' or 'dir' ) ( default is 'file' )")
	renderCmd.Flags().StringP("data", "d", "", "Data variables path ( json / csv file )")
	renderCmd.Flags().StringP("left-delimiter", "l", "{{", "Left variable delimiter ( default is {{ )")
	renderCmd.Flags().StringP("right-delimiter", "r", "}}", "Right variable delimiter ( default is }} )")
	renderCmd.Flags().StringP("left-loop-variable-delimiter", "", "(", "Left loop variable delimiter ( default is '(' )")
	renderCmd.Flags().StringP("right-loop-variable-delimiter", "", ")", "Right loop variable delimiter ( default is ')' )")
	renderCmd.Flags().StringP("left-loop-block-delimiter", "", "[", "Left loop block delimiter ( default is '[' )")
	renderCmd.Flags().StringP("right-loop-block-delimiter", "", "]", "Right loop block delimiter ( default is ']' )")
	renderCmd.Flags().StringP("panic-if-no-match", "p", "}}", "Panic if a variable is not found in the data")
	renderCmd.Flags().StringP("key-column", "k", "id", "Key column ( for .csv variable file ) ( default is 'id' }} )")
	renderCmd.Flags().StringP("wrapping-loop-variable", "w", "root", "Name of the root loop variables in template ( for .csv variable file ) ( default is 'loop' }} )")
	renderCmd.Flags().StringP("multiple-output", "", "false", "Whether to generate multiple files from input template and an input data array ( default is 'false' }} )")
	renderCmd.Flags().StringP("multiple-output-filename-pattern", "", "{0}_{i}", "Naming pattern of the generated files in case of multiple-output set to true. Example {0}_{i}_{variable_name} ( default is {0}_{i} }} with {0} : the current file name, {i} : the current file index and {variable_name} : a variable from the data)")
	renderCmd.MarkFlagRequired("in")
	renderCmd.MarkFlagRequired("out")
	renderCmd.MarkFlagRequired("data")
}
