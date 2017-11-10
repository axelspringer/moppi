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
	path := universePkgPath(p.Prefix, req.Universe, req.Name, req.Revision)
	var pkg provider.Package

	if err := Transcode(&pkg, path, p.kvClient); err != nil {
		return &pkg, err
	}

	return &pkg, nil
}

// Packages returns all the packages in a universe
func (p *Provider) Packages(req *provider.Request) (*provider.Packages, error) {
	pkgs := make(provider.Packages, 0)
	path := p.Prefix + provider.MoppiUniverses + leadingSlash(req.Universe) + provider.MoppiPackages
	path = trailingSlash(path)

	// could this be refactored?
	kvPackages, err := p.kvClient.List(path)
	if err != nil {
		return nil, err
	}

	for _, pkg := range kvPackages {
		pkgs = append(pkgs, strings.TrimPrefix(pkg.Key, path))
	}

	return &pkgs, nil
}

// Revisions returns all the revisions of a package in a univers
func (p *Provider) Revisions(req *provider.Request) (*provider.PackageRevisions, error) {
	revs := make(provider.PackageRevisions, 0)
	path := p.Prefix + provider.MoppiUniverses + leadingSlash(req.Universe) + provider.MoppiPackages + leadingSlash(req.Name)
	path = trailingSlash(path)

	// could this be refactored?
	kvRevisions, err := p.kvClient.List(path)
	if err != nil {
		return nil, err
	}

	for _, rev := range kvRevisions {
		revs = append(revs, strings.TrimPrefix(rev.Key, path))
	}

	return &revs, nil
}

// Universes return all available universes
func (p *Provider) Universes() (*provider.Universes, error) {
	universes := make(provider.Universes, 0)
	path := universesPath(p.Prefix)

	// could this be refactored?
	kvUniverses, err := p.kvClient.List(path)
	if err != nil {
		return nil, err
	}

	for _, universe := range kvUniverses {
		universes = append(universes, strings.TrimPrefix(universe.Key, path))
	}

	return &universes, nil
}

// Universe return the meta infos of a universe
func (p *Provider) Universe(req *provider.Request) (*provider.Universe, error) {
	path := universeMetaPath(p.Prefix, req.Universe)
	var universe provider.Universe

	if err := Transcode(&universe, path, p.kvClient); err != nil {
		return &universe, err
	}

	return &universe, nil
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
