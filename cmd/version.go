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
	"fmt"

	"github.com/axelspringer/moppi/version"
	"github.com/spf13/cobra"
)

// newVersionCmd returns a new version command
func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Moppi",
		Long:  `All software has versions. This is Moppis`,
		RunE:  runVersionE,
	}
}

// runVersionE prints the current version of Moppi
func runVersionE(c *cobra.Command, args []string) error {
	// print the version string
	fmt.Println("v" + version.Version)

	return nil // noop
}
