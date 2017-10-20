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
	"github.com/axelspringer/moppi/install"
	"github.com/axelspringer/moppi/provider"
	"github.com/axelspringer/moppi/provider/kv"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/etcd"
)

var _ provider.Provider = (*Provider)(nil)

// Provider holds configurations of the provider.
type Provider struct {
	kv.Provider
}

// New creates a new etcd provider
func New(bucket string) error {
	return mustNew(&bucket)
}

func mustNew(bucket *string) error {
	p := new(Provider)
	store, err := p.CreateStore(*bucket)
	if err != nil {
		return err
	}
	p.SetKVClient(store)

	return nil
}

// Package gets a package from etcd
func (p *Provider) Package(name string, version int) (*install.Package, error) {
	return p.Provider.Package(name, version)
}

// CreateStore creates the etcd store
func (p *Provider) CreateStore(bucket string) (store.Store, error) {
	p.SetStoreType(store.ETCD)
	etcd.Register()

	return p.Provider.CreateStore(bucket)
}
