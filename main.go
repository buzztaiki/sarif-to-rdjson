package main

// https://docs.oasis-open.org/sarif/sarif/v2.1.0/os/sarif-v2.1.0-os.html

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/owenrumney/go-sarif/v2/sarif"
	"github.com/reviewdog/reviewdog/proto/rdf"
	"google.golang.org/protobuf/encoding/protojson"
)

func or[T any](x *T, y T) T {
	if x == nil {
		return y
	}
	return *x
}

func findRuleByID(rules []*sarif.ReportingDescriptor, id string) *sarif.ReportingDescriptor {
	// https://docs.oasis-open.org/sarif/sarif/v2.1.0/os/sarif-v2.1.0-os.html#_Toc34317866

	if id == "" {
		return nil
	}

	for _, rule := range rules {
		if rule == nil {
			continue
		}

		if rule.ID == id {
			return rule
		}
		if strings.HasPrefix(rule.ID, id) && rule.ID[len(id)] == '/' {
			return rule
		}
	}

	return nil
}

func findRuleByIndex(rules []*sarif.ReportingDescriptor, index uint) *sarif.ReportingDescriptor {
	if len(rules) <= int(index) {
		return nil
	}

	if rule := rules[index]; rule != nil {
		return rule
	}

	return nil
}

func findRule(rules []*sarif.ReportingDescriptor, index *uint, id *string) *sarif.ReportingDescriptor {
	if index != nil {
		if rule := findRuleByIndex(rules, *index); rule != nil {
			return rule
		}
	}

	if id != nil {
		if rule := findRuleByID(rules, *id); rule != nil {
			return rule
		}
	}

	return nil
}

func findRuleFromResult(rules []*sarif.ReportingDescriptor, res *sarif.Result) *sarif.ReportingDescriptor {
	// result.ruleId
	//   https://docs.oasis-open.org/sarif/sarif/v2.1.0/os/sarif-v2.1.0-os.html#_Toc34317643
	// result.ruleIndex
	//   https://docs.oasis-open.org/sarif/sarif/v2.1.0/os/sarif-v2.1.0-os.html#_Toc34317644
	// result.rule
	//   https://docs.oasis-open.org/sarif/sarif/v2.1.0/os/sarif-v2.1.0-os.html#_Toc34317645
	// reportingDescriptorReference
	//   https://docs.oasis-open.org/sarif/sarif/v2.1.0/os/sarif-v2.1.0-os.html#_Toc34317862

	if res == nil {
		return nil
	}

	// TODO: sarif-web-component doesn't seems to use res.Rule, do the same?
	if res.Rule != nil {
		if rule := findRule(rules, res.Rule.Index, res.Rule.Id); rule != nil {
			return rule
		}
	}

	if rule := findRule(rules, res.RuleIndex, res.RuleID); rule != nil {
		return rule
	}

	if res.Rule != nil && res.Rule.Id != nil {
		return &sarif.ReportingDescriptor{ID: *res.Rule.Id}
	}

	if res.RuleID != nil {
		return &sarif.ReportingDescriptor{ID: *res.RuleID}
	}

	return nil
}

func rdfSeverity(res *sarif.Result, rule *sarif.ReportingDescriptor) rdf.Severity {
	// https://docs.oasis-open.org/sarif/sarif/v2.1.0/os/sarif-v2.1.0-os.html#_Toc34317648
	// https://docs.oasis-open.org/sarif/sarif/v2.1.0/os/sarif-v2.1.0-os.html#_Ref508894469

	kind := or(res.Kind, "fail")
	level := or(res.Level, "")
	if level == "" && rule != nil && rule.DefaultConfiguration != nil {
		level = rule.DefaultConfiguration.Level
	}

	if level == "" {
		if kind != "fail" {
			return rdf.Severity_UNKNOWN_SEVERITY
		}
		return rdf.Severity_WARNING
	}

	switch level {
	case "warning":
		return rdf.Severity_WARNING
	case "error":
		return rdf.Severity_ERROR
	case "note":
		return rdf.Severity_INFO
	default:
		return rdf.Severity_UNKNOWN_SEVERITY
	}
}

func rdfLocation(loc *sarif.Location) *rdf.Location {
	// https://docs.oasis-open.org/sarif/sarif/v2.1.0/os/sarif-v2.1.0-os.html#_Toc34317670

	if loc == nil || loc.PhysicalLocation == nil {
		return nil
	}

	start := &rdf.Position{
		Line:   int32(or(loc.PhysicalLocation.Region.StartLine, 0)),
		Column: int32(or(loc.PhysicalLocation.Region.StartColumn, 0)),
	}

	var end *rdf.Position
	if loc.PhysicalLocation.Region.EndLine != nil {
		end = &rdf.Position{
			Line:   int32(*loc.PhysicalLocation.Region.EndLine),
			Column: int32(or(loc.PhysicalLocation.Region.EndColumn, 0)),
		}
	}

	return &rdf.Location{
		Path:  or(loc.PhysicalLocation.ArtifactLocation.URI, ""),
		Range: &rdf.Range{Start: start, End: end},
	}
}

func rdfCode(rule *sarif.ReportingDescriptor) *rdf.Code {
	if rule == nil {
		return &rdf.Code{}
	}

	return &rdf.Code{
		Value: rule.ID,
		Url:   or(rule.HelpURI, ""),
	}
}

func rdfMessage(res *sarif.Result, rule *sarif.ReportingDescriptor) string {
	// https://docs.oasis-open.org/sarif/sarif/v2.1.0/os/sarif-v2.1.0-os.html#_Toc34317459
	// TODO: improve it
	// - read spec
	// - append rule's description?
	// - apppend rule's help.text? e.g. ansible-lint
	return or(res.Message.Text, "")
}

func SarifToRdf(report *sarif.Report) *rdf.DiagnosticResult {
	diags := make([]*rdf.Diagnostic, 0)
	var source rdf.Source

	for _, run := range report.Runs {
		for _, res := range run.Results {
			if res == nil {
				continue
			}

			// TODO: What if there are multiple runs?
			source.Name = run.Tool.Driver.Name
			source.Url = or(run.Tool.Driver.InformationURI, "")

			rule := findRuleFromResult(run.Tool.Driver.Rules, res)
			for _, loc := range res.Locations {
				diag := rdf.Diagnostic{
					Message:  rdfMessage(res, rule),
					Severity: rdfSeverity(res, rule),
					Location: rdfLocation(loc),
					Code:     rdfCode(rule),
				}
				diags = append(diags, &diag)
			}
		}
	}

	return &rdf.DiagnosticResult{Source: &source, Diagnostics: diags}
}

func abort(message string) {
	fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}

func abortIfError(message string, err error) {
	if err != nil {
		abort(fmt.Sprintf("%s: %s", message, err))
	}
}

func main() {
	content, err := io.ReadAll(os.Stdin)
	abortIfError("failed to read sarif", err)

	report, err := sarif.FromBytes(content)
	abortIfError("failed to parse sarif", err)

	out, err := protojson.Marshal(SarifToRdf(report))
	if err != nil {
		abortIfError("failed to marshal result", err)
	}
	os.Stdout.Write(out)
}
