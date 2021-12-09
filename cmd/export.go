package cmd

import (
	"log"

	"github.com/b1zzu/reportportal-dashboards-as-code/pkg/reportportal"
	"github.com/b1zzu/reportportal-dashboards-as-code/pkg/rpdac"
	"github.com/spf13/cobra"
)

type Test struct {
	Dashboard *reportportal.Dashboard
	Widgets   []*reportportal.Widget
}

var (
	exportFile      string
	exportProject   string
	exportDashboard int

	exportCmd = &cobra.Command{
		Use:   "export",
		Short: "Exprt a ReportPortal dashboard to YAML",
		RunE: func(cmd *cobra.Command, args []string) error {

			c, err := requireReportPortalClient()
			if err != nil {
				return err
			}

			r := rpdac.NewReportPortal(c)

			// retrieve the Dashboard and Widgets in a single reusable object
			d, err := r.GetDashboard(exportProject, exportDashboard)
			if err != nil {
				return err
			}

			// write the Dashboard object to file in YAML
			err = d.WriteToFile(exportFile)
			if err != nil {
				return err
			}

			log.Printf("Dashboard \"%s\" exported to \"%s\"", d.Name, exportFile)
			return nil
		},
	}
)

func init() {
	exportCmd.Flags().StringVarP(&exportFile, "file", "f", "", "YAML File")
	exportCmd.Flags().StringVarP(&exportProject, "project", "p", "", "ReportPortal Project")
	exportCmd.Flags().IntVarP(&exportDashboard, "dashboard", "d", -1, "ReportPortal Dashboard ID")

	exportCmd.MarkFlagRequired("file")
	exportCmd.MarkFlagRequired("project")
	exportCmd.MarkFlagRequired("dashboard")

	rootCmd.AddCommand(exportCmd)
}
