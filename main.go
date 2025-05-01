package main

import (
	"fmt"
	"time"

	"github.com/shameel-n/stock/pkg/upstox/client"
	"github.com/shameel-n/stock/pkg/upstox/config"
	"github.com/shameel-n/stock/pkg/upstox/models"
)

func main() {
	cfg := config.DefaultConfig()
	upstoxClient := client.NewClient(cfg)

	req := models.HistoricalDataRequest{
		Symbol:     "NSE_EQ|INE619B01017",
		Interval:   string(models.OneMinute),
		StartDate:  time.Now().AddDate(0, 0, -2),
		EndDate:    time.Now().AddDate(0, 0, -1),
		IsIntraday: false,
	}

	response, err := upstoxClient.GetHistoricalData(req)
	if err != nil {
		fmt.Println("Error fetching historical data:", err)
		return
	}

	fmt.Println("Historical data fetched successfully!")
	fmt.Println("Response:", response)
}
