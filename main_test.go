package main_test

import (
	"io"
	"os"
	"path"
	"testing"

	"github.com/owenrumney/go-sarif/v2/sarif"
)

func loadSarif(t *testing.T, dir string) *sarif.Report {
	t.Helper()

	f, err := os.Open(path.Join(dir, "sarif.json"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	report, err := sarif.FromBytes(content)
	if err != nil {
		t.Fatal(err)
	}

	return report
}
