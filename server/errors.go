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

import "fmt"

// PackageRequestFieldMissing is a new type that inherits error
type PackageRequestFieldMissing string

// Error returns a custom error
func (err PackageRequestFieldMissing) Error() string {
	return fmt.Sprintf("A field is missing: %v", string(err))
}
