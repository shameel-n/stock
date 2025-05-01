package strategy

import (
	"fmt"
	"time"

	"github.com/shameel-n/stock/pkg/upstox/models"
)

// StockConfig holds configuration for a single stock
type StockConfig struct {
	Symbol              string
	UpperDeviationLimit float64
	LowerDeviationLimit float64
}

// VWAPStrategy represents the VWAP-based trading strategy
type VWAPStrategy struct {
	stockConfig StockConfig
}

// Signal represents a trading signal
type Signal struct {
	Symbol    string
	Type      SignalType
	Price     float64
	VWAP      float64
	Deviation float64
	Time      time.Time
}

// SignalType represents the type of trading signal
type SignalType string

const (
	BuySignal  SignalType = "BUY"
	SellSignal SignalType = "SELL"
)

// NewVWAPStrategy creates a new VWAP strategy for a single stock
func NewVWAPStrategy(config StockConfig) *VWAPStrategy {
	return &VWAPStrategy{
		stockConfig: config,
	}
}

// UpdateStockConfig updates the configuration for the stock
func (s *VWAPStrategy) UpdateStockConfig(config StockConfig) {
	s.stockConfig = config
}

// CalculateVWAP calculates the Volume Weighted Average Price for a list of candles
func CalculateVWAP(candles []models.Candle) float64 {
	var cumulativePV float64 // Price * Volume
	var cumulativeVolume float64

	for _, candle := range candles {
		// Use typical price (High + Low + Close) / 3
		typicalPrice := (candle.High + candle.Low + candle.Close) / 3
		cumulativePV += typicalPrice * float64(candle.Volume)
		cumulativeVolume += float64(candle.Volume)
	}

	if cumulativeVolume == 0 {
		return 0
	}

	return cumulativePV / cumulativeVolume
}

// AnalyzeStock analyzes a single stock for VWAP deviation signals
func (s *VWAPStrategy) AnalyzeStock(candles []models.Candle) *Signal {
	if len(candles) == 0 {
		return nil
	}

	// Get current price from the last candle
	currentPrice := candles[len(candles)-1].Close

	vwap := CalculateVWAP(candles)
	if vwap == 0 {
		return nil
	}

	// Calculate deviation percentage
	deviation := ((currentPrice - vwap) / vwap) * 100

	// Check for trading signals
	if deviation >= s.stockConfig.UpperDeviationLimit {
		return &Signal{
			Symbol:    s.stockConfig.Symbol,
			Type:      SellSignal,
			Price:     currentPrice,
			VWAP:      vwap,
			Deviation: deviation,
			Time:      time.Now(),
		}
	} else if deviation <= -s.stockConfig.LowerDeviationLimit {
		return &Signal{
			Symbol:    s.stockConfig.Symbol,
			Type:      BuySignal,
			Price:     currentPrice,
			VWAP:      vwap,
			Deviation: deviation,
			Time:      time.Now(),
		}
	}

	return nil
}

// String returns a string representation of a Signal
func (s Signal) String() string {
	return fmt.Sprintf("[%s] %s Signal for %s at %.2f (VWAP: %.2f, Deviation: %.2f%%)",
		s.Time.Format("15:04:05"),
		s.Type,
		s.Symbol,
		s.Price,
		s.VWAP,
		s.Deviation,
	)
}
