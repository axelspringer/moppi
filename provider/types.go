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
	Setup() (bool, error)
	CreateUniverse(u *Universe) error
	GetUniverse(req *Request) (*Universe, error)
	GetUniverses() (*Universes, error)
	GetRevisions(req *Request) (*PackageRevisions, error)
	GetPackage(req *Request) (*Package, error)
	GetPackages(req *Request) (*Packages, error)
	// Packages() (map[string]map[int]*install.Package, error)
}

// Version describes the version of a moppi installation
type Version string

// Meta
type Meta struct {
	Version string `kvstructure:"version"`
}

// Config
type Config struct {
	*Meta `kvstructure:"meta"`
}

// Universe describes a universe from endpoint /meta
type Universe struct {
	Description string `kvstructure:"description" json:"description" validate:"required,min=1"`
	Name        string `kvstructure:"name" json:"name" validate:"required,min=1"`
	Version     string `kvstructure:"version" json:"version" validate:"required"`
	Href        string `json:"href"`
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
	Chronos   []chronos.Job          `kvstructure:"chronos,json" json:"chronos" validate:"required"`
	Marathon  []marathon.Application `kvstructure:"marathon,json" json:"marathon" validate:"required"`
	Install   Install                `kvstructure:"install,json" json:"install" validate:"required"`
	Uninstall Uninstall              `kvstructure:"uninstall,json" json:"uninstall" validate:"required"`
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
