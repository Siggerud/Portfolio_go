package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/olekukonko/tablewriter"
)

func printSectorMetrics(sectorMetrics []*SectorMetrics) {
	sortSectorMetrics(sectorMetrics)

	table := tablewriter.NewWriter(os.Stdout)

	table.Header([]string{
		"Name",
		"Ideal Wt",
		"Current Wt",
		"Diff Wt",
		"Ideal Val",
		"Current Val",
		"Diff Val",
	})

	for _, s := range sectorMetrics {
		table.Append([]string{
			s.name,
			fmt.Sprintf("%.2f", s.idealWeight),
			fmt.Sprintf("%.2f", s.currentWeight),
			fmt.Sprintf("%+.2f", s.diffFromIdealWeight),
			fmt.Sprintf("%.2f", s.idealValue),
			fmt.Sprintf("%.2f", s.currentValue),
			fmt.Sprintf("%+.2f", s.diffFromIdealValue),
		})
	}

	table.Render()
}

func sortSectorMetrics(sectorMetrics []*SectorMetrics) {
	sort.Slice(sectorMetrics, func(i, j int) bool {
		return sectorMetrics[i].diffFromIdealValue < sectorMetrics[j].diffFromIdealValue
	})
}
