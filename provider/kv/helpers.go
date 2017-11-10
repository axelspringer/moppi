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
	"strings"

	"github.com/axelspringer/moppi/provider"
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

// universesPath gets the path to the universes from a prefix
func universesPath(prefix string) string {
	return prefix + provider.MoppiUniverses
}

// universePath gets a universe path from a prefix
func universePath(prefix string, universe string) string {
	return prefix + provider.MoppiUniverses + leadingSlash(universe)
}

// universeMetaPath gets the meta path of a universe
func universeMetaPath(prefix string, universe string) string {
	return universePath(prefix, universe) + provider.MoppiUniversesMeta
}

// universePkgPath gets a package from a universe
func universePkgPath(prefix string, universe string, pkg string, rev string) string {
	return prefix + provider.MoppiUniverses + leadingSlash(universe) + provider.MoppiPackages + leadingSlash(pkg) + leadingSlash(rev)
}
