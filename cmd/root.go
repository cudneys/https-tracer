/*
Copyright © 2022 Scott Cudney <scott@cudneys.net>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	tracer "github.com/cudneys/http-tracer/tracer"
	"github.com/spf13/cobra"
	"os"
)

var url string
var logTLSInfo bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "http-tracer",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		t, err := tracer.Trace(url, logTLSInfo)
		if err != nil {
			fmt.Printf("HTTP Request Error: %s\n", err.Error())
			return
		}
		fmt.Println(t.AsJSON())
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.http-tracer.yaml)")

	rootCmd.Flags().StringVarP(&url, "url", "u", "", "The URL to test")
	rootCmd.Flags().BoolVarP(&logTLSInfo, "show-tls-info", "s", false, "Display verbose TLS information")
	rootCmd.MarkFlagRequired("url")
}
