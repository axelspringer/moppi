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

package cmd

import (
	"github.com/spf13/cobra"
)

// newSetupCmd returns the setup subcommand
func newSetupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   setupCmd,
		Short: setupShort,
		Long:  setupLong,
		RunE:  runSetupE,
	}
}

// runSetupE contains the main functionality to setup Moppi
// in supported KVs (kvlib)
func runSetupE(c *cobra.Command, args []string) error {
	if ok, err := config.Etcd.Setup(); !ok {
		return err
	}

	return nil // noop
}
