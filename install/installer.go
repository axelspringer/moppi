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

package install

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/gambol99/go-marathon"
)

// Request is requesting an install
func (i *Installer) Request(req *Request) error {
	store := i.stores[req.Universe]
	if store == nil {
		return errors.New("Could not fullfil the request")
	}

	prefix := store.prefix + moppiUniversePackages + "/" + req.Name + "/" + req.Revision
	// package
	val, err := store.kv.Get(prefix + moppiUniversePackage)
	if err != nil {
		return err
	}

	fmt.Println(string(val.Value))

	// chronos
	val, err = store.kv.Get(prefix + moppiUniverseChronos)
	if err != nil {
		return err
	}

	fmt.Println(string(val.Value))

	// marathon
	val, err = store.kv.Get(prefix + moppiUniverseMarathon)
	if err != nil {
		return err
	}

	fmt.Println(string(val.Value))

	app := &marathon.Application{}
	err = json.Unmarshal(val.Value, &app)

	if _, err := i.marathon.CreateApplication(app); err != nil {
		log.Fatalf("Failed to create application: %s, error: %s", app, err)
	} else {
		log.Printf("Created the application: %s", app)
	}

	return nil
}
