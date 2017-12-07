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

package cfg

import (
	"net"

	log "github.com/sirupsen/logrus"
)

// New returns a new config
// cmdCfg: is a CmdConfig that is used to create a config
func New(cmdCfg *CmdConfig) (*Config, error) {
	return mustNew(cmdCfg)
}

// mustNew wraps the creation of a new config
func mustNew(cmdCfg *CmdConfig) (*Config, error) {
	// config
	var cfg Config

	// check for logger
	cfg.Logger = cmdCfg.Logger
	if cfg.Logger == nil {
		cfg.Logger = log.New()
		// this is the standard setting
		cfg.Logger.SetLevel(log.WarnLevel)
		if cmdCfg.Verbose {
			cfg.Logger.SetLevel(log.DebugLevel)
		}
	}

	// create listener
	listener, err := net.Listen("tcp", cmdCfg.Listen)
	if err != nil {
		return nil, err
	}
	cfg.Listener = listener

	// map additional values
	cfg.Chronos = cmdCfg.Chronos
	cfg.Listen = cmdCfg.Listen
	cfg.Listener = listener
	cfg.Marathon = cmdCfg.Marathon
	cfg.Mesos = cmdCfg.Mesos
	cfg.Verbose = cmdCfg.Verbose
	cfg.Zookeeper = cmdCfg.Zookeeper
	cfg.Etcd = cmdCfg.Etcd

	// initialize provider
	store, err := cfg.Etcd.CreateStore(DefaultBucket)
	if err != nil {
		cfg.Logger.Fatalf("Initializing provider: %v", err)
	}
	cfg.Etcd.SetKVClient(store)

	return &cfg, nil
}
