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

package install

import "fmt"

var (
	workerQueue chan chan WorkRequest
)

// StartDispatcher is providing the WorkerQueue
func StartDispatcher(workers int, installer *Installer) {
	// First, initialize the channel we are going to but the workers' work channels into.
	workerQueue = make(chan chan WorkRequest, workers)

	// Now, create all of our workers.
	for i := 0; i < workers; i++ {
		fmt.Println("Starting worker", i+1)
		worker := NewWorker(i+1, workerQueue, installer)
		worker.Start()
	}

	go func() {
		for {
			select {
			case work := <-workQueue:
				go func() {
					worker := <-workerQueue
					worker <- work
				}()
			}
		}
	}()
}
