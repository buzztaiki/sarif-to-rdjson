package main_test

import (
	"testing"

	"github.com/buzztaiki/sarif-to-rdjson"
	"github.com/google/go-cmp/cmp"
	"github.com/reviewdog/reviewdog/proto/rdf"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestCodeQL(t *testing.T) {
	report := loadSarif(t, "testcases/codeql")

	want := rdf.DiagnosticResult{
		Source: &rdf.Source{
			Name: "CodeQL",
		},
		Diagnostics: []*rdf.Diagnostic{
			{
				Code: &rdf.Code{
					Value: "go/reflected-xss",
				},
				Severity: rdf.Severity_ERROR,
				Message:  "Cross-site scripting vulnerability due to [user-provided value](1).",
				Location: &rdf.Location{
					Path: "codeql.go",
					Range: &rdf.Range{
						Start: &rdf.Position{Line: 13, Column: 43},
						End:   &rdf.Position{Column: 51},
					},
				},
			},
			{
				Code: &rdf.Code{
					Value: "go/sql-injection",
				},
				Severity: rdf.Severity_ERROR,
				Message:  "This query depends on a [user-provided value](1).",
				Location: &rdf.Location{
					Path: "codeql.go",
					Range: &rdf.Range{
						Start: &rdf.Position{Line: 21, Column: 11},
						End:   &rdf.Position{Column: 12},
					},
				},
			},
		},
	}

	if diff := cmp.Diff(&want, main.SarifToRdf(report), protocmp.Transform()); diff != "" {
		t.Errorf("tflint/simple mismatch (-want +got):\n%s", diff)
	}
}
