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

package queue

import "time"

func install(i *Install) error {
	installer := i.Installer
	pkg := i.Package

	// deploy marathon
	if pkg.Install.Marathon {
		for _, marathon := range pkg.Marathon {
			if _, err := installer.Marathon.CreateApplication(&marathon); err != nil {
				return err
			}

			if err := installer.Marathon.WaitOnApplication(marathon.ID, 10*time.Second); err != nil {
				return err
			}
		}
	}

	// deploy chronos
	if pkg.Install.Chronos {
		for _, chronos := range pkg.Chronos {
			if ok, _, err := installer.Chronos.Job.New(&chronos); !ok {
				return err
			}
		}
	}

	return nil
}

func uninstall(i *Uninstall) error {

	return nil
}
