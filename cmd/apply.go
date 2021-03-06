package cmd

import (
	"github.com/b1zzu/reportportal-dashboards-as-code/pkg/rpdac"
	"github.com/spf13/cobra"
)

var (
	applyFile      string
	applyProject   string
	applyRecursive bool

	applyCmd = &cobra.Command{
		Use:   "apply",
		Short: "create or replace ReportPortal object from a YAML definition",
		RunE: func(cmd *cobra.Command, args []string) error {

			c, err := requireReportPortalClient()
			if err != nil {
				return err
			}
			r := rpdac.NewReportPortal(c)

			return r.Apply(applyProject, applyFile, applyRecursive)
		},
	}
)

func init() {
	applyCmd.Flags().StringVarP(&applyFile, "file", "f", "", "YAML file")
	applyCmd.Flags().StringVarP(&applyProject, "project", "p", "", "ReportPortal Project")
	applyCmd.Flags().BoolVarP(&applyRecursive, "recursive", "r", false, "If file is a directory it will recusive apply all objects in it")

	applyCmd.MarkFlagRequired("file")
	applyCmd.MarkFlagRequired("project")

	rootCmd.AddCommand(applyCmd)
}
