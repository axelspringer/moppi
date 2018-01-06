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
	"os"
	"os/signal"
	"syscall"

	"github.com/axelspringer/pitti/cfg"
)

// gracefulShutdown handles a graceful shutdown
func waitGracefulShutdown() {
	cfg.Log.Infof("Graceful shutdown")

	close(exit) // close all routines

	wg.Wait() // wait for exit
	close(shutdown)
}

// watchdog is watching important syscalls
func watchdog() {
	sys := make(chan os.Signal, 1) // create new channel for syscalls
	signal.Notify(sys, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)

	defer close(sys)
	defer waitGracefulShutdown()

	for {
		select {
		case <-exit: // wait for signal on exit channel
			return
		case sig := <-sys:
			switch sig {
			case syscall.SIGUSR1, syscall.SIGINT:
				return
			default:
			}
		}
	}
}
