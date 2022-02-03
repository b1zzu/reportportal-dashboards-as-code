package cmd

import (
	"fmt"

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

			file, err := rpdac.LoadFile(createFile)
			if err != nil {
				return err
			}

			object, err := rpdac.LoadObjectFromFile(file)
			if err != nil {
				return fmt.Errorf("error loading '%s': %w", createFile, err)
			}

			c, err := requireReportPortalClient()
			if err != nil {
				return err
			}

			r := rpdac.NewReportPortal(c)

			switch object.Kind {
			case rpdac.DashboardKind:

				d, err := rpdac.LoadDashboardFromFile(file)
				if err != nil {
					return err
				}

				err = r.Dashboard.CreateDashboard(createProject, d)
				if err != nil {
					return err
				}

			case rpdac.FilterKind:

				f, err := rpdac.LoadFilterFromFile(file)
				if err != nil {
					return err
				}

				err = r.Filter.CreateFilter(createProject, f)
				if err != nil {
					return err
				}

			default:
				return fmt.Errorf("unknown Kind '%s' in file '%s'", object.Kind, createFile)

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
