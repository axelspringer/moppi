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
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

const (
	okString        = "OK"
	defaultListener = "localhost:8080"
)

// New returns a new instance of a Server
func New(listen string, universes []*Universe) *Server {
	if listen == "" {
		listen = defaultListener
	}

	return mustNew(&listen, universes)
}

// mustNew wraps the creation of a new Server
func mustNew(listen *string, universes []*Universe) *Server {
	signals := make(chan os.Signal, 1)

	return &Server{listen, universes, signals}
}

// ping just returns a 200 and the okString
func (server *Server) ping(c web.C, w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, okString)
}

// health returns health infos about the api
func (server *Server) health(c web.C, w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, &Health{
		Universes: server.universes,
	})
	return
}

// universes returns the current configured universes
func (server *Server) allUniverses(c web.C, w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, server.universes)
}

// installPackage tries to install a package
func (server *Server) installPackage(w http.ResponseWriter, req *http.Request) {
	req.Header.Add("Accept", "application/json")

	pkgRequest, err := parsePackageRequest(req.Body)
	if err != nil {
		writeErrorJSON(w, "Could not parse the package request", 400, err)
		return
	}

	fmt.Println(pkgRequest.Name)

	w.WriteHeader(201)
}

// uninstallPackage tries to uninstall a package
func (server *Server) uninstallPackage(w http.ResponseWriter, req *http.Request) {

}

// Start is starting to serve the api
func (server *Server) Start() {
	// setup signals
	server.configSignals()
	go server.watchSignals()

	// create listener
	listener, err := net.Listen("tcp", *server.listen)
	if err != nil {
		fmt.Printf("Could not create listener")
	}

	// setup routes
	goji.Get("/ping", server.ping)
	goji.Get("/health", server.health)
	goji.Get("/universes", server.allUniverses)

	goji.Post("/install", server.installPackage)
	goji.Post("/uninstall", server.uninstallPackage)

	// create server
	goji.ServeListener(listener)
}
