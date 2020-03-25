/*
Copyright © 2020 Julien SENON

Credit @aerostitch

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

//Package cmd handle command line server
package cmd

import (
	"sync"

	"github.com/jsenon/aws-cleanup/configs"
	"github.com/jsenon/aws-cleanup/internal/cloudwatch"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	mylog "github.com/jsenon/aws-cleanup/internal/log"
)

var logsgroupname string
var dayslg int
var daysls int

const deltaN = 2

// rootCmd represents the base command when called without any subcommands
var cloudwatchCmd = &cobra.Command{
	Use:   "cloudwatch",
	Short: "cloudwatch cleaning",
	Long: `Perform cleaning of loggroup and logstream.
	Default values are:
	 - LogStrean cleaning if last event timestamp is over 30 days
	 - LogGroup cleaning if no logStream and are older than 90 days`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Logger = log.With().Str("Service", configs.Service).Logger()
		log.Logger = log.With().Str("Version", configs.Version).Logger()

		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		if loglevel {
			err := mylog.SetDebug()
			if err != nil {
				log.Error().Msgf("Could not set loglevel to debug: %v", err)
			}
			log.Debug().Msg("Log level set to Debug")
		}
		cloudwatchLaunch()
	},
}

func init() {
	rootCmd.AddCommand(cloudwatchCmd)
	cloudwatchCmd.PersistentFlags().StringVar(&logsgroupname, "logsgroupname", "all",
		"Logs group name to delete")

	err := viper.BindPFlag("logsgroupname", cloudwatchCmd.PersistentFlags().Lookup("logsgroupname"))
	if err != nil {
		log.Error().Msgf("Error binding logsgroupname: %v", err.Error())
	}

	cloudwatchCmd.PersistentFlags().IntVar(&dayslg, "dayslg", 90,
		"Cleaning loggroup after this number of days")

	err = viper.BindPFlag("daysloggroup", cloudwatchCmd.PersistentFlags().Lookup("dayslg"))
	if err != nil {
		log.Error().Msgf("Error binding dayslg: %v", err.Error())
	}

	cloudwatchCmd.PersistentFlags().IntVar(&daysls, "daysls", 30,
		"Cleaning logstream after this number of days")

	err = viper.BindPFlag("dayslogstream", cloudwatchCmd.PersistentFlags().Lookup("daysls"))
	if err != nil {
		log.Error().Msgf("Error binding daysls: %v", err.Error())
	}

	cobra.OnInitialize(initConfig)
}

func cloudwatchLaunch() {
	log.Info().Msgf("Debug Mode: %t", viper.GetBool("LOGLEVEL"))

	logSName := viper.GetString("LOGSGROUPNAME")
	daysLg := viper.GetInt("DAYSLOGGROUP")
	daysLs := viper.GetInt("DAYSLOGSTREAM")

	if logSName == "all" {
		log.Debug().Msgf("LogGroup deletion after: %v, LogStream deletion after: %v, LogsGroupName deletion: %v",
			daysLg, daysLs, logSName)
	} else {
		log.Debug().Msgf("LogStream deletion after: %v, LogsGroupName deletion: %v",
			daysLs, logSName)
	}

	// Launch cloudwatch cleaning
	// if variable has been set
	p := cloudwatch.NewCwProcessor(daysLs, daysLg, logSName)
	wg := sync.WaitGroup{}
	wg.Add(deltaN)

	go func() {
		p.CleanupLogStreams()
		wg.Done()
	}()

	go func() {
		p.CleanupLogGroups()
		wg.Done()
	}()

	p.ProcessLogGroups()
	wg.Wait()
}
