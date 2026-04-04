package core

import (
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
