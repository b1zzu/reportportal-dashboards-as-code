package cmd

import (
	"github.com/b1zzu/reportportal-dashboards-as-code/pkg/rpdac"
	"github.com/spf13/cobra"
)

var (
	createFile    string
	createProject string

	createCmd = &cobra.Command{
		Use:   "create",
		Short: "create ReportPortal object from a YAML definition",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := requireReportPortalClient()
			if err != nil {
				return err
			}
			r := rpdac.NewReportPortal(c)

			return r.Create(createProject, createFile)
		},
	}
)

func init() {
	createCmd.Flags().StringVarP(&createFile, "file", "f", "", "YAML file")
	createCmd.Flags().StringVarP(&createProject, "project", "p", "", "ReportPortal Project")

	createCmd.MarkFlagRequired("file")
	createCmd.MarkFlagRequired("project")

	rootCmd.AddCommand(createCmd)
}
