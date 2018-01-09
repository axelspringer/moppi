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
	"encoding/json"
	"net/http"

	"github.com/axelspringer/moppi/provider"
	"github.com/zenazn/goji/web"
)

func (server *Server) getUniverse(c web.C, w http.ResponseWriter, _ *http.Request) {
	var pkgRequest provider.Request
	pkgRequest.Universe = c.URLParams["universe"]

	universe, err := server.provider.GetUniverse(&pkgRequest)
	if err != nil {
		writeErrorJSON(w, "Could not retrieve the universes", 400, err)
		return
	}

	writeJSON(w, universe)
	return
}

// createUniverse is creating a new universe in the KV
func (server *Server) createUniverse(c web.C, w http.ResponseWriter, req *http.Request) {
	var universe provider.Universe
	var err error
	if req.Body == nil {
		writeErrorJSON(w, "Please send a request body", http.StatusBadRequest, err)
		return
	}

	err = json.NewDecoder(req.Body).Decode(&universe)
	if err != nil {
		writeErrorJSON(w, "Could not parse the request", http.StatusBadRequest, err)
		return
	}

	err = server.validator.Struct(universe)
	if err != nil {
		writeErrorJSON(w, "Could not create a new universe", http.StatusBadRequest, err)
		return
	}

	err = server.provider.CreateUniverse(&universe)
	if err != nil {
		writeErrorJSON(w, "Could not create a new universe", http.StatusBadGateway, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// deleteUniverse is deleting a universe in the KV
func (server *Server) deleteUniverse(c web.C, w http.ResponseWriter, req *http.Request) {
	var pkgRequest provider.Request
	pkgRequest.Universe = c.URLParams["universe"]

	if err := server.provider.DeleteUniverse(&pkgRequest); err != nil {
		writeErrorJSON(w, "Could not delete the universe", 400, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

// getUniverses returns all the known universes
func (server *Server) getUniverses(w http.ResponseWriter, _ *http.Request) {
	universes, err := server.provider.GetUniverses()
	if err != nil {
		writeErrorJSON(w, "Could not retrieve the universes", 400, err)
		return
	}

	writeJSON(w, universes)
	return
}
