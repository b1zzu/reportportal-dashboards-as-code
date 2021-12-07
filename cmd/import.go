package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	importFile string

	importCmd = &cobra.Command{
		Use:   "import",
		Short: "Import a YAML dashboard to ReportPortal",
		RunE: func(cmd *cobra.Command, args []string) error {

			byteFile, err := ioutil.ReadFile(importFile)
			if err != nil {
				return err
			}

			data := make(map[interface{}]interface{})

			err = yaml.Unmarshal(byteFile, data)
			if err != nil {
				return err
			}

			fmt.Println(data)

			return nil
		},
	}
)

func init() {
	importCmd.Flags().StringVarP(&importFile, "file", "f", "", "YAML file")

	rootCmd.AddCommand(importCmd)
}
