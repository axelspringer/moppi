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

package kv

import (
	"github.com/axelspringer/moppi/install"
	"github.com/axelspringer/moppi/provider"
	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
)

// Provider holds common configurations of key-value providers.
type Provider struct {
	provider.Provider
	Endpoint  []string
	Prefix    string
	TLS       *ClientTLS
	Username  string
	Password  string
	storeType store.Backend
	kvClient  store.Store
}

// SetStoreType storeType setter, inherit from libkv
func (p *Provider) SetStoreType(storeType store.Backend) {
	p.storeType = storeType
}

// SetKVClient kvClient setter
func (p *Provider) SetKVClient(kvClient store.Store) {
	p.kvClient = kvClient
}

// Package return a package
func (p *Provider) Package(name string, version int) (*install.Package, error) {
	return &install.Package{}, nil
}

// Packages returns all available packages
func (p *Provider) Packages() ([]*install.Package, error) {
	return []*install.Package{}, nil
}

// CreateStore creates the K/V store
func (p *Provider) CreateStore(bucket string) (store.Store, error) {
	storeConfig := &store.Config{
		ConnectionTimeout: connectionTimeout,
		Bucket:            bucket,
		Username:          p.Username,
		Password:          p.Password,
	}

	if p.TLS != nil {
		var err error
		storeConfig.TLS, err = p.TLS.CreateTLSConfig()
		if err != nil {
			return nil, err
		}
	}

	return libkv.NewStore(
		p.storeType,
		p.Endpoint,
		storeConfig,
	)
}

// Provide provides the configuration to fernmelder via the configuration channel
// func (p *Provider) Provide(cfgChan chan<- types.ConfigMessage, pool *safe.Pool) error {
// 	watch kv for changes
// 	watcher := func() error {
// 		if _, err := p.kvClient.Exists(TOKEN); err != nil {
// 			return fmt.Errorf("Failed to test KV store connection: %v", err)
// 		}

// 		if we should watch the config in the KV
// 		if p.Watch {
// 			pool.Go(func(stop chan bool) {
// 				err := p.watchKv(cfgChan, p.Prefix, stop)
// 				if err != nil {
// 					log.Errorf("Cannot watch KV store: %v", err)
// 				}
// 			})
// 		}

// 		cfg := p.loadConfig()
// 		cfgChan <- types.ConfigMessage{
// 			ProviderName:  string(p.storeType),
// 			Configuration: cfg,
// 		}

// 		return nil
// 	}

// 	notify := func(err error, time time.Duration) {
// 		log.Errorf("KV connection error: %+v, retrying in %s", err, time)
// 	}

// 	err := backoff.RetryNotify(safe.OperationWithRecover(watcher), job.NewBackOff(backoff.NewExponentialBackOff()), notify)
// 	if err != nil {
// 		return fmt.Errorf("Cannot connect to KV server: %v", err)
// 	}

// 	return nil
// }

// func (p *Provider) watchKv(configurationChan chan<- types.ConfigMessage, prefix string, stop chan bool) error {
// 	create watcher
// 	watcher := func() error {
// 		events, err := p.kvClient.WatchTree(p.Prefix, make(chan struct{}))
// 		if err != nil {
// 			return fmt.Errorf("Failed to KV WatchTree: %v", err)
// 		}

// 		loop
// 		for {
// 			select {
// 			case <-stop:
// 				return nil
// 			case _, ok := <-events:
// 				if !ok {
// 					return errors.New("watchtree channel closed")
// 				}
// 				configuration := p.loadConfig()
// 				if configuration != nil {
// 					configurationChan <- types.ConfigMessage{
// 						ProviderName:  string(p.storeType),
// 						Configuration: configuration,
// 					}
// 				}
// 			}
// 		}
// 	}

// 	notify := func(err error, time time.Duration) {
// 		log.Errorf("KV connection error: %+v, retrying in %s", err, time)
// 	}

// 	err := backoff.RetryNotify(safe.OperationWithRecover(watcher), job.NewBackOff(backoff.NewExponentialBackOff()), notify)
// 	if err != nil {
// 		return fmt.Errorf("Cannot connect to KV server: %v", err)
// 	}

// 	return nil
// }

// func (p *Provider) loadConfig() *types.Configuration {
// 	return &types.Configuration{}
// }
