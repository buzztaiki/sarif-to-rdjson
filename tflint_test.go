package main_test

import (
	"testing"

	"github.com/buzztaiki/sarif-to-rdjson"
	"github.com/google/go-cmp/cmp"
	"github.com/reviewdog/reviewdog/proto/rdf"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestTflintSimple(t *testing.T) {
	report := loadSarif(t, "testcases/tflint/simple")

	want := rdf.DiagnosticResult{
		Source: &rdf.Source{
			Name: "tflint",
			Url:  "https://github.com/terraform-linters/tflint",
		},
		Diagnostics: []*rdf.Diagnostic{
			{
				Code: &rdf.Code{
					Value: "aws_instance_invalid_type",
				},
				Severity: rdf.Severity_ERROR,
				Message:  `"t100.micro" is an invalid value as instance_type`,
				Location: &rdf.Location{
					Path: "main.tf",
					Range: &rdf.Range{
						Start: &rdf.Position{Line: 17, Column: 19},
						End:   &rdf.Position{Line: 17, Column: 31},
					},
				},
			},
			{
				Code: &rdf.Code{
					Value: "terraform_required_providers",
					Url:   "https://github.com/terraform-linters/tflint-ruleset-terraform/blob/v0.2.2/docs/rules/terraform_required_providers.md",
				},
				Severity: rdf.Severity_WARNING,
				Message:  `Missing version constraint for provider "aws" in "required_providers"`,
				Location: &rdf.Location{
					Path: "main.tf",
					Range: &rdf.Range{
						Start: &rdf.Position{Line: 15, Column: 1},
						End:   &rdf.Position{Line: 15, Column: 37},
					},
				},
			},
		},
	}

	if diff := cmp.Diff(&want, main.SarifToRdf(report), protocmp.Transform()); diff != "" {
		t.Errorf("tflint/simple mismatch (-want +got):\n%s", diff)
	}
}

func TestTflintSyntaxError(t *testing.T) {
	report := loadSarif(t, "testcases/tflint/syntax-error")

	want := rdf.DiagnosticResult{
		Source: &rdf.Source{
			Name: "tflint-errors",
			Url:  "https://github.com/terraform-linters/tflint",
		},
		Diagnostics: []*rdf.Diagnostic{
			{
				Code: &rdf.Code{
					Value: "Unclosed configuration block",
				},
				Severity: rdf.Severity_ERROR,
				Message:  `There is no closing brace for this block before the end of the file. This may be caused by incorrect brace nesting elsewhere in this file.`,
				Location: &rdf.Location{
					Path: "main.tf",
					Range: &rdf.Range{
						Start: &rdf.Position{Line: 1, Column: 11},
						End:   &rdf.Position{Line: 1, Column: 12},
					},
				},
			},
		},
	}

	if diff := cmp.Diff(&want, main.SarifToRdf(report), protocmp.Transform()); diff != "" {
		t.Errorf("tflint/simple mismatch (-want +got):\n%s", diff)
	}
}
