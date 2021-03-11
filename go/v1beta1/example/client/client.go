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

func retrieveSonarServerInfo(sonarUrl string, staticanalysis *pb.Occurrence_StaticAnalysis) {

}

func retrieveOverallMeasures(sonarUrl string, staticanalysis *pb.Occurrence_StaticAnalysis) {

}

func retrieveAllIssues(sonarUrl string, staticanalysis *pb.Occurrence_StaticAnalysis) {
	// TODO implement pagination
}

func handleWebhook(w http.ResponseWriter, request *http.Request) {
	event := &Event{}
	if err := json.NewDecoder(request.Body).Decode(event); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("error reading webhook event")
		return
	}
	body, _ := json.Marshal(event)
	log.Printf(string(body))

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

// Sonar stuff-----------------------------------------------------------
type Event struct {
	TaskID      string            `json:"taskid"`
	Status      string            `json:"status"`
	AnalyzedAt  string            `json:"analyzedat"`
	GitCommit   string            `json:"revision"`
	Project     *Project          `json:"project"`
	QualityGate *QualityGate      `json:"qualityGate"`
	Branch      *Branch           `json:"branch"`
	Properties  map[string]string `json:"properties"`
}

// Branch is...
type Branch struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	IsMain bool   `json:"isMain"`
	URL    string `json:"url"`
}

// Project is
type Project struct {
	Key  string `json:"key"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

// QualityGate is...
type QualityGate struct {
	Conditions []*Condition `json:"conditions"`
	Name       string       `json:"name"`
	Status     string       `json:"status"`
}

// Condition is...
type Condition struct {
	ErrorThreshold string `json:"errorThreshold"`
	Metric         string `json:"metric"`
	OnLeakPeriod   bool   `json:"onLeakPeriod"`
	Operator       string `json:"operator"`
	Status         string `json:"status"`
}

// Overall Measures
type OverallMeasures struct {
	Component struct {
		ID        string `json:"id"`
		Key       string `json:"key"`
		Name      string `json:"name"`
		Qualifier string `json:"qualifier"`
		Measures  []struct {
			Metric    string `json:"metric"`
			Value     string `json:"value,omitempty"`
			BestValue bool   `json:"bestValue,omitempty"`
			Periods   []struct {
				Index     int    `json:"index"`
				Value     string `json:"value"`
				BestValue bool   `json:"bestValue"`
			} `json:"periods,omitempty"`
			Period struct {
				Index     int    `json:"index"`
				Value     string `json:"value"`
				BestValue bool   `json:"bestValue"`
			} `json:"period,omitempty"`
		} `json:"measures"`
	} `json:"component"`
	Metrics []struct {
		Key                   string `json:"key"`
		Name                  string `json:"name"`
		Description           string `json:"description"`
		Domain                string `json:"domain"`
		Type                  string `json:"type"`
		HigherValuesAreBetter bool   `json:"higherValuesAreBetter,omitempty"`
		Qualitative           bool   `json:"qualitative"`
		Hidden                bool   `json:"hidden"`
		Custom                bool   `json:"custom"`
		BestValue             string `json:"bestValue,omitempty"`
		DecimalScale          int    `json:"decimalScale,omitempty"`
		WorstValue            string `json:"worstValue,omitempty"`
	} `json:"metrics"`
	Periods []struct {
		Index int    `json:"index"`
		Mode  string `json:"mode"`
		Date  string `json:"date"`
	} `json:"periods"`
}

// Individual Issues
type Issues struct {
	Total  int `json:"total"`
	P      int `json:"p"`
	Ps     int `json:"ps"`
	Paging struct {
		PageIndex int `json:"pageIndex"`
		PageSize  int `json:"pageSize"`
		Total     int `json:"total"`
	} `json:"paging"`
	EffortTotal int `json:"effortTotal"`
	DebtTotal   int `json:"debtTotal"`
	Issues      []struct {
		Key       string `json:"key"`
		Rule      string `json:"rule"`
		Severity  string `json:"severity"`
		Component string `json:"component"`
		Project   string `json:"project"`
		Line      int    `json:"line"`
		Hash      string `json:"hash"`
		TextRange struct {
			StartLine   int `json:"startLine"`
			EndLine     int `json:"endLine"`
			StartOffset int `json:"startOffset"`
			EndOffset   int `json:"endOffset"`
		} `json:"textRange"`
		Flows        []interface{} `json:"flows"`
		Status       string        `json:"status"`
		Message      string        `json:"message"`
		Effort       string        `json:"effort"`
		Debt         string        `json:"debt"`
		Tags         []string      `json:"tags"`
		Transitions  []interface{} `json:"transitions"`
		Actions      []interface{} `json:"actions"`
		Comments     []interface{} `json:"comments"`
		CreationDate string        `json:"creationDate"`
		UpdateDate   string        `json:"updateDate"`
		Type         string        `json:"type"`
		Organization string        `json:"organization"`
		Scope        string        `json:"scope"`
	} `json:"issues"`
	Components []struct {
		Organization string `json:"organization"`
		Key          string `json:"key"`
		UUID         string `json:"uuid"`
		Enabled      bool   `json:"enabled"`
		Qualifier    string `json:"qualifier"`
		Name         string `json:"name"`
		LongName     string `json:"longName"`
		Path         string `json:"path,omitempty"`
	} `json:"components"`
	Rules []struct {
		Key      string `json:"key"`
		Name     string `json:"name"`
		Lang     string `json:"lang"`
		Status   string `json:"status"`
		LangName string `json:"langName"`
	} `json:"rules"`
	Users     []interface{} `json:"users"`
	Languages []struct {
		Key  string `json:"key"`
		Name string `json:"name"`
	} `json:"languages"`
	Facets []struct {
		Property string `json:"property"`
		Values   []struct {
			Val   string `json:"val"`
			Count int    `json:"count"`
		} `json:"values"`
	} `json:"facets"`
}
