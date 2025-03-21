package svcTag_test

import (
	"testing"

	"github.com/kigland/HPC-Scheduler/lib/svcTag"
)

func TestSvcTag(t *testing.T) {
	data := "KHS-kevin-test-csSAnAaE"
	svgT, err := svcTag.Parse(data)
	if err != nil {
		t.Fatalf("Failed to parse svcTag: %v", err)
	}
	if svgT.Owner != "kevin" {
		t.Fatalf("Owner is not kevin")
	}
	if svgT.Project != "test" {
		t.Fatalf("Project is not test")
	}
	if svgT.Rand != "csSAnAaE" {
		t.Fatalf("Rand is not csSAnAaE")
	}
}
