package core

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/joho/godotenv"
)

type configProvider interface {
	readConfig(path string) (*Config, error)
}

type configReader struct{}

func (c configReader) readConfig(path string) (*Config, error) {
	return ReadConfigYaml(path)
}

func getRegexPatterns(c configProvider) (map[string]string, error) {
	config, err := c.readConfig("../../configs/config.yml")
	if err != nil {
		return map[string]string{}, fmt.Errorf("Error reading config: %v\n", err)
	}

	// read the regex patterns from the config.yml file
	stockListPattern := config.Patterns.StockFilenamePattern
	fundListPattern := config.Patterns.FundFilenamepattern

	return map[string]string{
		"stock": stockListPattern,
		"fund":  fundListPattern,
	}, nil
}

func getAllCsvFinancials() ([]string, error) {
	path, err := getPathToCsvFinancials()

	if err != nil {
		return []string{}, err
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return []string{}, errors.New("Error: could not open path: " + path)
	}

	regexPatterns, err := getRegexPatterns(configReader{})
	if err != nil {
		return []string{}, fmt.Errorf("Error: Something went wrong retrieving regex patterns: %v", err)
	}

	// read the regex patterns
	stockListPattern := regexPatterns["stock"]
	fundListPattern := regexPatterns["fund"]

	csvFinancialFiles := []string{}
	// loop over the content of the directory
	for _, entry := range entries {
		// skip subdirectories
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()
		// check if file matches the regexp pattern for stocks
		if checkIfFilenameMatchesPattern(stockListPattern, fileName) {
			csvFinancialFiles = append(csvFinancialFiles, fileName)
			continue
		}

		// check if file matches the regexp pattern for funds
		if checkIfFilenameMatchesPattern(fundListPattern, fileName) {
			csvFinancialFiles = append(csvFinancialFiles, fileName)
		}
	}

	return csvFinancialFiles, nil
}

func getNewestCsvFinancials() ([]string, error) {
	csvFinancialFiles, err := getAllCsvFinancials()
	if err != nil {
		return []string{}, fmt.Errorf("Error: Something went wrong retrieving csv financials: %v", err)
	}

	path, err := getPathToCsvFinancials()

	regexPatterns, err := getRegexPatterns(configReader{})
	if err != nil {
		return []string{}, fmt.Errorf("Error: Something went wrong retrieving regex patterns: %v", err)
	}

	// read the regex patterns
	stockListPattern := regexPatterns["stock"]
	fundListPattern := regexPatterns["fund"]

	stockDataFilePath := ""
	stockDataFileNewestDate := time.Time{}

	fundDataFilePath := ""
	fundDataFileNewestDate := time.Time{}
	for _, fileName := range csvFinancialFiles {
		// check if file matches the regexp pattern for stocks
		if checkIfFilenameMatchesPattern(stockListPattern, fileName) {
			time, err := extractDateFromFilename(stockListPattern, fileName)
			if err != nil {
				return []string{}, err
			}

			// if the file date is greater than the previous checked file, this file is currently the newest
			if checkIfFilenameDateIsNewest(time, stockDataFileNewestDate) {
				stockDataFileNewestDate = time
				stockDataFilePath = path + "/" + fileName
			}
			continue
		}

		// check if file matches the regexp pattern for funds
		if checkIfFilenameMatchesPattern(fundListPattern, fileName) {
			time, err := extractDateFromFilename(fundListPattern, fileName)

			if err != nil {
				return []string{}, err
			}
			// if the file date is greater than the previous checked file, this file is currently the newest
			if checkIfFilenameDateIsNewest(time, fundDataFileNewestDate) {
				fundDataFileNewestDate = time
				fundDataFilePath = path + "/" + fileName
			}
		}
	}

	// add csv files for stocks and funds if they exist
	csvFinancials := []string{}
	if fundDataFilePath != "" {
		csvFinancials = append(csvFinancials, fundDataFilePath)
	}

	if stockDataFilePath != "" {
		csvFinancials = append(csvFinancials, stockDataFilePath)
	}

	return csvFinancials, nil
}

func checkIfFilenameDateIsNewest(fileNameDate time.Time, currentNewestDate time.Time) bool {
	if fileNameDate.After(currentNewestDate) {
		return true
	}
	return false
}

func extractDateFromFilename(pattern string, fileName string) (time.Time, error) {
	re := getRegexp(pattern)
	matches := re.FindStringSubmatch(fileName)
	if len(matches) == 1 {
		return time.Time{}, errors.New("No date could be extracted from pattern")
	} else if len(matches) > 2 {
		return time.Time{}, errors.New("More than one match group in pattern. Inconclusive date.")
	}

	dateStr := matches[1]

	// Parse to time.Time
	t, _ := time.Parse("2.1.2006", dateStr)

	return t, nil
}

func checkIfFilenameMatchesPattern(pattern string, fileName string) bool {
	re := getRegexp(pattern)
	if re.MatchString(fileName) {
		return true
	}
	return false
}

func getRegexp(pattern string) *regexp.Regexp {
	return regexp.MustCompile(pattern)
}

func getPathToCsvFinancials() (string, error) {
	// Loads .env file from the current directory
	err := godotenv.Load("../../.env")
	if err != nil {
		return "", errors.New("Error loading .env file")
	}

	return os.Getenv("CSV_FINANCIALS_PATH"), nil
}
