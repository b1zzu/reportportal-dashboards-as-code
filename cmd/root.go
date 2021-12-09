package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/b1zzu/reportportal-dashboards-as-code/pkg/reportportal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

const (
	endpointKey = "endpoint"
	tokenKey    = "token"
)

var (
	rootConfigFile string

	rootCmd = &cobra.Command{
		Use:   "rpdac",
		Short: "Import and export ReportPortal dashboards and widget in YAML",
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&rootConfigFile, "config", ".rpdac.toml", "Config file (default: .rpdac.toml)")

	rootCmd.PersistentFlags().StringP(endpointKey, "e", "", "ReportPortal endpoint (example: https://reportportal.example.com)")
	rootCmd.PersistentFlags().StringP(tokenKey, "t", "", "ReportPortal access token")

	viper.BindPFlag("endpoint", rootCmd.PersistentFlags().Lookup(endpointKey))
	viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup(tokenKey))
}

func initConfig() {
	viper.SetConfigFile(rootConfigFile)

	viper.AutomaticEnv()
	viper.SetEnvPrefix("rpdac")

	err := viper.ReadInConfig()
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func requireValue(key string) (string, error) {
	v := viper.GetString(key)
	if v == "" {
		return "", fmt.Errorf("required flag/env/conf \"%s\" is not set", key)
	}
	return v, nil
}

func requireEndpoint() (string, error) {
	return requireValue(endpointKey)
}

func requireToken() (string, error) {
	return requireValue(tokenKey)
}

func requireReportPortalClient() (*reportportal.Client, error) {

	endpoint, err := requireEndpoint()
	if err != nil {
		return nil, err
	}

	token, err := requireToken()
	if err != nil {
		return nil, err
	}

	oc := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))

	// initizlie the ReportPortal client
	rc, err := reportportal.NewClient(oc, endpoint)
	if err != nil {
		return nil, err
	}

	return rc, nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
