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

package cfg

// New returns a new Config
func New() (*Config, error) {
	// config
	var cfg = new(Config)

	return cfg, nil // noop
}

func (c *Config) Init() {
	// initialize provider
	store, err := c.Etcd.CreateStore(defaultBucket)
	if err != nil {
		Log.Fatalf("Initializing provider: %v", err)
	}
	c.Etcd.SetKVClient(store)
}
