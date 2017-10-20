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

package server

import (
	"os"

	"github.com/axelspringer/moppi/mesos"

	marathon "github.com/gambol99/go-marathon"
)

// Server holds the state of a new Server
type Server struct {
	listen    *string
	universes []*Universe
	signals   chan os.Signal
	marathon  marathon.Marathon
	mesos     *mesos.Mesos
}

// Universe describes a complete universe
type Universe struct {
	Path *string
}

// Health describes the health of the api
type Health struct {
	Universes []*Universe
}

// Error contains an error of the api
type Error struct {
	Msg string
	Err string
}
