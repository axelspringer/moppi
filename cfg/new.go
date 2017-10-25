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

	"github.com/axelspringer/moppi/provider/etcd"
	log "github.com/sirupsen/logrus"
)

// New returns a new config
// cmdCfg: is a CmdConfig that is used to create a config
func New(cmdCfg *CmdConfig) (*Config, error) {
	return mustNew(cmdCfg)
}

// mustNew wraps the creation of a new config
func mustNew(cmdCfg *CmdConfig) (*Config, error) {
	var logger = cmdCfg.Logger
	if logger == nil {
		logger = log.New()
		// this is the standard setting
		logger.SetLevel(log.WarnLevel)
		if cmdCfg.Verbose {
			logger.SetLevel(log.DebugLevel)
		}
	}

	// default Etcd
	var defaultEtcd etcd.Provider
	// subject to change
	defaultEtcd.Prefix = cmdCfg.Etcd.Prefix
	defaultEtcd.Endpoint = cmdCfg.Etcd.Endpoint

	// providers
	providers := &Providers{
		Etcd: defaultEtcd,
	}

	// create listener
	listener, err := net.Listen("tcp", cmdCfg.Listen)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		CmdConfig: cmdCfg,
		Listener:  listener,
		Logger:    logger,
		Verbose:   cmdCfg.Verbose,
		Providers: providers,
	}

	return cfg, nil
}
