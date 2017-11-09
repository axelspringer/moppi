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
	"net/http"

	"github.com/axelspringer/moppi/provider"
	"github.com/zenazn/goji/web"
)

// metaUniversesPkg returns a universe package defintion
func (server *Server) pkgRevisions(c web.C, w http.ResponseWriter, _ *http.Request) {
	var pkgRequest provider.Request
	pkgRequest.Universe = c.URLParams["universe"]
	pkgRequest.Name = c.URLParams["name"]

	revs, err := server.provider.Revisions(&pkgRequest)
	if err != nil {
		writeErrorJSON(w, "Could not retrieve packages", 400, err)
		return
	}

	writeJSON(w, revs)
	return
}
