package cmd

import (
	"github.com/spf13/cobra"
	"user-service/internal/app"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start",
	Long:  `Start`,
	Run: func(cmd *cobra.Command, args []string) {
		app.RunServer()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
