package main

import (
	"fmt"
	"log"
	"time"

	"github.com/shameel-n/stock/pkg/upstox/client"
	"github.com/shameel-n/stock/pkg/upstox/config"
	"github.com/shameel-n/stock/pkg/upstox/models"
)

func main() {
	// Create a new client with default configuration
	cfg := config.DefaultConfig()
	cfg.APIKey = "your-api-key-here" // Replace with your actual API key

	upstoxClient := client.NewClient(cfg)

	// Create a request for historical data
	req := models.HistoricalDataRequest{
		Symbol:     "NSE_EQ|INE619B01017", // Example symbol
		Interval:   string(models.OneMinute),
		StartDate:  time.Now().AddDate(0, 0, -2), // Yesterday
		EndDate:    time.Now().AddDate(0, 0, -1), // Today
		IsIntraday: false,
	}

	// Fetch historical data
	response, err := upstoxClient.GetHistoricalData(req)
	if err != nil {
		log.Fatalf("Error fetching historical data: %v", err)
	}

	// Print the response
	fmt.Printf("Status: %s\n", response.Status)
	fmt.Printf("Number of candles: %d\n", len(response.Data.Candles))
}
