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

package queue

import "github.com/axelspringer/moppi/provider"
import "github.com/axelspringer/moppi/installer"

// Queue describes a queue
type Queue struct {
	Work   WorkQueue
	Worker WorkerQueue
}

// Worker is describing a worker to which work can be send
type Worker struct {
	ID       int
	Work     WorkQueue
	Worker   WorkerQueue
	QuitChan chan bool
}

// WorkQueue describes the queue for the work
type WorkQueue chan interface{}

// WorkerQueue describes the queue for the workers
type WorkerQueue chan chan interface{}

// Install describes an installment
type Install struct {
	Package   *provider.Package
	Installer *installer.Installer
}

// Uninstall describes an uninstallment
type Uninstall struct {
	Package   *provider.Package
	Installer *installer.Installer
}
