// Copyright © 2019 Red Hat
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"path"

	"github.com/fusor/cpma/env"
	"github.com/fusor/cpma/internal/config"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var debugLogLevel bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cpma",
	Short: "Helps migration cluster configuration of a OCP 3.x cluster to OCP 4.x",
	Long:  `Helps migration cluster configuration of a OCP 3.x cluster to OCP 4.x`,
	Run: func(cmd *cobra.Command, args []string) {
		// Set loglevel to debug
		if debugLogLevel {
			log.SetLevel(log.DebugLevel)
		}

		err := config.InitConfig()
		if err != nil {
			log.Fatal(err, "\n")
		}

		envConfig := env.New()

		if envConfig.LocalOnly {
			envConfig.LoadSrc()
		} else {
			envConfig.FetchSrc()
		}

		envConfig.Parse()

		log.Print(envConfig.Show())
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize()
	rootCmd.PersistentFlags().StringVar(&config.ConfigFile, "config", "", "config file (default is $HOME/.cpma.yaml)")
	rootCmd.PersistentFlags().BoolVar(&debugLogLevel, "debug", false, "set log level to debug")

	rootCmd.Flags().StringP("output-dir", "o", "", "set the directory to store extracted configuration.")
	config.Config().BindPFlag("outputPath", rootCmd.Flags().Lookup("output-dir"))
	config.Config().SetDefault("outputPath", path.Dir(""))

	rootCmd.Flags().BoolP("local-only", "l", false, "do not fetch files, use only local files.")
	config.Config().BindPFlag("localOnly", rootCmd.Flags().Lookup("local-only"))
	config.Config().SetDefault("localOnly", false)
}
