/*
Copyright © 2022 Seb P sebpsdev@gmail.com

*/
package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/sebps/template-engine/rendering"
	"github.com/spf13/cobra"
)

func readFileContent(filePath string) string {
	fileRawContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return string(fileRawContent)
}

func renderAndWrite(template string, variables map[string]interface{}, path string) {
	// render result
	rendered := rendering.Render(template, variables)

	// prepare dir
	if _, err := os.Stat(filepath.Dir(path)); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(path), 0777)
	}

	// write result
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.WriteString(rendered)
	if err != nil {
		panic(err)
	}
}

func readVariables(filePath string) map[string]interface{} {
	// parse template variables
	dataRawContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	var variables map[string]interface{}
	err = json.Unmarshal(dataRawContent, &variables)
	if err != nil {
		log.Fatal(err)
	}

	return variables
}

// renderCmd represents the render command
var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// render flags
		mode, _ := cmd.Flags().GetString("mode")
		in, _ := cmd.Flags().GetString("in")
		out, _ := cmd.Flags().GetString("out")
		data, _ := cmd.Flags().GetString("data")
		variables := readVariables(data)

		switch mode {
		case "dir", "DIR":
			filepath.Walk(in, func(pathIn string, info os.FileInfo, err error) error {
				if err != nil {
					log.Fatalf(err.Error())
				}

				if info.IsDir() {
					return nil
				}
				pathOut := strings.Replace(pathIn, in, out, 1)
				template := readFileContent(pathIn)
				renderAndWrite(template, variables, pathOut)

				return nil
			})
		case "file", "FILE":
			template := readFileContent(in)
			renderAndWrite(template, variables, out)
		}
	},
}

func init() {
	rootCmd.AddCommand(renderCmd)

	// serveCmd.Flags().StringP("address", "a", "", "Template engine server address")
	renderCmd.Flags().StringP("in", "i", "", "Input path ( file or dir )")
	renderCmd.Flags().StringP("out", "o", "", "Output path ( file or dir )")
	renderCmd.Flags().StringP("mode", "m", "file", "Parsing mode ( 'file' or 'dir' ) ( default is 'file' )")
	renderCmd.Flags().StringP("data", "d", "", "Data variables path ( json file )")
	renderCmd.MarkFlagRequired("in")
	renderCmd.MarkFlagRequired("out")
	renderCmd.MarkFlagRequired("data")
}
