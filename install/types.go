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

import marathon "github.com/gambol99/go-marathon"

// PackageRequest describes an request for installing/uninstalling
// a package from different universes
type PackageRequest struct {
	Name     string               `json:"name"`
	Universe string               `json:"universe"`
	Version  string               `json:"version"`
	Config   PackageRequestConfig `json:"config"`
}

// PackageRequestConfig describes the config of a package installment
type PackageRequestConfig struct {
}

// PackageInstall describes a running package installment
type PackageInstall struct {
	Request *PackageRequest
}

// Install describes an installment (contained in install.json)
type Install struct {
}

// Uninstall describes an Un-installment (contained in uninstall.json)
type Uninstall struct {
}

// Info describes the package (contained in package.json)
type Info struct {
}

// @TODO Chronos needs to be integrated, type, package and client

// Package describes a package from an universe
type Package struct {
	Marathon  marathon.Application
	Info      Info
	Install   Install
	Uninstall Uninstall
}
