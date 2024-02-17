/*
Copyright © 2022 Seb P sebpsdev@gmail.com
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sebps/template-engine/internal/filtering"
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
		in, _ := cmd.Flags().GetString("in")
		out, _ := cmd.Flags().GetString("out")
		dataPath, _ := cmd.Flags().GetString("data")
		dataFilter, _ := cmd.Flags().GetString("data-filter")
		panicIfNoMatch, _ := cmd.Flags().GetBool("panic-if-no-match")
		leftDelimiter, _ := cmd.Flags().GetString("left-delimiter")
		rightDelimiter, _ := cmd.Flags().GetString("right-delimiter")
		leftLoopVariableDelimiter, _ := cmd.Flags().GetString("left-loop-variable-delimiter")
		rightLoopVariableDelimiter, _ := cmd.Flags().GetString("right-loop-variable-delimiter")
		leftLoopBlockDelimiter, _ := cmd.Flags().GetString("left-loop-block-delimiter")
		rightLoopBlockDelimiter, _ := cmd.Flags().GetString("right-loop-block-delimiter")
		keyColumn, _ := cmd.Flags().GetString("key-column")
		loopVariable, _ := cmd.Flags().GetString("injection-loop-variable")
		multipleOutput, _ := cmd.Flags().GetString("multiple-output")
		multipleOutputFilenamePattern, _ := cmd.Flags().GetString("multiple-output-filename-pattern")

		var isMultipleOutput bool
		if multipleOutput == "true" {
			isMultipleOutput = true
		} else {
			isMultipleOutput = false
		}

		/* rules start */
		inFileInfo, err := os.Stat(in)
		if err != nil {
			panic(err)
		}

		outFileInfo, err := os.Stat(out)
		if err == nil {
			if inFileInfo.IsDir() != outFileInfo.IsDir() {
				panic(errors.New("input and output must be of the same type ( file or directory )"))
			}
		}

		if dataFilter != "" {
			if !filtering.IsJsonPathCompliant(dataFilter) {
				panic(errors.New("wrong data filter format"))
			}
		}
		/* rules end */

		variables, err := parsing.ParseVariablesFile(dataPath, dataFilter, keyColumn, isMultipleOutput, loopVariable)
		if err != nil {
			panic(err)
		}

		if inFileInfo.IsDir() {
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
					variables,
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
		} else {
			template, err := utils.ReadFileContent(in)
			if err != nil {
				panic(err)
			}

			renderAndWrite(
				template,
				variables,
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
	renderCmd.Flags().StringP("data", "d", "", "Data variables path ( json / csv file )")
	renderCmd.Flags().StringP("data-filter", "f", "", "JSONPath filtering expression on data to reduce the input data before rendering")
	renderCmd.Flags().StringP("left-delimiter", "l", "{{", "Left variable delimiter ( default is {{ )")
	renderCmd.Flags().StringP("right-delimiter", "r", "}}", "Right variable delimiter ( default is }} )")
	renderCmd.Flags().StringP("left-loop-variable-delimiter", "", "(", "Left loop variable delimiter ( default is '(' )")
	renderCmd.Flags().StringP("right-loop-variable-delimiter", "", ")", "Right loop variable delimiter ( default is ')' )")
	renderCmd.Flags().StringP("left-loop-block-delimiter", "", "[", "Left loop block delimiter ( default is '[' )")
	renderCmd.Flags().StringP("right-loop-block-delimiter", "", "]", "Right loop block delimiter ( default is ']' )")
	renderCmd.Flags().StringP("panic-if-no-match", "p", "true", "Panic if a variable is not found in the data")
	renderCmd.Flags().StringP("key-column", "k", "id", "Key column ( for .csv variable file ) ( default is 'id' }} )")
	renderCmd.Flags().StringP("injection-loop-variable", "w", "$", "Name of the root loop variable in single file template ( default is '$' }} )")
	renderCmd.Flags().StringP("multiple-output", "", "false", "Whether to generate multiple files from input template and an input data array ( default is 'false' }} )")
	renderCmd.Flags().StringP("multiple-output-filename-pattern", "", "{0}_{i}", "Naming pattern of the generated files in case of multiple-output set to true. Example {0}_{i}_{variable_name} ( default is {0}_{i} }} with {0} : the current file name, {i} : the current file index and {variable_name} : a variable from the data)")
	renderCmd.MarkFlagRequired("in")
	renderCmd.MarkFlagRequired("out")
	renderCmd.MarkFlagRequired("data")
}
