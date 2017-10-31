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
	"github.com/axelspringer/moppi/provider"
	"github.com/docker/libkv/store"
)

// ClientTLS holds TLS specific configurations as client
// CA, Cert and Key can be either path or file contents
type ClientTLS struct {
	CA                 string
	Cert               string
	Key                string
	InsecureSkipVerify bool
}

// Provider holds common configurations of key-value providers.
type Provider struct {
	provider.AbstractProvider `mapstructure:",squash" export:"true"`
	Endpoint                  string
	Prefix                    string
	TLS                       *ClientTLS
	Username                  string
	Password                  string
	storeType                 store.Backend
	kvClient                  store.Store
}
