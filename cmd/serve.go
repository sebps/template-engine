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
		// address, _ := cmd.Flags().GetString("address")
		port, _ := cmd.Flags().GetInt("port")
		server.Serve(port)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// serveCmd.Flags().StringP("address", "a", "", "Template engine server address")
	serveCmd.Flags().IntP("port", "p", -1, "Template engine server port")
	serveCmd.MarkFlagRequired("address")
	serveCmd.MarkFlagRequired("port")
}
