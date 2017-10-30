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

import (
	"log"
	"time"
)

// NewWorker creates, and returns a new Worker object. Its only argument
// is a channel that the worker can add itself to whenever it is done its
// work.
func NewWorker(id int, workerQueue chan chan WorkRequest, installer *Installer) Worker {
	// Create, and return the worker.
	worker := Worker{
		ID:          id,
		Work:        make(chan WorkRequest),
		WorkerQueue: workerQueue,
		Installer:   installer,
		QuitChan:    make(chan bool)}

	return worker
}

// Start is starting the worker by starting a goroutine.
func (w *Worker) Start() {
	go func() {
		for {
			// Add ourselves into the worker queue.
			w.WorkerQueue <- w.Work

			select {
			case work := <-w.Work:
				time.Sleep(work.Delay)

				// deploy marathon
				if work.Install.Marathon {
					if _, err := w.Installer.marathon.CreateApplication(&work.Marathon); err != nil {
						log.Printf("Failed to create application: %s, error: %s", work.Marathon.ID, err)
					} else {
						log.Printf("Created the application: %s", work.Marathon.ID)
					}
				}

				// deploy chronos
				if work.Install.Chronos {
					if ok, _, _ := w.Installer.chronos.Job.New(&work.Chronos); !ok {
						log.Printf("Failed to create job: %s, error: %s", work.Chronos.Name)
					} else {
						log.Printf("Created the application: %s", work.Chronos.Name)
					}
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
