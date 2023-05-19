package main_test

import (
	"testing"

	"github.com/buzztaiki/sarif-to-rdjson"
	"github.com/google/go-cmp/cmp"
	"github.com/reviewdog/reviewdog/proto/rdf"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestAnsibleLint(t *testing.T) {
	report := loadSarif(t, "testcases/ansible-lint")

	want := rdf.DiagnosticResult{
		Source: &rdf.Source{
			Name: "ansible-lint",
			Url:  "https://github.com/ansible/ansible-lint",
		},
		Diagnostics: []*rdf.Diagnostic{
			{
				Code: &rdf.Code{
					Value: "name[play]",
					Url:   "https://ansible-lint.readthedocs.io/rules/name/",
				},
				Severity: rdf.Severity_ERROR,
				Message:  "All plays should be named.",
				Location: &rdf.Location{
					Path: "site.yaml",
					Range: &rdf.Range{
						Start: &rdf.Position{Line: 3},
					},
				},
			},
			{
				Code: &rdf.Code{
					Value: "name[missing]",
					Url:   "https://ansible-lint.readthedocs.io/rules/name/",
				},
				Severity: rdf.Severity_ERROR,
				Message:  "Task/Handler: command cmd=pwd",
				Location: &rdf.Location{
					Path: "site.yaml",
					Range: &rdf.Range{
						Start: &rdf.Position{Line: 5},
					},
				},
			},
			{
				Code: &rdf.Code{
					Value: "no-changed-when",
					Url:   "https://ansible-lint.readthedocs.io/rules/no-changed-when/",
				},
				Severity: rdf.Severity_ERROR,
				Message:  "Task/Handler: command cmd=pwd",
				Location: &rdf.Location{
					Path: "site.yaml",
					Range: &rdf.Range{
						Start: &rdf.Position{Line: 5},
					},
				},
			},
		},
	}

	if diff := cmp.Diff(&want, main.SarifToRdf(report), protocmp.Transform()); diff != "" {
		t.Errorf("tflint/simple mismatch (-want +got):\n%s", diff)
	}
}
