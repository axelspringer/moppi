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
	"log"
	"time"

	chronos "github.com/axelspringer/go-chronos"
	"github.com/axelspringer/moppi/cfg"
	"github.com/axelspringer/moppi/mesos"
	"github.com/docker/libkv/store"
	marathon "github.com/gambol99/go-marathon"
)

// Request describes an request for installing/uninstalling
// a package from different universes
type Request struct {
	Name     string        `json:"name"`
	Universe string        `json:"universe"`
	Revision string        `json:"revision"`
	Config   RequestConfig `json:"config"`
}

// RequestConfig describes a passed in config for installer Request
type RequestConfig struct {
}

// Installer describes the config of a package installment
type Installer struct {
	cfg      *cfg.Config
	chronos  *chronos.Client
	log      *log.Logger
	marathon marathon.Marathon
	mesos    *mesos.Mesos
	stores   map[string]*Store
}

// Store describes a store
type Store struct {
	prefix   string
	endpoint []string
	kv       store.Store
}

// Install describes an installment (contained in install.json)
type Install struct {
	Marathon bool `json:"marathon"`
	Chronos  bool `json:"chronos"`
}

// Uninstall describes an uninstallment
type Uninstall struct{}

// Package describes an package
type Package struct {
}

// WorkRequest is describing a workload to do
type WorkRequest struct {
	Name      string
	Install   Install
	Uninstall Uninstall
	Marathon  marathon.Application
	Chronos   chronos.Job
	Delay     time.Duration
}

// Worker is describing a worker to which work can be send
type Worker struct {
	ID          int
	Work        chan WorkRequest
	WorkerQueue chan chan WorkRequest
	QuitChan    chan bool
	Installer   *Installer
}
