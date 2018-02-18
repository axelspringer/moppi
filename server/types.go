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
	"net"
	"os"
	"sync"

	"github.com/axelspringer/moppi/installer"

	"github.com/axelspringer/moppi/provider/etcd"
	"github.com/axelspringer/moppi/queue"
	validator "gopkg.in/go-playground/validator.v9"
)

// Server holds the state of a new Server
type Server struct {
	signals   chan os.Signal
	installer *installer.Installer
	listener  net.Listener
	provider  etcd.Provider
	queue     queue.WorkQueue
	exit      chan bool
	validator *validator.Validate
	wg        *sync.WaitGroup
}

// Error contains an error of the api
type Error struct {
	Msg string
	Err string
}
