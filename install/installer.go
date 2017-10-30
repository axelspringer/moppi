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
)

// Request is requesting an install
func (i *Installer) Request(req *Request) error {

	// prepare
	work := WorkRequest{}

	store := i.stores[req.Universe]
	if store == nil {
		return errors.New("Could not fullfil the request")
	}
	prefix := store.prefix + moppiUniversePackages + "/" + req.Name + "/" + req.Revision

	// package
	kvPkg, err := store.kv.Get(prefix + moppiUniversePackage)
	if err != nil {
		return err
	}
	err = json.Unmarshal(kvPkg.Value, &work.Package)
	if err != nil {
		return err
	}

	// read install
	kvInstall, err := store.kv.Get(prefix + moppiUniverseInstall)
	if err != nil {
		return err
	}
	err = json.Unmarshal(kvInstall.Value, &work.Install)
	if err != nil {
		return err
	}

	// read uninstall
	kvUninstall, err := store.kv.Get(prefix + moppiUniverseUninstall)
	if err != nil {
		return err
	}
	err = json.Unmarshal(kvUninstall.Value, &work.Uninstall)
	if err != nil {
		return err
	}

	// read chronos
	kvChronos, err := store.kv.Get(prefix + moppiUniverseChronos)
	if err != nil {
		return err
	}
	err = json.Unmarshal(kvChronos.Value, &work.Chronos)
	if err != nil {
		return err
	}

	// read marathon
	kvMarathon, err := store.kv.Get(prefix + moppiUniverseMarathon)
	if err != nil {
		return err
	}
	err = json.Unmarshal(kvMarathon.Value, &work.Marathon)
	if err != nil {
		return err
	}

	workQueue <- work

	return nil
}
