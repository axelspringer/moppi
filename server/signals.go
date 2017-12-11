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
	"os/signal"
	"syscall"
)

// configSignals configures signals the server should listen to
func (server *Server) configSignals() {
	signal.Notify(server.signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)
}

// watchSignals is watching configured signals
func (server *Server) watchSignals() {
	log := server.log
	for {
		sig := <-server.signals
		switch sig {
		case syscall.SIGUSR1:
			log.Infof("Nothing to see here yet")
		default:
			log.Infof("Shutting gracefully down")
			server.Stop()
		}
	}
}
