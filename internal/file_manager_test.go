package core

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

	// Create nested directories: tempDir/level1/level2
	level1 := filepath.Join(tempDir, "level1")
	level2 := filepath.Join(level1, "level2")

	err := os.MkdirAll(level2, 0o755)
	assert.NoError(t, err)

	// Get current working directory
	oldWd, err := os.Getwd()
	assert.NoError(t, err)

	// Change working directory to the deepest level
	err = os.Chdir(level2)
	assert.NoError(t, err)

	// Restore working directory after test
	t.Cleanup(func() {
		_ = os.Chdir(oldWd)
	})

	// Create .env file TWO levels above (in tempDir)
	expectedPath := "/tmp/csvs/csv-financials"
	envPath := filepath.Join(tempDir, ".env")

	err = os.WriteFile(envPath, []byte("CSV_FINANCIALS_PATH="+expectedPath+"\n"), 0o644)
	assert.NoError(t, err)

	actualPath, err := getPathToCsvFinancials()

	assert.Equal(t, expectedPath, actualPath)
	assert.NoError(t, err)
}

func TestGetPathToCsvFinancialsThrowsErros(t *testing.T) {
	tempDir := t.TempDir()

	// Get current working directory
	oldWd, err := os.Getwd()
	assert.NoError(t, err)

	// change working directory to testing tmp directory which will not contain an .env file
	err = os.Chdir(tempDir)
	assert.NoError(t, err)

	// change working directory back on ending of test
	t.Cleanup(func() {
		_ = os.Chdir(oldWd)
	})

	expectedError := errors.New("Error loading .env file")
	_, actualError := getPathToCsvFinancials()

	assert.Equal(t, expectedError, actualError)
}

func TestExtractDateFromFilenameParsesDate(t *testing.T) {
	pattern := `^aksjelister_konto-\d+_(\d{1,2}\.\d{1,2}\.\d{4})(?: \(\d+\))?\.csv$`
	fileName := "aksjelister_konto-123_5.4.2026 (1).csv"

	expectedDate := time.Date(2026, time.April, 5, 0, 0, 0, 0, time.UTC)
	actualDate, err := extractDateFromFilename(pattern, fileName)

	assert.Equal(t, expectedDate, actualDate)
	assert.NoError(t, err)
}

func TestExtractDateFromFilenameParsesDateThrowsErrorWithNoMatch(t *testing.T) {
	pattern := `^aksjelister_konto-\d+_\d{1,2}\.\d{1,2}\.\d{4}(?: \(\d+\))?\.csv$`
	fileName := "aksjelister_konto-123_5.4.2026 (1).csv"

	expectedError := errors.New("No date could be extracted from pattern")
	_, actualError := extractDateFromFilename(pattern, fileName)

	assert.Equal(t, expectedError, actualError)
}

func TestExtractDateFromFilenameParsesDateThrowsErrorWithMoreThanOneMatch(t *testing.T) {
	pattern := `^(.*?)_(konto-\d+_[\d.]+).*\.csv$`
	fileName := "aksjelister_konto-123_5.4.2026 (1).csv"

	expectedError := errors.New("More than one match group in pattern. Inconclusive date.")
	_, actualError := extractDateFromFilename(pattern, fileName)

	assert.Equal(t, expectedError, actualError)
}

func TestCheckIfFilenameDateIsNewest(t *testing.T) {
	// Define test cases
	tests := []struct {
		name              string
		fileNameDate      time.Time
		currentNewestDate time.Time
		expected          bool
	}{

		{
			name:              "fileNameDate is newer",
			fileNameDate:      time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC),
			currentNewestDate: time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC),
			expected:          true,
		},
		{
			name:              "fileNameDate is older",
			fileNameDate:      time.Date(2026, 4, 29, 0, 0, 0, 0, time.UTC),
			currentNewestDate: time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC),
			expected:          false,
		},
		{
			name:              "fileNameDate is the same",
			fileNameDate:      time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC),
			currentNewestDate: time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC),
			expected:          false,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkIfFilenameDateIsNewest(tt.fileNameDate, tt.currentNewestDate)
			assert.Equal(t, tt.expected, result, "Unexpected result for test case: %s", tt.name)
		})
	}
}

// Create a mock ConfigReader to mock the readConfig method
type mockConfigReader struct{ mock.Mock }

func newMockConfigReader() *mockConfigReader { return &mockConfigReader{} }

func (m *mockConfigReader) readConfig(path string) (*Config, error) {
	args := m.Called(path)

	// Retrieve the *Config value from the mock arguments
	config, _ := args.Get(0).(*Config)

	return config, nil
}

func TestDivByRand(t *testing.T) {
	// get our mock object
	m := newMockConfigReader()
	// specify our return value when mock object is called with specified path
	m.On("readConfig", "../../configs/config.yml").Return(&Config{
		Patterns: Pattern{
			StockFilenamePattern: "stock_\\d{4}-\\d{2}-\\d{2}\\.csv",
			FundFilenamepattern:  "fund_\\d{4}-\\d{2}-\\d{2}\\.csv",
		},
	}, nil)

	patterns, _ := getRegexPatterns(m)

	assert.Equal(t, patterns["stock"], "stock_\\d{4}-\\d{2}-\\d{2}\\.csv")
	assert.Equal(t, patterns["fund"], "fund_\\d{4}-\\d{2}-\\d{2}\\.csv")
}
