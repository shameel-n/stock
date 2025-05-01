package processor

import (
	"fmt"
	"sync"
	"time"

	"github.com/shameel-n/stock/pkg/instrument"
	"github.com/shameel-n/stock/pkg/strategy"
	"github.com/shameel-n/stock/pkg/upstox/client"
	"github.com/shameel-n/stock/pkg/upstox/models"
)

// Analyzer handles stock data processing and analysis
type Analyzer struct {
	upstoxClient *client.Client
	infoManager  *instrument.Manager
	limitManager *instrument.LimitManager
	workerCount  int
}

// NewAnalyzer creates a new stock analyzer
func NewAnalyzer(upstoxClient *client.Client, infoManager *instrument.Manager, limitManager *instrument.LimitManager) *Analyzer {
	return &Analyzer{
		upstoxClient: upstoxClient,
		infoManager:  infoManager,
		limitManager: limitManager,
		workerCount:  10, // Default number of workers
	}
}

// SetWorkerCount sets the number of concurrent workers
func (a *Analyzer) SetWorkerCount(count int) {
	if count > 0 {
		a.workerCount = count
	}
}

// ProcessStock processes a single stock and returns any trading signals
func (a *Analyzer) ProcessStock(instrumentKey string) (*strategy.Signal, error) {
	// Get stock info
	info, exists := a.infoManager.GetInfo(instrumentKey)
	if !exists {
		return nil, fmt.Errorf("no info found for %s", instrumentKey)
	}

	// Get stock limits
	limits, exists := a.limitManager.GetLimit(instrumentKey)
	if !exists {
		return nil, fmt.Errorf("no limits found for %s", info.Symbol)
	}

	// Create VWAP strategy with stock-specific limits
	vwapStrategy := strategy.NewVWAPStrategy(strategy.StockConfig{
		Symbol:              instrumentKey,
		UpperDeviationLimit: limits.UpperLimit,
		LowerDeviationLimit: limits.LowerLimit,
	})

	// Fetch historical data
	req := models.HistoricalDataRequest{
		Symbol:     instrumentKey,
		Interval:   string(models.OneMinute),
		StartDate:  time.Now().AddDate(0, 0, -2),
		EndDate:    time.Now().AddDate(0, 0, -1),
		IsIntraday: true,
	}

	response, err := a.upstoxClient.GetHistoricalData(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching data for %s: %v", info.Symbol, err)
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

	if len(candles) == 0 {
		return nil, fmt.Errorf("no candle data for %s", info.Symbol)
	}

	// Analyze the stock
	return vwapStrategy.AnalyzeStock(candles), nil
}

// ProcessAllStocks processes all stocks concurrently and returns a map of signals
func (a *Analyzer) ProcessAllStocks() map[string]*strategy.Signal {
	signals := make(map[string]*strategy.Signal)
	var mu sync.Mutex // Mutex to protect the signals map

	// Create channels for work distribution
	jobs := make(chan string, a.workerCount)
	results := make(chan struct {
		key    string
		signal *strategy.Signal
		err    error
	}, a.workerCount)

	// Start worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < a.workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for instrumentKey := range jobs {
				signal, err := a.ProcessStock(instrumentKey)
				results <- struct {
					key    string
					signal *strategy.Signal
					err    error
				}{instrumentKey, signal, err}
			}
		}()
	}

	// Start a goroutine to close results channel when all workers are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Send jobs to workers
	go func() {
		for instrumentKey := range a.infoManager.GetAllInstruments() {
			jobs <- instrumentKey
		}
		close(jobs)
	}()

	// Process results
	for result := range results {
		if result.err != nil {
			fmt.Printf("Error processing %s: %v\n", result.key, result.err)
			continue
		}

		if result.signal != nil {
			mu.Lock()
			signals[result.key] = result.signal
			mu.Unlock()
		}
	}

	return signals
}
