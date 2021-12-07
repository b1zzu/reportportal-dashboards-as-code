package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	reportportalBaseURL = "TODO"
)

var rootCmd = &cobra.Command{
	Use:   "rpdac",
	Short: "Import and export ReportPortal dashboards and widget in YAML",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
