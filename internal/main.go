package main

import "fmt"

func main() {
	financialData := getFinancialData()

	config, err := ReadConfig("sectors.yml")
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
	}

	sectorHoldings := getSectorHoldingValue(financialData, *config)

	depositAmount := getDepositAmount()

	sectorMetrics := getSectorMetrics(sectorHoldings, *config, depositAmount)

	printSectorMetrics(sectorMetrics)
}

func getDepositAmount() int {
	var deposit int

	// takes input value for name
	fmt.Print("Enter deposit amount: ")
	fmt.Scan(&deposit)

	return deposit
}
