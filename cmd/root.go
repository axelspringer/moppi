// Copyright 2017 Axel Springer SE
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
	"fmt"
	"os"

	"github.com/axelspringer/moppi/server"
	"github.com/spf13/pflag"

	log "github.com/Sirupsen/logrus"
	"github.com/axelspringer/moppi/version"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	verbose bool
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "moppi",
	Short: "",
	// TODO: should be added later
	Long: ``,
	// Run the root command
	Run: server.Run,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// init is run to initialize cobra
func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is fernmelder.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Enables verbose output
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// Some more specific flags
	RootCmd.PersistentFlags().String("marathon", "", "Marathon Endpoint")
	RootCmd.PersistentFlags().String("mesos", "", "Mesos Endpoint")
	RootCmd.PersistentFlags().String("zookeeper", "", "List of Zookeepers")

	// Add all commands
	addCommands(RootCmd)
}

// addCommand is adding and initializing additional commands to the rootCmd
func addCommands(cmd *cobra.Command) {
	// adding version command
	cmd.AddCommand(version.NewCmd())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	var err error

	for _, flags := range []*pflag.FlagSet{RootCmd.PersistentFlags()} {
		err = viper.BindPFlags(flags)
		if err != nil {
			log.WithField("error", err).Fatal("could not bind flags")
		}
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name "config" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err = viper.ReadInConfig()
	if err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
