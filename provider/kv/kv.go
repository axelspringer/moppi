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
	"strconv"
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

// Version returns the version of the
func (p *Provider) Version() (*store.KVPair, error) {
	ver, err := p.kvClient.Get(p.Prefix + provider.MoppiMetaVersion)
	return ver, err
}

// GetPackage return a package
func (p *Provider) GetPackage(req *provider.Request) (*provider.Package, error) {
	path := universePkgPath(p.Prefix, req.Universe, req.Name, req.Revision)
	var pkg provider.Package

	if err := Transdecode(&pkg, path, p.kvClient); err != nil {
		return &pkg, err
	}

	return &pkg, nil
}

// GetPackages returns all the packages in a universe
func (p *Provider) GetPackages(req *provider.Request) (*provider.Packages, error) {
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

// GetRevisions returns all the revisions of a package in a univers
func (p *Provider) GetRevisions(req *provider.Request) (*provider.PackageRevisions, error) {
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

// CreateUniverse creates a new universe
func (p *Provider) CreateUniverse(u *provider.Universe) error {
	path := universeMetaPath(p.Prefix, strings.ToLower(u.Name))

	return Transcode(&u, path, p.kvClient)
}

// DeleteUniverse deletes a universe
func (p *Provider) DeleteUniverse(req *provider.Request) error {
	path := universePath(p.Prefix, req.Universe)

	return p.kvClient.DeleteTree(path)
}

// DeletePackage deletes a package form a universe
func (p *Provider) DeletePackage(req *provider.Request) error {
	path := universePkgBasePath(p.Prefix, req.Universe, req.Name)

	return p.kvClient.DeleteTree(path)
}

// CreatePackageRevision creates a new package revision in a universe
// and returns the new revision and any error
func (p *Provider) CreatePackageRevision(req *provider.Request, pkg *provider.Package) (*int, error) {
	// check if package exists
	path := universePkgBasePath(p.Prefix, req.Universe, req.Name)
	rev := 1

	// if not exists, create package
	if ok, _ := p.kvClient.Exists(path); !ok {
		// create package in universe
		// @todo check if universe exists
		revPath := universePkgPath(p.Prefix, req.Universe, req.Name, fmt.Sprint(rev))

		// Transcode
		if err := Transcode(&pkg, revPath, p.kvClient); err != nil {
			return nil, err
		}

		return &rev, nil
	}

	// get all revisions
	revs, err := p.kvClient.List(path) // is sorted from etcd
	if err != nil {
		return nil, err
	}

	kvRevPath := strings.Split(string(revs[len(revs)-1].Key), "/")
	kvRev := kvRevPath[len(kvRevPath)-1]

	rev, err = strconv.Atoi(kvRev)
	if err != nil {
		return nil, err
	}
	revPath := universePkgPath(p.Prefix, req.Universe, req.Name, fmt.Sprint(rev+1))

	// Transcode
	if err := Transcode(&pkg, revPath, p.kvClient); err != nil {
		return nil, err
	}

	return &rev, nil
}

// DeletePackageRevision deletes a package revision from a universe
func (p *Provider) DeletePackageRevision(req *provider.Request) error {
	path := universePkgPath(p.Prefix, req.Universe, req.Name, req.Revision)

	return p.kvClient.DeleteTree(path)
}

// GetUniverses return all available universes
func (p *Provider) GetUniverses() (*provider.Universes, error) {
	universes := make(provider.Universes, 0)
	universesPath := universesPath(p.Prefix)

	// could this be refactored?
	kvUniverses, err := p.kvClient.List(universesPath)
	if err != nil {
		return nil, err
	}

	for _, kvUniverse := range kvUniverses {
		universePath := kvUniverse.Key + provider.MoppiUniversesMeta
		var universe provider.Universe

		err := Transdecode(&universe, universePath, p.kvClient)
		if err != nil {
			return &universes, err
		}

		universe.Href = kvUniverse.Key
		universes = append(universes, universe)
	}

	return &universes, nil
}

// GetUniverse return the meta infos of a universe
func (p *Provider) GetUniverse(req *provider.Request) (*provider.Universe, error) {
	path := universeMetaPath(p.Prefix, req.Universe)
	var universe provider.Universe

	if err := Transdecode(&universe, path, p.kvClient); err != nil {
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
