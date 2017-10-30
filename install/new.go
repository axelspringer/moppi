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

package install

import (
	"net/http"
	"os"
	"strings"

	chronos "github.com/axelspringer/go-chronos"
	"github.com/axelspringer/moppi/cfg"
	"github.com/axelspringer/moppi/mesos"
	"github.com/axelspringer/moppi/provider/etcd"
	"github.com/axelspringer/moppi/version"
	marathon "github.com/gambol99/go-marathon"
)

// New is creates a new installer
func New(config *cfg.Config) (*Installer, error) {
	return mustNew(config)
}

// mustNew wraps the creation of a new installer
func mustNew(config *cfg.Config) (*Installer, error) {
	stores := make(map[string]*Store)

	// create kv stores
	for universe, storeCfg := range config.Universes {
		if storeCfg.Etcd.Endpoint != "" {
			var provider etcd.Provider
			provider.Prefix = strings.TrimPrefix(storeCfg.Etcd.Prefix, "/")
			provider.Endpoint = strings.Split(storeCfg.Etcd.Endpoint, ",")
			kvStore, err := provider.CreateStore(defaultBucket)

			// keep this for now
			if err != nil {
				panic(err)
			}

			// get version
			kvVersion, err := kvStore.Get(provider.Prefix + moppiUniverseVersion)
			if err != nil {
				panic(err) // TODO: should be an internal error
			}
			universeVersion := string(kvVersion.Value)

			// do intermediate checks
			if universeVersion != version.Version {
				config.Logger.Fatalf("Incompatible Universe found: %v (version: %v)", provider.Prefix, universeVersion)
				os.Exit(-1)
			}

			// Todo: check universe

			// map store to universe
			stores[universe] = &Store{
				prefix:   provider.Prefix,
				endpoint: provider.Endpoint,
				kv:       kvStore,
			}

			// for now, only configure one
			break
		}
	}

	// create new Marathon client
	marathonConfig := marathon.NewDefaultConfig()
	marathonConfig.URL = config.CmdConfig.Marathon
	marathonClient, err := marathon.NewClient(marathonConfig)
	if err != nil {
		return nil, err
	}

	// creating new Mesos client
	httpClient := &http.Client{}
	mesosClient := mesos.New(httpClient, config.CmdConfig.Mesos)

	// creating new Chronos client
	chronosClient := chronos.New(config.CmdConfig.Chronos, httpClient)

	// creating installer
	installer := &Installer{
		cfg:      config,
		marathon: marathonClient,
		mesos:    mesosClient,
		chronos:  chronosClient,
		stores:   stores,
	}
	StartDispatcher(4, installer)

	return installer, nil
}
