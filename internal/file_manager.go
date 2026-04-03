package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/joho/godotenv"
)

func getNewestCsvFinancials() []string {
	path := getPathToCsvFinancials()

	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("Error: could not open path: %s", path)
	}

	stockListPattern := `^aksjelister_konto-\d+_(\d{1,2}\.\d{1,2}\.\d{4})( \(\d+\))?\.csv$`
	fundListPattern := `^fondslister_konto-\d+_(\d{1,2}\.\d{1,2}\.\d{4})( \(\d+\))?\.csv$`

	// loop over the content of the directory
	stockDataFilePath := ""
	stockDataFileNewestDate := time.Time{}

	fundDataFilePath := ""
	fundDataFileNewestDate := time.Time{}
	for _, entry := range entries {
		// skip subdirectories
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()

		// check if file matches the regexp pattern for stocks
		if checkIfFilenameMatchesPattern(stockListPattern, fileName) {
			time := extractDateFromFilename(stockListPattern, fileName)

			// if the file date is greater than the previous checked file, this file is currently the newest
			if checkIfFilenameDateIsNewest(time, stockDataFileNewestDate) {
				stockDataFileNewestDate = time
				stockDataFilePath = path + "/" + fileName
			}
			continue
		}

		// check if file matches the regexp pattern for funds
		if checkIfFilenameMatchesPattern(fundListPattern, fileName) {
			time := extractDateFromFilename(fundListPattern, fileName)

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

	return csvFinancials
}

func checkIfFilenameDateIsNewest(fileNameDate time.Time, currentNewestDate time.Time) bool {
	if fileNameDate.After(currentNewestDate) {
		return true
	}
	return false
}

func extractDateFromFilename(pattern string, fileName string) time.Time {
	re := getRegexp(pattern)
	matches := re.FindStringSubmatch(fileName)

	dateStr := matches[1]

	// Parse to time.Time
	t, _ := time.Parse("2.1.2006", dateStr)

	return t
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

func getPathToCsvFinancials() string {
	// Loads .env file from the current directory
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv("CSV_FINANCIALS_PATH")
}
