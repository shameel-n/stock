package main

import (
	"fmt"
	"log"
	"time"

	"github.com/shameel-n/stock/pkg/strategy"
	"github.com/shameel-n/stock/pkg/upstox/client"
	"github.com/shameel-n/stock/pkg/upstox/config"
	"github.com/shameel-n/stock/pkg/upstox/models"
)

func ExampleVWAPStrategy() {
	// Initialize Upstox client
	cfg := config.DefaultConfig()
	cfg.APIKey = "your-api-key-here" // Replace with your actual API key
	upstoxClient := client.NewClient(cfg)

	// Create stock configuration with dynamic limits
	stockConfig := strategy.StockConfig{
		Symbol:              "NSE_EQ|INE619B01017",
		UpperDeviationLimit: 1.2, // 1.2%
		LowerDeviationLimit: 1.0, // 1.0%
	}

	// Initialize VWAP strategy for the stock
	vwapStrategy := strategy.NewVWAPStrategy(stockConfig)

	// Fetch historical data
	req := models.HistoricalDataRequest{
		Symbol:     stockConfig.Symbol,
		Interval:   string(models.OneMinute),
		StartDate:  time.Now().AddDate(0, 0, -2),
		EndDate:    time.Now().AddDate(0, 0, -1),
		IsIntraday: false,
	}

	response, err := upstoxClient.GetHistoricalData(req)
	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}

	// Convert response to candles
	var candles []models.Candle
	for _, rawCandle := range response.Data.Candles {
		candle := models.Candle{
			Timestamp: time.Now(), // Replace with actual timestamp from response
			Open:      rawCandle[1].(float64),
			High:      rawCandle[2].(float64),
			Low:       rawCandle[3].(float64),
			Close:     rawCandle[4].(float64),
			Volume:    int64(rawCandle[5].(float64)),
		}
		candles = append(candles, candle)
	}

	// Use the latest close price as current price
	if len(candles) > 0 {
		latestCandle := candles[len(candles)-1]
		currentPrice := latestCandle.Close

		// Analyze the stock
		signal := vwapStrategy.AnalyzeStock(candles)
		if signal != nil {
			fmt.Println(signal)
		} else {
			fmt.Printf("No trading signal for %s at price %.2f\n", stockConfig.Symbol, currentPrice)
		}

		// Example of updating limits dynamically
		newConfig := strategy.StockConfig{
			Symbol:              stockConfig.Symbol,
			UpperDeviationLimit: 1.5, // Increased to 1.5%
			LowerDeviationLimit: 1.2, // Increased to 1.2%
		}
		vwapStrategy.UpdateStockConfig(newConfig)

		// Re-analyze with new limits
		signal = vwapStrategy.AnalyzeStock(candles)
		if signal != nil {
			fmt.Println("Signal with updated limits:", signal)
		}
	}
}
