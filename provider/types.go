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

package provider

import (
	"github.com/axelspringer/go-chronos"
	"github.com/docker/libkv/store"
	"github.com/gambol99/go-marathon"
)

// Provider defines the interface to a Provider (e.g. etcd)
type Provider interface {
	Version() (*store.KVPair, error)
	Universe(req *Request) (*Universe, error)
	Universes() (*Universes, error)
	Revisions(req *Request) (*PackageRevisions, error)
	Package(req *Request) (*Package, error)
	Packages(req *Request) (*Packages, error)
	// Packages() (map[string]map[int]*install.Package, error)
}

// Universe describes a universe from endpoint /meta
type Universe struct {
	Description string `kvstructure:"description"`
	Name        string `kvstructure:"name"`
	Version     bool   `kvstructure:"version"`
	Test        Test   `kvstructure:"test,json"`
}

type Test struct {
	Test string `json:"test"`
}

// AbstractProvider is the base provider from which every provider inherits
type AbstractProvider struct {
	Enable bool
}

// RequestConfig describes a passed in config for installer Request
type RequestConfig struct {
}

// Request describes an request to the installer
type Request struct {
	Name     string        `json:"name"`
	Universe string        `json:"universe"`
	Revision string        `json:"revision"`
	Config   RequestConfig `json:"config"`
}

// Package describes a package in the universe
type Package struct {
	Version   string `json:"version"`
	Chronos   []chronos.Job
	Marathon  []marathon.Application
	Install   Install
	Uninstall Uninstall
}

// Install describes an installment (contained in install.json)
type Install struct {
	Marathon bool `json:"marathon"`
	Chronos  bool `json:"chronos"`
	Update   bool `json:"update"`
}

// Uninstall describes an uninstallment
type Uninstall struct {
	Marathon bool `json:"marathon"`
	Chronos  bool `json:"chronos"`
}

// Universes describes known universes
type Universes []Universe

// Packages describes the known packages in a universe
type Packages []string

// PackageRevisions describes known package revisions
type PackageRevisions []string
