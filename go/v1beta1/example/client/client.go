// Copyright 2019 The Grafeas Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	pb "github.com/grafeas/grafeas/proto/v1beta1/grafeas_go_proto"

	"log"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/webhook/event", handleWebhook)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w, "I'm healthy") })
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", 8080),
		Handler: mux,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("could not start http server...")
		}
	}()

	fmt.Println("listening for events")

}

func handleWebhook(w http.ResponseWriter, request *http.Request) {
	event := &Event{}
	if err := json.NewDecoder(request.Body).Decode(event); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("error reading webhook event")
		return
	}

	occ := &pb.Occurrence{}

}

type Event struct {
	Type     WebhookEvent `json:"type"`
	OccurAt  int64        `json:"occur_at"`
	Operator string       `json:"operator"`
	Data     *EventData   `json:"event_data"`
}
