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
	"context"
	"io"
	"net"
	"net/http"
	"os"
	"sync"

	pb "github.com/axelspringer/moppi/api/v1"
	"google.golang.org/grpc"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/axelspringer/moppi/cfg"
	"github.com/axelspringer/moppi/installer"
	"github.com/axelspringer/moppi/queue"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/zenazn/goji/web"
)

// New returns a new instance of the server
func New(config *cfg.Config, exit chan bool, wg *sync.WaitGroup) (*Server, error) {
	var validate *validator.Validate

	// check version of repo in provider
	// if ok, err := config.Etcd.CheckVersion(version.Version); !ok {
	// 	return nil, err
	// }

	installer, err := installer.New(config)
	if err != nil {
		return nil, err
	}

	queue := queue.New(1)
	signals := make(chan os.Signal, 1)

	validate = validator.New()

	// server config
	server := &Server{
		listen:    config.Listen,
		installer: installer,
		signals:   signals,
		provider:  config.Etcd,
		queue:     queue,
		exit:      exit,
		validator: validate,
		wg:        wg,
	}

	return server, nil
}

// ping just returns a 200 and the okString
func (server *Server) ping(c web.C, w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, okString)
}

// health returns health infos about the api
func (server *Server) health(c web.C, w http.ResponseWriter, _ *http.Request) {
	return
}

// installPackage tries to install a package
func (server *Server) installPackage(w http.ResponseWriter, req *http.Request) {
	req.Header.Add("Accept", "application/json")

	packageRequest, err := parseRequest(req.Body)
	if err != nil {
		writeErrorJSON(w, "Could not parse the package request", 400, err)
		return
	}

	pkg, err := server.provider.GetPackage(packageRequest)
	if err != nil {
		writeErrorJSON(w, "Could not parse the package request", 400, err)
		return
	}

	// queue installment
	server.queue <- &queue.Install{Package: pkg, Installer: server.installer}

	w.WriteHeader(201)
}

// uninstallPackage tries to uninstall a package
func (server *Server) uninstallPackage(w http.ResponseWriter, req *http.Request) {
	req.Header.Add("Accept", "application/json")

	packageRequest, err := parseRequest(req.Body)
	if err != nil {
		writeErrorJSON(w, "Could not parse the package request", 400, err)
		return
	}

	pkg, err := server.provider.GetPackage(packageRequest)
	if err != nil {
		writeErrorJSON(w, "Could not parse the package request", 400, err)
		return
	}

	// queue installment
	server.queue <- &queue.Uninstall{Package: pkg, Installer: server.installer}

	w.WriteHeader(204)
}

// version returns the version of the repository
func (server *Server) version(w http.ResponseWriter, req *http.Request) {
	ver, err := server.provider.Version()
	if err != nil {
		w.WriteHeader(400)
		return
	}
	io.WriteString(w, string(ver.Value))
}

// Start is starting to serve the api
func (server *Server) Start(grpcAddr string) {
	var err error

	defer server.wg.Done()

	// add to waiting group
	server.wg.Add(1)

	// connect gateway to gRPC server
	mux := http.NewServeMux()
	// Start the REST gateway service, connecting to the gRPC server
	gwmux := runtime.NewServeMux()
	ctx := context.Background()
	dopts := []grpc.DialOption{grpc.WithInsecure()}
	err = pb.RegisterUniversesHandlerFromEndpoint(ctx, gwmux, grpcAddr, dopts)
	if err != nil {
		server.Stop()

		return
	}
	mux.Handle("/", gwmux)

	// register server listener
	httpLis, err := net.Listen("tcp", server.listen)
	if err != nil {
		panic(err)
	}

	srv := &http.Server{
		Addr:    server.listen,
		Handler: mux,
	}

	go func() {
		cfg.Log.Infof("Starting http server on %s\n", server.listen)
		err = srv.Serve(httpLis)
		if err != nil {
			server.Stop()

			return
		}
	}()

	// wait for exit
	<-server.exit
}

// Stop is stoping to serve the api
func (server *Server) Stop() {
	server.exit <- true

	return // noop
}
