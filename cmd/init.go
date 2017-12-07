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

// NewInitCmd returns the new subcommand version
func NewInitCmd() *cobra.Command {
	// suprise, create a new version command
	return mustNewInitCmd()
}

// mustNew wraps the creation of a new version cmd
func mustNewInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initializes Moppi in a support KV",
		Long:  ``,
		RunE:  runInitE,
	}
}

// RunE is the function to be executed to run the command
func runInitE(c *cobra.Command, args []string) error {
	if ok, err := config.Etcd.Setup(); !ok {
		return err
	}

	// save for later
	return nil
}
