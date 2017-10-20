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

package mesos

import "github.com/dghubble/sling"

// Mesos holds a Mesos Client
type Mesos struct {
	httpClient *sling.Sling
}

// State holds the current State of Mesos
type State struct {
	CompletedFrameworks    []*Framework `json:"completed_frameworks"`
	Frameworks             []*Framework `json:"frameworks"`
	UnregisteredFrameworks []string     `json:"unregistered_frameworks"`
	Flags                  Flags        `json:"flags"`
}

// Flags holds the flags provided to Mesos
type Flags struct {
	Authenticate       string `json:"authenticate"`
	AuthenticateSlaves string `json:"authenticate_slaves"`
}

// Framework holds a Mesos Framework
type Framework struct {
	Name             string  `json:"name"`
	ID               string  `json:"id"`
	PID              string  `json:"pid"`
	Active           bool    `json:"active"`
	Hostname         string  `json:"hostname"`
	User             string  `json:"user"`
	RegisteredTime   float64 `json:"registered_time"`
	ReregisteredTime float64 `json:"reregistered_time"`
	Tasks            []*Task `json:"tasks"`
}

// Task holds a Mesos Framework Task
type Task struct {
	FrameworkID string `json:"framework_id"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	SlaveID     string `json:"slave_id"`
	State       string `json:"state"`
}
