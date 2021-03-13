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
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

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

func retrieveSonarServerInfo(sonarURL string, staticanalysis *pb.Occurrence_StaticAnalysis) {
	url := fmt.Sprintf("%s/api/server/version", sonarURL)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	staticanalysis.StaticAnalysis.AnalysisResults.ToolVersion = string(body)
	staticanalysis.StaticAnalysis.AnalysisResults.Tool = "Sonarqube"
}

func retrieveOverallMeasures(sonarURL string, staticanalysis *pb.Occurrence_StaticAnalysis) {

	// parm := url.Values{}
	// parm.Add("additionalFields", "metrics,periods")
	// parm.Add("component", "app")
	// parm.Add("metricKeys", "sqale_debt_ratio,duplicated_files,duplicated_lines,alert_status,functions,quality_gate_details,bugs,new_bugs,reliability_rating,new_reliability_rating,vulnerabilities,new_vulnerabilities,security_rating,new_security_rating,security_hotspots,new_security_hotspots,security_hotspots_reviewed,new_security_hotspots_reviewed,security_review_rating,new_security_review_rating,code_smells,new_code_smells,sqale_rating,new_maintainability_rating,sqale_index,new_technical_debt,coverage,new_coverage,lines_to_cover,new_lines_to_cover,tests,duplicated_lines_density,new_duplicated_lines_density,duplicated_blocks,ncloc,ncloc_language_distribution,projects,lines,new_lines,complexity,cognitive_complexity,reliability_remediation_effort,security_remediation_effort,classes,comment_lines,comment_lines_density,directories,files,statements,new_code_smells,new_technical_debt,new_sqale_debt_ratio")
	// https://docs.sonarqube.org/latest/user-guide/metric-definitions/
	url := fmt.Sprintf("%s/api/measures/component?additionalFields=metrics,periods&component=app&metricKeys=sqale_debt_ratio,duplicated_files,duplicated_lines,alert_status,functions,quality_gate_details,bugs,new_bugs,reliability_rating,new_reliability_rating,vulnerabilities,new_vulnerabilities,security_rating,new_security_rating,security_hotspots,new_security_hotspots,security_hotspots_reviewed,new_security_hotspots_reviewed,security_review_rating,new_security_review_rating,code_smells,new_code_smells,sqale_rating,new_maintainability_rating,sqale_index,new_technical_debt,coverage,new_coverage,lines_to_cover,new_lines_to_cover,tests,duplicated_liness_density,new_duplicated_lines_density,duplicated_blocks,ncloc,ncloc_language_distribution,projects,lines,new_lines,complexity,cognitive_complexity,reliability_remediation_effort,security_remediation_effort,classes,comment_lines,comment_lines_density,directories,files,statements,new_code_smells,new_technical_debt,new_sqale_debt_ratio,new_reliability_remediation_effort,new_security_remediation_effort,violations,new_violations,blocker_violations,critical_violations,major_violations,minor_violations,info_violations,new_blocker_violations,new_critical_violations,new_major_violations,new_minor_violations,new_info_violations,false_positive_issues,open_issues,confirmed_issues,reopened_issues", sonarURL)
	// url := fmt.Sprintf("%s/api/measures/component", sonarURL)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	// req, err := http.NewRequest(method, url, strings.NewReader(parm.Encode()))
	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	overallMeasures := &OverallMeasures{}
	if err := json.NewDecoder(res.Body).Decode(overallMeasures); err != nil {
		fmt.Println("error reading response from sonarqube")
		return
	}

	for i := range overallMeasures.Component.Measures {
		switch overallMeasures.Component.Measures[i].Metric {
		case "complexity":
			staticanalysis.StaticAnalysis.AnalysisResults.Summary.Complexity.Cyclomatic = stringToUint32(overallMeasures.Component.Measures[i].Value)
		case "cognitive_complexity":
			staticanalysis.StaticAnalysis.AnalysisResults.Summary.Complexity.Cognitive = stringToUint32(overallMeasures.Component.Measures[i].Value)
			//-----------------------------------------------------------------------
		case "duplicated_lines":
			staticanalysis.StaticAnalysis.AnalysisResults.Summary.Duplication.Lines = stringToUint32(overallMeasures.Component.Measures[i].Value)
		case "duplicated_files":
			staticanalysis.StaticAnalysis.AnalysisResults.Summary.Duplication.Files = stringToUint32(overallMeasures.Component.Measures[i].Value)
		case "duplicated_lines_density":
			staticanalysis.StaticAnalysis.AnalysisResults.Summary.Duplication.LinesDensity = stringToFloat(overallMeasures.Component.Measures[i].Value)
		case "duplicated_blocks":
			staticanalysis.StaticAnalysis.AnalysisResults.Summary.Duplication.Blocks = stringToUint32(overallMeasures.Component.Measures[i].Value)
			// ---------------------------------------------------------------------
		case "code_smells":
			staticanalysis.StaticAnalysis.AnalysisResults.Summary.Maintainability.CodeSmells = stringToUint64(overallMeasures.Component.Measures[i].Value)
		case "sqale_rating":
			staticanalysis.StaticAnalysis.AnalysisResults.Summary.Maintainability.SqaleRating = static_analysis_go_proto.Rating(stringToInt(overallMeasures.Component.Measures[i].Value))
		case "sqale_index":
			staticanalysis.StaticAnalysis.AnalysisResults.Summary.Maintainability.SqaleIndex = stringToUint32(overallMeasures.Component.Measures[i].Value)
		case "sqale_debt_ratio":
			staticanalysis.StaticAnalysis.AnalysisResults.Summary.Maintainability.SqaleDebtRatio = stringToFloat(overallMeasures.Component.Measures[i].Value)
			// --------------------------------------------------------
		case "bugs":
			staticanalysis.StaticAnalysis.AnalysisResults.Summary.Reliability.Bugs = stringToUint32(overallMeasures.Component.Measures[i].Value)
		case "reliability_rating":
			staticanalysis.StaticAnalysis.AnalysisResults.Summary.Reliability.Rating = static_analysis_go_proto.Rating(stringToInt(overallMeasures.Component.Measures[i].Value))
		case "reliability_remediation_effort":
			staticanalysis.StaticAnalysis.AnalysisResults.Summary.Reliability.RemediationEffort = stringToUint32(overallMeasures.Component.Measures[i].Value)
			// -------------------------------------------------------
		case "vulnerabilities":
			staticanalysis.StaticAnalysis.AnalysisResults.Summary.Security.Vulnerabilities = stringToUint32(overallMeasures.Component.Measures[i].Value)
		case "security_rating":
			staticanalysis.StaticAnalysis.AnalysisResults.Summary.Security.SecurityRating = static_analysis_go_proto.Rating(stringToInt(overallMeasures.Component.Measures[i].Value))
		case "security_remediation_effort":
			staticanalysis.StaticAnalysis.AnalysisResults.Summary.Security.SecurityRemediationEffort = stringToUint32(overallMeasures.Component.Measures[i].Value)
		case "security_review_rating":
			staticanalysis.StaticAnalysis.AnalysisResults.Summary.Security.SecurityReviewRating = static_analysis_go_proto.Rating(stringToInt(overallMeasures.Component.Measures[i].Value))
			// -------------------------------------------------------------
		case "classes":
			staticanalysis.StaticAnalysis.AnalysisResults.Summary.CodeSize.Classes = stringToUint32(overallMeasures.Component.Measures[i].Value)
		case "comment_lines":
			staticanalysis.StaticAnalysis.AnalysisResults.Summary.CodeSize.CommentLines = stringToUint32(overallMeasures.Component.Measures[i].Value)
		case "comment_lines_density":
			staticanalysis.StaticAnalysis.AnalysisResults.Summary.CodeSize.CommentLinesDensity = stringToFloat(overallMeasures.Component.Measures[i].Value)
		case "directories":
			staticanalysis.StaticAnalysis.AnalysisResults.Summary.CodeSize.Directories = stringToUint32(overallMeasures.Component.Measures[i].Value)
		case "files":
			staticanalysis.StaticAnalysis.AnalysisResults.Summary.CodeSize.Files = stringToUint64(overallMeasures.Component.Measures[i].Value)
		case "lines":
			staticanalysis.StaticAnalysis.AnalysisResults.Summary.CodeSize.Lines = stringToUint64(overallMeasures.Component.Measures[i].Value)
		case "functions":
			staticanalysis.StaticAnalysis.AnalysisResults.Summary.CodeSize.Functions = stringToUint64(overallMeasures.Component.Measures[i].Value)
		case "statements":
			staticanalysis.StaticAnalysis.AnalysisResults.Summary.CodeSize.Statements = stringToUint64(overallMeasures.Component.Measures[i].Value)
		}

	}
}

func stringToUint32(val string) uint32 {
	result, err := strconv.ParseUint(val, 10, 32)
	if err != nil {
		log.Printf("Failed to convert string to uint32")
	}
	return uint32(result)
}

func stringToUint64(val string) uint64 {
	result, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		log.Printf("Failed to convert string to uint64")
	}
	return result
}

func stringToFloat(val string) float32 {
	result, err := strconv.ParseFloat(val, 32)
	if err != nil {
		log.Printf("Failed to convert string to float32")
	}
	return float32(result)
}

func stringToInt(val string) int {
	if strings.Contains(val, ".") {
		return int(stringToFloat(val))
	}
	result, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		log.Printf("Failed to convert string to int")
	}
	return int(result)
}

func retrieveAllIssues(sonarURL string, staticanalysis *pb.Occurrence_StaticAnalysis) {
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

	retrieveSonarServerInfo(event.ServerURL, staticanalysis)
	retrieveOverallMeasures(event.ServerURL, staticanalysis)
	retrieveAllIssues(event.ServerURL, staticanalysis)

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
	ServerURL   string            `json:"serverUrl"`
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
	Component Component `json:"component"`
	Metrics   []Metrics `json:"metrics"`
	Periods   []Periods `json:"periods"`
}

type Periods struct {
	Index int    `json:"index"`
	Mode  string `json:"mode"`
	Date  string `json:"date"`
}

type Metrics struct {
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
}
type Component struct {
	ID        string     `json:"id"`
	Key       string     `json:"key"`
	Name      string     `json:"name"`
	Qualifier string     `json:"qualifier"`
	Measures  []Measures `json:"measures"`
}

type Measures struct {
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
