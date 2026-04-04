package main

import (
	core "Portfolio_go/internal"
	"fmt"
)

func main() {
	financialData := core.GetFinancialData()

	config, err := core.ReadConfig("sectors.yml")
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
	}

	sectorHoldings := core.GetSectorHoldingValue(financialData, *config)

	depositAmount := getDepositAmount()

	sectorMetrics := core.GetSectorMetrics(sectorHoldings, *config, depositAmount)

	core.PrintSectorMetrics(sectorMetrics)
}

func getDepositAmount() int {
	var deposit int

	// takes input value for name
	fmt.Print("Enter deposit amount: ")
	fmt.Scan(&deposit)

	return deposit
}
