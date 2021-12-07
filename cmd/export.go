package cmd

import (
	"context"
	"os"

	"github.com/b1zzu/reportportal-dashboards-as-code/pkg/reportportal"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
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

			c := oauth2.NewClient(context.TODO(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "TODO"}))

			rp, err := reportportal.NewClient(c, reportportalBaseURL)
			if err != nil {
				return err
			}

			t := &Test{}

			d, _, err := rp.Dashboard.Get(exportProject, exportDashboard)
			if err != nil {
				return err
			}

			t.Dashboard = d

			// retrieve all widgets
			ws := make([]*reportportal.Widget, len(d.Widgets))
			for i, w := range d.Widgets {
				wr, _, err := rp.Widget.Get(exportProject, w.WidgetId)
				if err != nil {
					return err
				}

				ws[i] = wr
			}
			t.Widgets = ws

			b, err := yaml.Marshal(t)
			if err != nil {
				return err
			}

			err = os.WriteFile(exportFile, b, 0660)
			if err != nil {
				return err
			}

			return nil
		},
	}
)

func init() {
	exportCmd.Flags().StringVarP(&exportFile, "file", "f", "", "YAML File")
	exportCmd.Flags().StringVarP(&exportProject, "project", "p", "", "ReportPortal Project")
	exportCmd.Flags().IntVarP(&exportDashboard, "dashboard", "d", -1, "ReportPortal Dashboard ID")

	rootCmd.AddCommand(exportCmd)
}
