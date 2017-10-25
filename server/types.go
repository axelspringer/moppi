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

package server

import (
	"os"

	"github.com/axelspringer/moppi/cfg"
	"github.com/axelspringer/moppi/mesos"

	chronos "github.com/axelspringer/go-chronos"
	marathon "github.com/gambol99/go-marathon"
	log "github.com/sirupsen/logrus"
)

// Server holds the state of a new Server
type Server struct {
	cfg      *cfg.Config
	chronos  *chronos.Client
	log      *log.Logger
	marathon marathon.Marathon
	mesos    *mesos.Mesos
	signals  chan os.Signal
}

// Health describes the health of the api
type Health struct {
	Universes []*cfg.Universe
}

// Error contains an error of the api
type Error struct {
	Msg string
	Err string
}
