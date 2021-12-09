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
		Short: "create ReportPortal dashboard from a YAML definition",
		RunE: func(cmd *cobra.Command, args []string) error {

			d, err := rpdac.LoadDashboardFromFile(createFile)
			if err != nil {
				return err
			}

			c, err := requireReportPortalClient()
			if err != nil {
				return err
			}

			r := rpdac.NewReportPortal(c)

			err = r.CreateDashboard(createProject, d)
			if err != nil {
				return err
			}

			return nil
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
