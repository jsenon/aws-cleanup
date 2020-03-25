/*
Copyright © 2020 Julien SENON

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

//Package cmd handle command line root
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/jsenon/aws-cleanup/configs"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var cfgFile string
var loglevel bool
var version bool

const exitOne = 1

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "aws-cleanup",
	Short: "A cleaner tool for AWS",
	Long:  `Cleaning AWS Object like cloudwatch log stream`,
	Run: func(cmd *cobra.Command, args []string) {
		if version {
			fmt.Printf("Version: %v, build from: %v, on: %v\n", configs.Version, configs.GitCommit, configs.BuildDate)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(exitOne)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.process-launcher.yaml)")
	rootCmd.PersistentFlags().BoolVar(&loglevel, "debug", false, "Set log level to Debug")
	rootCmd.PersistentFlags().BoolVar(&version, "version", false, "version")

	err := viper.BindPFlag("loglevel", rootCmd.PersistentFlags().Lookup("debug"))
	if err != nil {
		log.Error().Msgf("Error binding loglevel value: %v", err.Error())
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(exitOne)
		}

		// Search config in home directory with name ".aws-cleanup" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".aws-cleanup")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
