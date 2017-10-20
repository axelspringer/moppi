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
	"github.com/axelspringer/moppi/provider/etcd"
	log "github.com/sirupsen/logrus"
)

// Config holds the persistent config of Moppi
type Config struct {
	Logger    *log.Logger
	Verbose   bool
	Providers *Providers
}

// Providers holds all available providers
type Providers struct {
	Etcd etcd.Provider
}

// CmdConfig holds the needed config of the command
type CmdConfig struct {
	Logger  *log.Logger
	Verbose bool
}