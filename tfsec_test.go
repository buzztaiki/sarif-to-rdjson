package main_test

import (
	"testing"

	"github.com/buzztaiki/sarif-to-rdjson"
	"github.com/google/go-cmp/cmp"
	"github.com/reviewdog/reviewdog/proto/rdf"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestTfsec(t *testing.T) {
	report := loadSarif(t, "testcases/tfsec")

	want := rdf.DiagnosticResult{
		Source: &rdf.Source{
			Name: "defsec",
			Url:  "https://github.com/aquasecurity/defsec",
		},
		Diagnostics: []*rdf.Diagnostic{
			{
				Code: &rdf.Code{
					Value: "aws-ec2-enable-at-rest-encryption",
					Url:   "https://aquasecurity.github.io/tfsec/latest/checks/aws/ec2/enable-at-rest-encryption/",
				},
				Severity: rdf.Severity_ERROR,
				Message:  `Root block device is not encrypted.`,
				Location: &rdf.Location{
					Path: "main.tf",
					Range: &rdf.Range{
						Start: &rdf.Position{Line: 15},
						End:   &rdf.Position{Line: 22},
					},
				},
			},
			{
				Code: &rdf.Code{
					Value: "aws-ec2-enforce-http-token-imds",
					Url:   "https://aquasecurity.github.io/tfsec/latest/checks/aws/ec2/enforce-http-token-imds/",
				},
				Severity: rdf.Severity_ERROR,
				Message:  "Instance does not require IMDS access to require a token",
				Location: &rdf.Location{
					Path: "main.tf",
					Range: &rdf.Range{
						Start: &rdf.Position{Line: 15},
						End:   &rdf.Position{Line: 22},
					},
				},
			},
		},
	}

	if diff := cmp.Diff(&want, main.SarifToRdf(report), protocmp.Transform()); diff != "" {
		t.Errorf("tflint/simple mismatch (-want +got):\n%s", diff)
	}
}
