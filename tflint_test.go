package main_test

import (
	"os/exec"
	"testing"

	"github.com/buzztaiki/sarif-to-rdjson"
	"github.com/google/go-cmp/cmp"
	"github.com/owenrumney/go-sarif/v2/sarif"
	"github.com/reviewdog/reviewdog/proto/rdf"
	"google.golang.org/protobuf/testing/protocmp"
)

func tflint(t *testing.T, dir string) *sarif.Report {
	t.Helper()
	cmd := exec.Command("tflint", "--format", "sarif", "--force")
	cmd.Dir = dir

	output, err := cmd.Output()
	if err != nil {
		t.Fatal(err, string(output))
	}

	report, err := sarif.FromBytes(output)
	if err != nil {
		t.Fatal(err)
	}
	return report
}

func TestTflintSimple(t *testing.T) {
	report := tflint(t, "testcases/tflint/simple/")

	want := rdf.DiagnosticResult{
		Source: &rdf.Source{
			Name: "tflint",
			Url:  "https://github.com/terraform-linters/tflint",
		},
		Diagnostics: []*rdf.Diagnostic{
			{
				Message: `"t100.micro" is an invalid value as instance_type`,
				Code: &rdf.Code{
					Value: "aws_instance_invalid_type",
				},
				Severity: rdf.Severity_ERROR,
				Location: &rdf.Location{
					Path: "main.tf",
					Range: &rdf.Range{
						Start: &rdf.Position{Line: 17, Column: 19},
						End:   &rdf.Position{Line: 17, Column: 31},
					},
				},
			},
			{
				Message: `Missing version constraint for provider "aws" in "required_providers"`,
				Code: &rdf.Code{
					Value: "terraform_required_providers",
					Url:   "https://github.com/terraform-linters/tflint-ruleset-terraform/blob/v0.2.2/docs/rules/terraform_required_providers.md",
				},
				Severity: rdf.Severity_WARNING,
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
