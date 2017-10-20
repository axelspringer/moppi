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

	log "github.com/sirupsen/logrus"
)

// ClientTLS holds TLS specific configurations as client
// CA, Cert and Key can be either path or file contents
type ClientTLS struct {
	CA                 string
	Cert               string
	Key                string
	InsecureSkipVerify bool
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
