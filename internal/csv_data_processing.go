package main

import (
	"fmt"
	"strconv"
	"strings"
)

type CsvFinancialData struct {
	name                      string
	currency                  string
	quantity                  float64
	changeTodayPercent        float64
	loanValue                 float64
	valueNOK                  float64
	returnOnInvestmentPercent float64
	returnOnInvestmentNOK     float64
}

type SectorMetrics struct {
	name string
	idealWeight float64
	currentWeight float64
	diffFromIdealWeight float64
	idealValue float64
	currentValue float64
	diffFromIdealValue float64
}

func getFundDataFromCsv(csvInfo *CsvInfo) []CsvFinancialData {
	var fundDataArray []CsvFinancialData

	for i, eachrecord := range csvInfo.content {
		if i == 0 { // skip header
			continue
		}

		quantity, _ := strconv.ParseFloat(strings.Replace(eachrecord[2], ",", ".", -1), 64)
		changePercent, _ := strconv.ParseFloat(strings.Replace(eachrecord[4], ",", ".", -1), 64)
		loanVal, _ := strconv.ParseFloat(strings.Replace(eachrecord[6], ",", ".", -1), 64)
		valNOK, _ := strconv.ParseFloat(strings.Replace(eachrecord[8], ",", ".", -1), 64)
		roiPercent, _ := strconv.ParseFloat(strings.Replace(eachrecord[9], ",", ".", -1), 64)
		roiNOK, _ := strconv.ParseFloat(strings.Replace(eachrecord[10], ",", ".", -1), 64)

		fundData := CsvFinancialData{
			name:                      eachrecord[0],
			currency:                  eachrecord[1],
			quantity:                  quantity,
			changeTodayPercent:        changePercent,
			loanValue:                 loanVal,
			valueNOK:                  valNOK,
			returnOnInvestmentPercent: roiPercent,
			returnOnInvestmentNOK:     roiNOK,
		}

		fundDataArray = append(fundDataArray, fundData)
	}

	return fundDataArray
}

func getStockDataFromCsv(csvInfo *CsvInfo) []CsvFinancialData {
	var stockDataArray []CsvFinancialData

	for i, eachrecord := range csvInfo.content {
		if i == 0 { // skip header
			continue
		}

		// Parse quantity, replace comma with dot for parsing
		quantity, _ := strconv.ParseFloat(strings.Replace(eachrecord[2], ",", ".", -1), 64)
		changePercent, _ := strconv.ParseFloat(strings.Replace(eachrecord[4], ",", ".", -1), 64)
		loanVal, _ := strconv.ParseFloat(strings.Replace(eachrecord[6], ",", ".", -1), 64)
		valNOK, _ := strconv.ParseFloat(strings.Replace(eachrecord[8], ",", ".", -1), 64)
		roiPercent, _ := strconv.ParseFloat(strings.Replace(eachrecord[9], ",", ".", -1), 64)
		roiNOK, _ := strconv.ParseFloat(strings.Replace(eachrecord[10], ",", ".", -1), 64)

		stockData := CsvFinancialData{
			name:                      eachrecord[0],
			currency:                  eachrecord[1],
			quantity:                  quantity,
			changeTodayPercent:        changePercent,
			loanValue:                 loanVal,
			valueNOK:                  valNOK,
			returnOnInvestmentPercent: roiPercent,
			returnOnInvestmentNOK:     roiNOK,
		}

		stockDataArray = append(stockDataArray, stockData)
	}

	return stockDataArray
}

func getFinancialData() []CsvFinancialData {
	fileNames := [2]string{"fondslister_konto-24260812_30.3.2026.csv",
		"aksjelister_konto-24260812_30.3.2026.csv"}

	var financialData []CsvFinancialData
	for _, fileName := range fileNames {
		csvInfo := getCsvInfo(fileName)
		switch csvInfo.financialType {
		case "aksje":
			financialData = append(financialData, getStockDataFromCsv(csvInfo)...)
		case "fond":
			financialData = append(financialData, getFundDataFromCsv(csvInfo)...)
		}
	}

	return financialData
}

func main() {
	financialData := getFinancialData()

	config, err := ReadConfig("config.yml")
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
	}

	sectorHoldings := getSectorHoldingValue(financialData, *config)

	getSectorMetrics(sectorHoldings, *config)
}

func getSectorMetrics(sectorHoldings map[string]float64, config Config) []SectorMetrics {
	var sectorMetrics []SectorMetrics
	sumOfAllHoldings := getSumOfAllHoldings(sectorHoldings)
	for sector, holdingValue := range sectorHoldings {
		idealWeight := config.Sectors[sector].Weight
		currentWeight := (holdingValue / sumOfAllHoldings) * 100
		diffWeight := idealWeight - currentWeight
		idealValue := (idealWeight / 100) * sumOfAllHoldings
		diffVal := idealValue - holdingValue

		sectorMetric := SectorMetrics{
			name: sector,
			idealWeight: idealWeight,
			currentWeight: currentWeight,
			diffFromIdealWeight: diffWeight,
			idealValue: idealValue,
			currentValue: holdingValue,
			diffFromIdealValue: diffVal,
		}

		sectorMetrics = append(sectorMetrics, sectorMetric)
	}
	return sectorMetrics
}

func getSumOfAllHoldings(sectorHoldings map[string]float64) float64 {
	var total float64
	for _, value := range sectorHoldings {
		total += value
	}
	return total
}

func checkIfNameMatchesFundOrStock(dataName string, keyword string) bool {
	if strings.Contains(strings.ToLower(dataName), strings.ToLower(keyword)) {
		return true
	}
	
	return false
}

func getSectorHoldingValue(financialData []CsvFinancialData, config Config) map[string]float64 {
	sectorHoldings := make(map[string]float64)
	for name := range config.Sectors {
		sectorHoldings[name] = 0
	}

	for name, sector := range config.Sectors {
		securities := append(sector.Funds, sector.Stocks...)
		for _, security := range securities {
			for _, data := range financialData {
				if checkIfNameMatchesFundOrStock(data.name, security) {
					sectorHoldings[name] += data.valueNOK
				}
			}
		}
	}

	return sectorHoldings
}
