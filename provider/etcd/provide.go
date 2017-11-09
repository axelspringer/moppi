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

package etcd

import (
	"github.com/axelspringer/moppi/provider"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/etcd"
)

var _ provider.Provider = (*Provider)(nil)

// Package gets a package from etcd
func (p *Provider) Package(req *provider.Request) (*provider.Package, error) {
	return p.Provider.Package(req)
}

// Packages gets all the packages from a universe in the etcd
func (p *Provider) Packages(req *provider.Request) (*provider.Packages, error) {
	return p.Provider.Packages(req)
}

// Revisions gets all the packages revisions from a universe
func (p *Provider) Revisions(req *provider.Request) (*provider.PackageRevisions, error) {
	return p.Provider.Revisions(req)
}

// Version gets the meta version of the moppi repo
func (p *Provider) Version() (*store.KVPair, error) {
	return p.Provider.Version()
}

// Universe gets the meta info of a universe
func (p *Provider) Universe(req *provider.Request) (*provider.Universe, error) {
	return p.Provider.Universe(req)
}

// Universes gets all the known universes
func (p *Provider) Universes() (*provider.Universes, error) {
	return p.Provider.Universes()
}

// CheckVersion checks the meta version of the moppi repo
func (p *Provider) CheckVersion(moppiVersion string) (bool, error) {
	version, err := p.Version()
	if err != nil {
		return false, err
	}
	// TODO: real error
	return string(version.Value) == moppiVersion, nil
}

// CreateStore creates the etcd store
func (p *Provider) CreateStore(bucket string) (store.Store, error) {
	p.SetStoreType(store.ETCD)
	etcd.Register()

	return p.Provider.CreateStore(bucket)
}
