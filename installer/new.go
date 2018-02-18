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

package installer

import (
	"net/http"

	chronos "github.com/axelspringer/go-chronos"
	"github.com/axelspringer/moppi/cfg"
	"github.com/axelspringer/moppi/mesos"
	marathon "github.com/gambol99/go-marathon"
)

// New is creates a new installer
func New(config *cfg.Config) (*Installer, error) {
	// create new Marathon client
	marathonConfig := marathon.NewDefaultConfig()
	marathonConfig.URL = config.Marathon
	marathonClient, err := marathon.NewClient(marathonConfig)
	if err != nil {
		return nil, err
	}

	// creating new Mesos client
	httpClient := &http.Client{}
	mesosClient := mesos.New(httpClient, config.Mesos)

	// creating new Chronos client
	chronosClient := chronos.New(config.Chronos, httpClient)

	// creating installer
	installer := &Installer{
		Chronos:  chronosClient,
		Log:      cfg.Log,
		Marathon: marathonClient,
		Mesos:    mesosClient,
	}

	return installer, nil
}
