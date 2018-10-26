// This script will perform black box tests on parseRecords

package main

import (
	"fmt"
	"testing"
)

func getTestData() (string, string) {
	// Returns paths to test data
	var dict, input strings
	wd, _ := os.Executable()
	fmt.Println(wd)
	return dict, input
}

func TestExtractDiagnosis(t *testing.T) {
	// Tests extractDiangosis output
	_, _ = getTestData()

}
