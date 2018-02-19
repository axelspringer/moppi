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
	"net"
	"os"
	"path/filepath"
	"sync"

	pb "github.com/axelspringer/moppi/api/v1"
	"github.com/axelspringer/moppi/cfg"
	"github.com/axelspringer/moppi/server"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultListener = "localhost:8080"
	defaultVerbose  = false
)

var (
	exit     = make(chan bool, 1)
	shutdown = make(chan bool, 1)
	wg       sync.WaitGroup
)

var (
	verbose  bool // make to new config
	config   *cfg.Config
	cfgFile  string
	listener string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "moppi",
	Short: "",
	// TODO: should be added later
	Long: ``,
	// Run the root command
	Run: run,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		cfg.Log.Fatal(err)
	}
}

// init is run to initialize cobra
func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./moppi.yml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Enables verbose output
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", defaultVerbose, "verbose output")

	// Bind to
	RootCmd.PersistentFlags().StringVarP(&listener, "listen", "", defaultListener, "Bind listener to (Default: localhost:8080)")

	// Some more specific flags
	RootCmd.PersistentFlags().String("chronos", "", "Chronos endpoints")
	RootCmd.PersistentFlags().String("marathon", "", "Marathon Endpoint")
	RootCmd.PersistentFlags().String("mesos", "", "Mesos Endpoint")
	RootCmd.PersistentFlags().String("zookeeper", "", "List of Zookeepers")

	// Add all commands
	addCommands(RootCmd)
}

// addCommand is adding and initializing additional commands to the rootCmd
func addCommands(cmd *cobra.Command) {
	// adding version command
	cmd.AddCommand(newVersionCmd())

	// adding setup command
	cmd.AddCommand(newSetupCmd())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	var err error

	// map all flags
	for _, flags := range []*pflag.FlagSet{RootCmd.PersistentFlags()} {
		err = viper.BindPFlags(flags)
		if err != nil {
			cfg.Log.WithField("error", err).Fatal("could not bind flags")
		}
	}

	if cfgFile != "" {
		// Use config file from the command line flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.

		// find current directory
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			cfg.Log.WithError(err).Fatal("Could not find current directory")
		}

		// Search config in home directory with name "config" (without extension).
		viper.AddConfigPath(dir)
		viper.AddConfigPath(".")
		viper.SetConfigName("moppi")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err = viper.ReadInConfig()
	if err == nil && verbose {
		cfg.Log.Infof("Using config file:", viper.ConfigFileUsed())
	}

	// Decode command config
	config, err = cfg.New()
	if err = viper.Unmarshal(&config); err != nil {
		cfg.Log.WithError(err).Fatal("Unable to decode into config")
	}

	cfg.Log.SetLevel(log.InfoLevel) // the standard log level is info
	if config.Verbose {             // if verbose extend log level to debug
		cfg.Log.SetLevel(log.DebugLevel)
	}

	// init provider
	config.Init()

	// listenAddr := "localhost:8080"
	grpcServer := grpc.NewServer()
	pb.RegisterUniversesServer(grpcServer, server.NewServer())
	grpcLis, err := net.Listen("tcp", "localhost:")
	if err != nil {
		panic(err)
	}

	go func() {
		fmt.Printf("Starting gRPC server on %s\n", grpcLis.Addr().String())
		err := grpcServer.Serve(grpcLis)
		if err != nil {
			panic(err)
		}
	}()

	// only log, when verbose is enabled
	cfg.Log.Info("Configuration initialized")
}

// run is running a server and is passing along the config
func run(cmd *cobra.Command, args []string) {
	var err error // handle error

	// watch relevant syscalls
	go watchdog()

	// create server
	server, err := server.New(config, exit, &wg)
	if err != nil {
		cfg.Log.Errorf("Failed to create server: %v", err)
		waitGracefulShutdown()
	}

	// start server
	go server.Start()

	// wait for graceful shutdown
	<-shutdown
}
