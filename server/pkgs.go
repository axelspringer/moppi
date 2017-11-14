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
	"io"
	"net/http"
	"strconv"

	"github.com/axelspringer/moppi/provider"
	"github.com/zenazn/goji/web"
)

// getPkg returns a universe package
func (server *Server) getPkg(c web.C, w http.ResponseWriter, _ *http.Request) {
	var pkgRequest provider.Request
	pkgRequest.Universe = c.URLParams["universe"]
	pkgRequest.Name = c.URLParams["name"]
	pkgRequest.Revision = c.URLParams["revision"]

	revs, err := server.provider.GetPackage(&pkgRequest)
	if err != nil {
		writeErrorJSON(w, "Could not retrieve packages", http.StatusBadRequest, err)
		return
	}

	writeJSON(w, revs)
	return
}

// getPkgRevisions returns a universe package defintion
func (server *Server) getPkgRevisions(c web.C, w http.ResponseWriter, _ *http.Request) {
	var pkgRequest provider.Request
	pkgRequest.Universe = c.URLParams["universe"]
	pkgRequest.Name = c.URLParams["name"]

	revs, err := server.provider.GetRevisions(&pkgRequest)
	if err != nil {
		writeErrorJSON(w, "Could not retrieve revisions", http.StatusBadRequest, err)
		return
	}

	writeJSON(w, revs)
	return
}

// getPkgs returns all the packages in a universe
func (server *Server) getPkgs(c web.C, w http.ResponseWriter, _ *http.Request) {
	var pkgRequest provider.Request
	pkgRequest.Universe = c.URLParams["universe"]

	pkgs, err := server.provider.GetPackages(&pkgRequest)
	if err != nil {
		writeErrorJSON(w, "Could not retrieve packages", http.StatusBadRequest, err)
		return
	}

	writeJSON(w, pkgs)
	return
}

// deletePkg deletes a package from a universe
func (server *Server) deletePkg(c web.C, w http.ResponseWriter, req *http.Request) {
	var pkgRequest provider.Request
	pkgRequest.Universe = c.URLParams["universe"]
	pkgRequest.Name = c.URLParams["name"]

	if err := server.provider.DeletePackage(&pkgRequest); err != nil {
		writeErrorJSON(w, "Could not delete the package", http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// deletePkgRevision deletes a revision of a package in a universe
func (server *Server) deletePkgRevision(c web.C, w http.ResponseWriter, req *http.Request) {
	var pkgRequest provider.Request
	pkgRequest.Universe = c.URLParams["universe"]
	pkgRequest.Name = c.URLParams["name"]
	pkgRequest.Revision = c.URLParams["revision"]

	if err := server.provider.DeletePackageRevision(&pkgRequest); err != nil {
		writeErrorJSON(w, "Could not delete the revision", http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// createPkgRevision creates a new revision of a package in a universe
func (server *Server) createPkgRevision(c web.C, w http.ResponseWriter, req *http.Request) {
	var err error
	var pkg provider.Package
	var pkgRequest provider.Request
	pkgRequest.Universe = c.URLParams["universe"]
	pkgRequest.Name = c.URLParams["name"]

	// empty body
	if req.Body == nil {
		writeErrorJSON(w, "Please send a request body", http.StatusBadRequest, err)
		return
	}

	// decode request
	if err := json.NewDecoder(req.Body).Decode(&pkg); err != nil {
		writeErrorJSON(w, "Could not parse the request", http.StatusBadRequest, err)
		return
	}

	// create new revision
	rev, err := server.provider.CreatePackageRevision(&pkgRequest, &pkg)
	if err != nil {
		writeErrorJSON(w, "Could not create a new package revision", http.StatusBadRequest, err)
	}

	// simply write the newly created revision
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, strconv.Itoa(*rev))
}
