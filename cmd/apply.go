package cmd

import (
	"fmt"

	"github.com/b1zzu/reportportal-dashboards-as-code/pkg/rpdac"
	"github.com/spf13/cobra"
)

var (
	applyFile    string
	applyProject string

	applyCmd = &cobra.Command{
		Use:   "apply",
		Short: "create or replace ReportPortal object from a YAML definition",
		RunE: func(cmd *cobra.Command, args []string) error {

			file, err := rpdac.LoadFile(applyFile)
			if err != nil {
				return err
			}

			object, err := rpdac.LoadObjectFromFile(file)
			if err != nil {
				return fmt.Errorf("error loading '%s': %w", applyFile, err)
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

				err = r.ApplyDashboard(applyProject, d)
				if err != nil {
					return err
				}

			case rpdac.FilterKind:

				f, err := rpdac.LoadFilterFromFile(file)
				if err != nil {
					return err
				}

				err = r.ApplyFilter(applyProject, f)
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
	applyCmd.Flags().StringVarP(&applyFile, "file", "f", "", "YAML file")
	applyCmd.Flags().StringVarP(&applyProject, "project", "p", "", "ReportPortal Project")

	applyCmd.MarkFlagRequired("file")
	applyCmd.MarkFlagRequired("project")

	rootCmd.AddCommand(applyCmd)
}
