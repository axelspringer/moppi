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

// TransdecoderConfig is the configuration that is used to create a new transdecoder
// and allows customization of various aspects of decoding.
type TransdecoderConfig struct {
	// mutex *sync.RWMutex

	// ZeroFields, if set to true, will zero fields before writing them.
	// For example, a map will be emptied before decoded values are put in
	// it. If this is false, a map will be merged.
	ZeroFields bool

	// should be done later
	WeaklyTypedInput bool

	// Metadata is the struct that will contain extra metadata about
	// the decoding. If this is nil, then no metadata will be tracked.
	Metadata *Metadata

	// Result is a pointer to the struct that will contain the decoded
	// value.
	Result interface{}

	// The tag name that kvstructure reads for field names. This
	// defaults to "kvstructure"
	TagName string

	// Prefix is the prefix of the store
	Prefix string

	// KV is the kv used to retrieve the needed infos
	KV store.Store
}

// TranscoderConfig is the configuration that is used to create a new transcoder
// and allows customization of various aspects of decoding.
type TranscoderConfig struct {
	// mutex *sync.RWMutex

	// Result is a pointer to the struct that will contain the decoded
	// value.
	Input interface{}

	// Metadata is the struct that will contain extra metadata about
	// the decoding. If this is nil, then no metadata will be tracked.
	Metadata *Metadata

	// The tag name that kvstructure reads for field names. This
	// defaults to "kvstructure"
	TagName string

	// Prefix is the prefix of the store
	Prefix string

	// KV is the kv used to retrieve the needed infos
	KV store.Store
}

// A Transdecoder takes a raw interface value and turns it into structured data
type Transdecoder struct {
	config *TransdecoderConfig
}

// A Transcoder takes a raw interface and puts it into a kv structure
type Transcoder struct {
	config *TranscoderConfig
}

// Metadata contains information about decoding a structure that
// is tedious or difficult to get otherwise.
type Metadata struct {
	// Keys are the keys of the structure which were successfully decoded
	Keys []string

	// Unused is a slice of keys that were found in the raw value but
	// weren't decoded since there was no matching field in the result interface
	Unused []string
}

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
