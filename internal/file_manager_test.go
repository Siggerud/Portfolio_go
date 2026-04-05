package core

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckIfFilenameMatchesPatternMatches(t *testing.T) {
	pattern := `^[A-Za-z]+_\d{2}\.\d{2}\.\d{4}\.txt$`
	fileName := "financials_01.12.2024.txt"

	assert.True(t, checkIfFilenameMatchesPattern(pattern, fileName), "Pattern should match filename.")
}

func TestCheckIfFilenameMatchesPatternNotMatches(t *testing.T) {
	pattern := `^[A-Za-z]+_\d{2}\.\d{2}\.\d{4}\.txt$`
	fileName := "financials_01.12.2024.csv"

	assert.False(t, checkIfFilenameMatchesPattern(pattern, fileName), "Pattern should match filename.")
}

func TestGetPathToCsvFinancialsLoadsValueFromDotEnv(t *testing.T) {
	tempDir := t.TempDir()

	// Get current working directory
	oldWd, err := os.Getwd()
	assert.NoError(t, err)

	// change working directory to testing tmp directory
	err = os.Chdir(tempDir)
	assert.NoError(t, err)

	// change working directory back on ending of test
	t.Cleanup(func() {
		_ = os.Chdir(oldWd)
	})

	// create a .env file in the testing tmp directory, and write the value for the environment variable
	expectedPath := "/tmp/csv-financials"
	err = os.WriteFile(filepath.Join(tempDir, ".env"), []byte("CSV_FINANCIALS_PATH="+expectedPath+"\n"), 0o644)
	assert.NoError(t, err)

	actualPath := getPathToCsvFinancials()

	assert.Equal(t, expectedPath, actualPath)
}
