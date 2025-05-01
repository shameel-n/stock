package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
)

type InstrumentInfo struct {
	CompanyName string `json:"companyName"`
	Symbol      string `json:"symbol"`
	Industry    string `json:"industry"`
	Group       string `json:"group"`
}

func main() {
	// Read NSE-500.csv
	csvFile, err := os.Open("data/NSE-500.csv")
	if err != nil {
		fmt.Printf("Error opening CSV file: %v\n", err)
		return
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	records, err := csvReader.ReadAll()
	if err != nil {
		fmt.Printf("Error reading CSV file: %v\n", err)
		return
	}

	// Read instruments.json
	jsonFile, err := os.ReadFile("data/instruments.json")
	if err != nil {
		fmt.Printf("Error reading JSON file: %v\n", err)
		return
	}

	var groupMap map[string]string
	if err := json.Unmarshal(jsonFile, &groupMap); err != nil {
		fmt.Printf("Error unmarshaling JSON: %v\n", err)
		return
	}

	// Create the final map
	instrumentInfo := make(map[string]InstrumentInfo)

	// Skip header row
	for _, record := range records[1:] {
		if len(record) < 5 {
			continue
		}

		companyName := record[0]
		industry := record[1]
		symbol := record[2]
		isin := record[4]

		// Create instrument key
		instrumentKey := fmt.Sprintf("NSE_EQ|%s", isin)

		// Get group from instruments.json
		group := groupMap[instrumentKey]

		// Add to map
		instrumentInfo[instrumentKey] = InstrumentInfo{
			CompanyName: companyName,
			Symbol:      symbol,
			Industry:    industry,
			Group:       group,
		}
	}

	// Write to output file
	outputFile, err := os.Create("data/instrument_info.json")
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer outputFile.Close()

	encoder := json.NewEncoder(outputFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(instrumentInfo); err != nil {
		fmt.Printf("Error writing JSON: %v\n", err)
		return
	}

	fmt.Println("Successfully generated instrument_info.json")
}
