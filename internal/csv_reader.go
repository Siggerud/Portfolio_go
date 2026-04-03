package main

import (
	//"fmt"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

type CsvInfo struct {
	content       [][]string
	financialType string
}

func get_csv_file(fileName string) *os.File {
	// open the file
	//TODO: denne må hentes på en annen måte
	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal("Error while reading the file", err)
	}

	// close the file
	//defer file.Close()

	return file
}

func get_csv_content(file *os.File) [][]string {
	// decode possible UTF-16 with BOM to UTF-8
	utf16Decoder := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewDecoder()
	decodedReader := transform.NewReader(file, utf16Decoder)

	// create a csv reader
	reader := csv.NewReader(decodedReader)

	// set tab as delimiter
	reader.Comma = '\t'

	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading records", err)
	}

	return records
}

func getCsvInfo(fileName string) *CsvInfo {
	// get csv file
	file := get_csv_file(fileName)

	defer file.Close()

	content := get_csv_content(file)

	var financialType string
	if strings.Contains(fileName, "fondslister") {
		financialType = "fond"
	} else if strings.Contains(fileName, "aksjelister") {
		financialType = "aksje"
	}

	return &CsvInfo{
		content:       content,
		financialType: financialType,
	}
}
