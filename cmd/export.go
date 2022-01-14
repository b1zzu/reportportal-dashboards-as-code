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
	exportFile        string
	exportProject     string
	exportDashboardID int
	exportFilterID    int

	exportCmd = &cobra.Command{
		Use: "export",
		RunE: func(cmd *cobra.Command, args []string) error {

			log.Printf("warning: 'rpdac export' is deprecated, please use 'rpdac export dashboard' instead")
			return exportDashboardCmd.RunE(cmd, args)
		},
	}

	exportDashboardCmd = &cobra.Command{
		Use:   "dashboard",
		Short: "Exprt a ReportPortal dashboard to YAML",
		RunE: func(cmd *cobra.Command, args []string) error {

			c, err := requireReportPortalClient()
			if err != nil {
				return err
			}

			r := rpdac.NewReportPortal(c)

			// retrieve the Dashboard and Widgets in a single reusable object
			d, err := r.GetDashboard(exportProject, exportDashboardID)
			if err != nil {
				return err
			}

			// write the Dashboard object to file in YAML
			err = rpdac.WriteToFile(d, exportFile)
			if err != nil {
				return err
			}

			log.Printf("Dashboard \"%s\" exported to \"%s\"", d.Name, exportFile)
			return nil
		},
	}

	exportFilterCmd = &cobra.Command{
		Use:   "filter",
		Short: "Export a ReportPortal filter to YAML",
		RunE: func(cmd *cobra.Command, args []string) error {

			c, err := requireReportPortalClient()
			if err != nil {
				return err
			}

			r := rpdac.NewReportPortal(c)

			// retrieve the Filter object
			f, err := r.GetFilter(exportProject, exportFilterID)
			if err != nil {
				return err
			}

			// write the Filter object to file in YAML
			err = rpdac.WriteToFile(f, exportFile)
			if err != nil {
				return err
			}

			log.Printf("Filter \"%s\" exported to \"%s\"", f.Name, exportFile)
			return nil
		},
	}
)

func decorateCommonOptions(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&exportFile, "file", "f", "", "YAML File")
	cmd.Flags().StringVarP(&exportProject, "project", "p", "", "ReportPortal Project")

	cmd.MarkFlagRequired("file")
	cmd.MarkFlagRequired("project")
}

func init() {
	// Export CMD
	exportCmd.Flags().StringVarP(&exportFile, "file", "f", "", "(Deprecated) YAML File")
	exportCmd.Flags().StringVarP(&exportProject, "project", "p", "", "(Deprecated) ReportPortal Project")
	exportCmd.Flags().IntVarP(&exportDashboardID, "dashboard", "d", -1, "(Deprecated) ReportPortal Dashboard ID")

	exportCmd.MarkFlagRequired("file")
	exportCmd.MarkFlagRequired("project")
	exportCmd.MarkFlagRequired("dashboard")

	rootCmd.AddCommand(exportCmd)

	// Export Dashboard CMD
	exportDashboardCmd.Flags().IntVar(&exportDashboardID, "id", -1, "ReportPortal Dashboard ID")
	exportDashboardCmd.MarkFlagRequired("id")
	decorateCommonOptions(exportDashboardCmd)

	exportCmd.AddCommand(exportDashboardCmd)

	// Export Filter CMD
	exportFilterCmd.Flags().IntVar(&exportFilterID, "id", -1, "ReportPortal Filter ID")
	exportFilterCmd.MarkFlagRequired("id")
	decorateCommonOptions(exportFilterCmd)

	exportCmd.AddCommand(exportFilterCmd)
}
