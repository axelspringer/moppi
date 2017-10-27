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
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/axelspringer/moppi/install"
)

// writeJson emits a success and writes the JSON
func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(data)
	io.WriteString(w, string(json))
}

// writeError emits an error with a message and the error
func writeError(w http.ResponseWriter, msg string, status int, err error) {
	w.WriteHeader(status)
	m := fmt.Sprintf("%s :%v", msg, err)
	io.WriteString(w, m)
}

// writeJSONError emits the error a JSON
func writeErrorJSON(w http.ResponseWriter, msg string, status int, err error) {
	w.WriteHeader(status)
	writeJSON(w, &Error{msg, err.Error()})
}

// parseRequest parses a package request to Request
func parseRequest(r io.Reader) (*install.Request, error) {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		// here we should log request
		return nil, err
	}

	req, err := install.NewRequest(body)
	if err != nil {
		return nil, err
	}

	if req.Name == "" {
		return nil, PackageRequestFieldMissing("Name")
	}

	return req, err
}

// getHostname returns the hostname
func getHostname() (hostname string, err error) {
	host, err := os.Hostname()
	return host, err
}
