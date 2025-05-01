package models

import "time"

// Candle represents a single candlestick data point
type Candle struct {
	Timestamp time.Time `json:"timestamp"`
	Open      float64   `json:"open"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Close     float64   `json:"close"`
	Volume    int64     `json:"volume"`
}

// HistoricalDataRequest represents the parameters for fetching historical data
type HistoricalDataRequest struct {
	Symbol     string
	Interval   string
	StartDate  time.Time
	EndDate    time.Time
	IsIntraday bool
}

// HistoricalDataResponse represents the response from the historical data API
type HistoricalDataResponse struct {
	Status string `json:"status"`
	Data   struct {
		Candles [][]interface{} `json:"candles"`
	} `json:"data"`
}

// TimeInterval represents the available time intervals for historical data
type TimeInterval string

const (
	OneMinute     TimeInterval = "1minute"
	FiveMinute    TimeInterval = "5minute"
	FifteenMinute TimeInterval = "15minute"
	ThirtyMinute  TimeInterval = "30minute"
	OneHour       TimeInterval = "1hour"
	OneDay        TimeInterval = "1day"
	OneWeek       TimeInterval = "1week"
	OneMonth      TimeInterval = "1month"
)
