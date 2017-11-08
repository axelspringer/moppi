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
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/axelspringer/moppi/provider"
	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/prometheus/common/log"
)

// SetStoreType storeType setter, inherit from libkv
func (p *Provider) SetStoreType(storeType store.Backend) {
	p.storeType = storeType
}

// SetKVClient kvClient setter
func (p *Provider) SetKVClient(kvClient store.Store) {
	p.kvClient = kvClient
}

// Version returns the version of the universe
func (p *Provider) Version() (*store.KVPair, error) {
	ver, err := p.kvClient.Get(p.Prefix + provider.MoppiMetaVersion)
	return ver, err
}

// Package return a package
func (p *Provider) Package(req *provider.Request) (*provider.Package, error) {
	pkg := &provider.Package{}
	path := p.Prefix + provider.MoppiUniverses + "/" + req.Universe + provider.MoppiPackages + "/" + req.Name + "/" + req.Revision

	// info
	kvInfo, err := p.kvClient.Get(path + provider.MoppiPackage)
	if err != nil {
		return pkg, err
	}
	err = json.Unmarshal(kvInfo.Value, &pkg)
	if err != nil {
		return pkg, err
	}

	// read install
	kvInstall, err := p.kvClient.Get(path + provider.MoppiInstall)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(kvInstall.Value, &pkg.Install)
	if err != nil {
		return nil, err
	}

	// read uninstall
	kvUninstall, err := p.kvClient.Get(path + provider.MoppiUninstall)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(kvUninstall.Value, &pkg.Uninstall)
	if err != nil {
		return nil, err
	}

	// read chronos
	kvChronos, err := p.kvClient.Get(path + provider.MoppiChronos)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(kvChronos.Value, &pkg.Chronos)
	if err != nil {
		return nil, err
	}

	// read marathon
	kvMarathon, err := p.kvClient.Get(path + provider.MoppiMarathon)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(kvMarathon.Value, &pkg.Marathon)
	if err != nil {
		return nil, err
	}

	return pkg, nil
}

// Universes return all available universes
func (p *Provider) Universes() (*provider.Universes, error) {
	universes := make(provider.Universes, 0)
	path := p.Prefix + provider.MoppiUniverses + "/"

	// could this be refactored
	kvUniverses, err := p.kvClient.List(path)
	if err != nil {
		return nil, err
	}
	for _, universe := range kvUniverses {
		universes = append(universes, universe.Key) // return full key to universes
	}

	return &universes, nil
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
		strings.Split(p.Endpoint, ","),
		storeConfig,
	)
}

// CreateTLSConfig creates a TLS config from ClientTLS structures
func (clientTLS *ClientTLS) CreateTLSConfig() (*tls.Config, error) {
	var err error
	if clientTLS == nil {
		log.Warnf("clientTLS is nil")
		return nil, nil
	}
	caPool := x509.NewCertPool()
	if clientTLS.CA != "" {
		var ca []byte
		if _, errCA := os.Stat(clientTLS.CA); errCA == nil {
			ca, err = ioutil.ReadFile(clientTLS.CA)
			if err != nil {
				return nil, fmt.Errorf("Failed to read CA. %s", err)
			}
		} else {
			ca = []byte(clientTLS.CA)
		}
		caPool.AppendCertsFromPEM(ca)
	}

	cert := tls.Certificate{}
	_, errKeyIsFile := os.Stat(clientTLS.Key)

	if !clientTLS.InsecureSkipVerify && (len(clientTLS.Cert) == 0 || len(clientTLS.Key) == 0) {
		return nil, fmt.Errorf("TLS Certificate or Key file must be set when TLS configuration is created")
	}

	if len(clientTLS.Cert) > 0 && len(clientTLS.Key) > 0 {
		if _, errCertIsFile := os.Stat(clientTLS.Cert); errCertIsFile == nil {
			if errKeyIsFile == nil {
				cert, err = tls.LoadX509KeyPair(clientTLS.Cert, clientTLS.Key)
				if err != nil {
					return nil, fmt.Errorf("Failed to load TLS keypair: %v", err)
				}
			} else {
				return nil, fmt.Errorf("tls cert is a file, but tls key is not")
			}
		} else {
			if errKeyIsFile != nil {
				cert, err = tls.X509KeyPair([]byte(clientTLS.Cert), []byte(clientTLS.Key))
				if err != nil {
					return nil, fmt.Errorf("Failed to load TLS keypair: %v", err)

				}
			} else {
				return nil, fmt.Errorf("tls key is a file, but tls cert is not")
			}
		}
	}

	TLSConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caPool,
		InsecureSkipVerify: clientTLS.InsecureSkipVerify,
	}
	return TLSConfig, nil
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
