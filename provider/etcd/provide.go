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
	"errors"

	pb "github.com/axelspringer/moppi/api/v1"
	"github.com/axelspringer/moppi/provider"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/etcd"
)

var _ provider.Provider = (*Provider)(nil)

// GetPackage gets a package from etcd
func (p *Provider) GetPackage(req *provider.Request) (*provider.Package, error) {
	return p.Provider.GetPackage(req)
}

// GetPackages gets all the packages from a universe in the etcd
func (p *Provider) GetPackages(req *provider.Request) (*provider.Packages, error) {
	return p.Provider.GetPackages(req)
}

// DeletePackage deletes a package from the universe in the etcd
func (p *Provider) DeletePackage(req *provider.Request) error {
	return p.Provider.DeletePackage(req)
}

// DeletePackageRevision deletes a revision of a package in the etcd
func (p *Provider) DeletePackageRevision(req *provider.Request) error {
	return p.Provider.DeletePackageRevision(req) // @todo, this is happy path
}

// GetRevisions gets all the packages revisions from a universe
func (p *Provider) GetRevisions(req *provider.Request) (*provider.PackageRevisions, error) {
	return p.Provider.GetRevisions(req)
}

// CreateUniverse creates a new universe
func (p *Provider) CreateUniverse(u *provider.Universe) error {
	return p.Provider.CreateUniverse(u)
}

// DeleteUniverse creates a new universe
func (p *Provider) DeleteUniverse(req *provider.Request) error {
	return p.Provider.DeleteUniverse(req)
}

// Version gets the meta version of the moppi repo
func (p *Provider) Version() (*store.KVPair, error) {
	return p.Provider.Version()
}

// GetUniverse gets the meta info of a universe
func (p *Provider) GetUniverse(req *pb.PackageRequest) (*pb.Universe, error) {
	return p.Provider.GetUniverse(req)
}

// GetUniverses gets all the known universes
func (p *Provider) GetUniverses() (*provider.Universes, error) {
	return p.Provider.GetUniverses()
}

// CheckVersion checks the meta version of the moppi repo
func (p *Provider) CheckVersion(moppiVersion string) (bool, error) {
	version, err := p.Version()
	if err != nil {
		return false, err
	}
	// TODO: real error
	return string(version.Value) == moppiVersion, errors.New("Version not matching")
}

// Setup checks for the setup of the moppi repo
func (p *Provider) Setup() (bool, error) {
	return p.Provider.Setup()
}

// CreateStore creates the etcd store
func (p *Provider) CreateStore(bucket string) (store.Store, error) {
	p.SetStoreType(store.ETCD)
	etcd.Register()

	return p.Provider.CreateStore(bucket)
}
