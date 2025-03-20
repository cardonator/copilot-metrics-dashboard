package services

import (
	"fmt"
	"os"
	"path/filepath"
)

// loadTestData loads test data from the testdata directory
func loadTestData(filename string) ([]byte, error) {
	// Try to find the testdata directory
	// First, check in the current directory
	if _, err := os.Stat("testdata"); err == nil {
		return os.ReadFile(filepath.Join("testdata", filename))
	}

	// Then check in parent directory
	if _, err := os.Stat("../testdata"); err == nil {
		return os.ReadFile(filepath.Join("../testdata", filename))
	}

	// Then check in the project root (assuming we're in cmd/dataingestion or similar)
	if _, err := os.Stat("../../testdata"); err == nil {
		return os.ReadFile(filepath.Join("../../testdata", filename))
	}

	return nil, fmt.Errorf("test data file not found: %s", filename)
}
