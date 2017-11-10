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

// NewWorker creates, and returns a new Worker object. Its only argument
// is a channel that the worker can add itself to whenever it is done its
// work.
func NewWorker(id int, queue Queue) Worker {
	// Create, and return the worker.
	worker := Worker{
		ID:       id,
		Work:     make(WorkQueue),
		Worker:   queue.Worker,
		QuitChan: make(chan bool)}

	return worker
}

// Start is starting the worker by starting a goroutine.
func (w *Worker) Start() {
	go func() {
		for {
			// Add ourselves into the worker queue.
			w.Worker <- w.Work

			select {
			case work := <-w.Work:
				switch work.(type) {
				case *Install:
					// TODO: error handling
					err := install(work.(*Install))
					if err != nil {
						fmt.Println(err)
					}
				case *Uninstall:
					// TODO: error handling
					err := uninstall(work.(*Uninstall))
					if err != nil {
						fmt.Println(err)
					}
				default:
					break
				}
			case <-w.QuitChan:
				return
			}
		}
	}()
}

// Stop tells the worker to stop listening for work requests.
// Note that the worker will only stop *after* it has finished its work.
func (w *Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}
