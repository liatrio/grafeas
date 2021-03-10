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
	"github.com/grafeas/grafeas/proto/v1beta1/source_go_proto"
	"github.com/grafeas/grafeas/proto/v1beta1/static_analysis_go_proto"
	"google.golang.org/protobuf/types/known/timestamppb"

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

	log.Fatal(server.ListenAndServe())

	fmt.Println("listening for events")

}

func handleWebhook(w http.ResponseWriter, request *http.Request) {
	// event := &Event{}
	// if err := json.NewDecoder(request.Body).Decode(event); err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	fmt.Println("error reading webhook event")
	// 	return
	// }

	staticanalysis := &pb.Occurrence_StaticAnalysis{
		StaticAnalysis: &static_analysis_go_proto.Details{
			AnalysisResults: &static_analysis_go_proto.StaticAnalysis{
				Tool:        "",
				ToolVersion: "",
				ToolConfig:  "",
				Summary: &static_analysis_go_proto.Stats{
					Complexity: &static_analysis_go_proto.Complexity{
						Cyclomatic: 0,
						Cognitive:  0,
						Findings:   nil,
					},
					Duplication: &static_analysis_go_proto.Duplication{
						Blocks:       0,
						Files:        0,
						Lines:        0,
						LinesDensity: 0.0,
						Findings:     nil,
					},
					Maintainability: &static_analysis_go_proto.Maintainability{
						CodeSmells:     0,
						SqaleRating:    0,
						SqaleIndex:     0,
						SqaleDebtRatio: 0.0,
						Findings:       nil,
					},
					Reliability: &static_analysis_go_proto.Reliability{
						Bugs:              0,
						Rating:            0,
						RemediationEffort: 0,
						Findings:          nil,
					},
					Security: &static_analysis_go_proto.Security{
						Vulnerabilities:           0,
						SecurityRating:            0,
						SecurityRemediationEffort: 0,
						SecurityReviewRating:      0,
						Findings:                  nil,
					},
					CodeSize: &static_analysis_go_proto.CodeSize{
						Classes:             0,
						CommentLines:        0,
						CommentLinesDensity: 0.0,
						Directories:         0,
						Files:               0,
						Lines:               0,
						Ncloc:               0,
						Functions:           0,
						Statements:          0,
						Findings:            nil,
					},
					Issues: &static_analysis_go_proto.Issues{
						Total:          0,
						Blocker:        0,
						Critical:       0,
						Major:          0,
						Minor:          0,
						Info:           0,
						FalsePositives: 0,
						Open:           0,
						Confirmed:      0,
						Reopened:       0,
						Findings:       nil,
					},
				},
				Context: &source_go_proto.SourceContext{
					Context: nil,
					Labels: map[string]string{
						"": "",
					},
				},
				StartTime: &timestamppb.Timestamp{
					Seconds: 0,
					Nanos:   0,
				},
				EndTime: &timestamppb.Timestamp{
					Seconds: 0,
					Nanos:   0,
				},
			},
		},
	}

	occ := &pb.Occurrence{
		Name: "",
		Resource: &pb.Resource{
			Name: "",
			Uri:  "",
		},
		NoteName:    "",
		Kind:        0,
		Remediation: "",
		CreateTime: &timestamppb.Timestamp{
			Seconds: 0,
			Nanos:   0,
		},
		UpdateTime: &timestamppb.Timestamp{
			Seconds: 0,
			Nanos:   0,
		},
		Details: staticanalysis,
	}
	out, _ := json.Marshal(occ)
	log.Printf(string(out))
}
