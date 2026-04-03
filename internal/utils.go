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
		fmt.Println("Error:", err)
	}

	stockListPattern := `^aksjelister_konto-\d+_(\d{1,2}\.\d{1,2}\.\d{4})( \(\d+\))?\.csv$`
	fundListPattern := `^fondslister^_konto-\d+_(\d{1,2}\.\d{1,2}\.\d{4})( \(\d+\))?\.csv$`

	reStock := regexp.MustCompile(stockListPattern)
	reFund := regexp.MustCompile(fundListPattern)

	// loop over the content of the directory
	stockDataFile := ""
	stockDataFileNewestDate := time.Time{}

	fundDataFile := ""
	fundDataFileNewestDate := time.Time{}
	for _, entry := range entries {
		// skip subdirectories
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()

		// check if file matches the regexp pattern for stocks
		if reStock.MatchString(fileName) {
			matches := reStock.FindStringSubmatch(fileName)

			dateStr := matches[1]
			fmt.Println("Extracted date:", dateStr)

			// Parse to time.Time
			t, _ := time.Parse("2.1.2006", dateStr)

			// if the file date is greater than the previous checked file, this file is currently the newest
			if t.After(stockDataFileNewestDate) {
				stockDataFileNewestDate = t
				stockDataFile = fileName
			}
			continue
		}

		// check if file matches the regexp pattern for funds
		if reFund.MatchString(fileName) {
			matches := reFund.FindStringSubmatch(fileName)

			dateStr := matches[1]
			fmt.Println("Extracted date:", dateStr)

			// Parse to time.Time
			t, _ := time.Parse("2.1.2006", dateStr)

			// if the file date is greater than the previous checked file, this file is currently the newest
			if t.After(fundDataFileNewestDate) {
				fundDataFileNewestDate = t
				fundDataFile = fileName
			}
		}
	}

	// add csv files for stocks and funds if they exist
	csvFinancials := []string{}
	if fundDataFile != "" {
		csvFinancials = append(csvFinancials, fundDataFile)
	}

	if stockDataFile != "" {
		csvFinancials = append(csvFinancials, stockDataFile)
	}

	return csvFinancials
}

func getPathToCsvFinancials() string {
	// Loads .env file from the current directory
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv("PATH")
}

func main() {
	csvs := getNewestCsvFinancials()
	fmt.Println(csvs)
}
