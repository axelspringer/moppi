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

import "fmt"

// New is providing the Queue
func New(workers int) WorkQueue {
	return mustNew(workers)
}

// mustNew wraps the creation of a new queue
func mustNew(workers int) WorkQueue {
	var queue Queue
	queue.Worker = make(chan chan interface{}, workers)
	queue.Work = make(chan interface{}, 100)

	// Now, create all of our workers.
	for i := 0; i < workers; i++ {
		// TODO: remove or substitude
		fmt.Println("Starting worker", i+1)
		worker := NewWorker(i+1, queue)
		worker.Start()
	}

	go func() {
		for {
			select {
			case work := <-queue.Work:
				go func() {
					worker := <-queue.Worker
					worker <- work
				}()
			}
		}
	}()

	return queue.Work
}
