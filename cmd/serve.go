/*
Copyright © 2022 Seb P sebpsdev@gmail.com
*/
package cmd

import (
	"github.com/sebps/template-engine/server"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		address, _ := cmd.Flags().GetString("address")
		port, _ := cmd.Flags().GetInt("port")
		leftDelimiter, _ := cmd.Flags().GetString("leftDelimiter")
		rightDelimiter, _ := cmd.Flags().GetString("rightDelimiter")
		leftLoopVariableDelimiter, _ := cmd.Flags().GetString("leftLoopVariableDelimiter")
		rightLoopVariableDelimiter, _ := cmd.Flags().GetString("rightLoopVariableDelimiter")
		leftLoopBlockDelimiter, _ := cmd.Flags().GetString("leftLoopBlockDelimiter")
		rightLoopBlockDelimiter, _ := cmd.Flags().GetString("rightLoopBlockDelimiter")

		server.Serve(
			address,
			port,
			leftDelimiter,
			rightDelimiter,
			leftLoopVariableDelimiter,
			rightLoopVariableDelimiter,
			leftLoopBlockDelimiter,
			rightLoopBlockDelimiter,
		)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringP("address", "a", "0.0.0.0", "Template engine server address ( default is 0.0.0.0 )")
	serveCmd.Flags().IntP("port", "p", -1, "Template engine server port")
	serveCmd.MarkFlagRequired("port")
	serveCmd.Flags().StringP("leftDelimiter", "l", "{{", "Left variable delimiter ( default is {{ )")
	serveCmd.Flags().StringP("rightDelimiter", "r", "}}", "Right variable delimiter ( default is }} )")
}
