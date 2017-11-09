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

package kv

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/axelspringer/moppi/provider"
	"github.com/docker/libkv/store"
)

// leadingSlash is adding a slash to the beginning
func leadingSlash(s string) string {
	prefix := "/"
	if !strings.HasPrefix(s, prefix) {
		return prefix + s
	}
	return s
}

// trailingSlash is adding a slash at the end
func trailingSlash(s string) string {
	suffix := "/"
	if !strings.HasSuffix(s, suffix) {
		return s + suffix
	}
	return s
}

// universePath gets a universe path from a prefix
func universePath(prefix string, universe string) string {
	return prefix + provider.MoppiUniverses + leadingSlash(universe)
}

// universeMetaPath gets the meta path of a universe
func universeMetaPath(prefix string, universe string) string {
	return universePath(prefix, universe) + provider.MoppiUniversesMeta
}

// reflectKVPair is reflecting a struct on a kv path
func reflectKVPair(s interface{}, path string, kv store.Store) {
	st := reflect.TypeOf(s)
	for i := 0; i < st.NumField(); i++ {
		key := trailingSlash(path) + strings.ToLower(st.Field(i).Name)
		if val, err := kv.Get(key); err == nil {
			f := reflect.ValueOf(s).Field(i)
			f.Set(reflect.ValueOf(val))
		}
	}

	fmt.Println(s)

	// fmt.Println(st.NumField())
	// for i := st.NumField(); i > 0; i-- {
	// 	fmt.Println(st.Field(i).Name)
	// }
}
