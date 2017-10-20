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

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dghubble/sling"
)

// New creates a new Mesos client
func New(httpClient *http.Client, url string) *Mesos {
	return mustNew(httpClient, &url)
}

// mustNew wraps the creation of a new Mesos client
func mustNew(httpClient *http.Client, url *string) *Mesos {
	sling := sling.New().Client(httpClient).Base(*url)

	return &Mesos{
		httpClient: sling,
	}
}

// State gets the state of the Mesos master
func (m *Mesos) State() (*State, error) {
	state := new(State)

	_, err := m.httpClient.Get("/master/state").ReceiveSuccess(state)
	if err != nil {
		return nil, err
	}

	return state, err
}

// Frameworks gets the mesos frameworks
func (m *Mesos) Frameworks() ([]*Framework, error) {
	state, err := m.State()

	if err != nil {
		return []*Framework{}, err
	}

	return state.Frameworks, nil
}

// CompletedFrameworks gets the completed frameworks
func (m *Mesos) CompletedFrameworks() ([]*Framework, error) {
	state, err := m.State()

	if err != nil {
		return []*Framework{}, err
	}

	return state.CompletedFrameworks, nil
}

// Shutdown is shutting down a framework
func (m *Mesos) Shutdown(frameworkID string) error {
	data := fmt.Sprintf("frameworkId=%s", frameworkID)
	body := strings.NewReader(data)

	res, err := m.httpClient.Post("/master/teardown").Body(body).Request()
	if res.Response.StatusCode == 200 {
		return nil
	}
	return err
}

// SearchFrameworks is searching frameworks in Mesos
func (m *Mesos) FindFrameworks(name string) ([]*Framework, error) {
	state, err := m.State()
	if err != nil {
		return []*Framework{}, err
	}

	matches := make(map[string]*Framework)
	for _, framework := range state.Frameworks {
		if strings.EqualFold(framework.Name, name) {
			matches[framework.ID] = framework
		}
	}

	var frameworks []*Framework
	for _, framework := range matches {
		frameworks = append(frameworks, framework)
	}

	return frameworks, nil
}

// FindFramework is trying to find a framework by name
func (m *Mesos) FindFramework(name string) (*Framework, error) {
	frameworks, err := m.FindFrameworks(name)
	if err != nil {
		return nil, err
	}

	if len(frameworks) == 0 {
		return nil, nil
	}

	if len(frameworks) > 1 {
		return nil, fmt.Errorf("There are %d %s frameworks", len(frameworks), name)
	}

	return frameworks[0], nil
}

// implement search || strings.Contains(framework.Name, name)
