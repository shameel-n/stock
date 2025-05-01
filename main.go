package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/shameel-n/stock/pkg/instrument"
	"github.com/shameel-n/stock/pkg/processor"
	"github.com/shameel-n/stock/pkg/upstox/client"
	"github.com/shameel-n/stock/pkg/upstox/config"
)

// ObserveStocksLive starts a live monitoring session for stock analysis
func ObserveStocksLive() error {
	// Initialize managers
	infoManager := instrument.NewManager()
	limitManager := instrument.NewLimitManager()

	// Load data
	if err := infoManager.LoadFromFile("data/instrument_info.json"); err != nil {
		return fmt.Errorf("error loading instrument info: %v", err)
	}

	if err := limitManager.LoadFromFile("data/limit.json"); err != nil {
		return fmt.Errorf("error loading limits: %v", err)
	}

	// Initialize Upstox client and analyzer
	cfg := config.DefaultConfig()
	upstoxClient := client.NewClient(cfg)
	analyzer := processor.NewAnalyzer(upstoxClient, infoManager, limitManager)
	analyzer.SetWorkerCount(runtime.NumCPU() * 2)

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Printf("Starting analysis with %d workers...\n", runtime.NumCPU()*2)
	fmt.Println("Press Ctrl+C to stop")

	// Run first analysis immediately
	runAnalysis(analyzer, infoManager)

	// Create ticker for periodic analysis
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	// Main processing loop
	for {
		select {
		case <-ticker.C:
			runAnalysis(analyzer, infoManager)
		case <-sigChan:
			fmt.Println("\nReceived shutdown signal. Cleaning up...")
			return nil
		}
	}
}

// runAnalysis performs one cycle of stock analysis
func runAnalysis(analyzer *processor.Analyzer, infoManager *instrument.Manager) {
	startTime := time.Now()
	fmt.Printf("\n[%s] Starting new analysis cycle...\n", startTime.Format("15:04:05"))

	signals := analyzer.ProcessAllStocks()

	fmt.Printf("\n[%s] Found %d trading signals:\n", time.Now().Format("15:04:05"), len(signals))
	for instrumentKey, signal := range signals {
		info, _ := infoManager.GetInfo(instrumentKey)
		fmt.Printf("\nSignal generated for %s (%s):\n", info.Symbol, info.CompanyName)
		fmt.Printf("  Type: %s\n", signal.Type)
		fmt.Printf("  Price: %.2f\n", signal.Price)
		fmt.Printf("  VWAP: %.2f\n", signal.VWAP)
		fmt.Printf("  Deviation: %.2f%%\n", signal.Deviation)
		fmt.Printf("  Time: %s\n", signal.Time.Format("15:04:05"))
	}

	fmt.Printf("\n[%s] Analysis cycle completed in %v\n",
		time.Now().Format("15:04:05"),
		time.Since(startTime))
}

func main() {
	if err := ObserveStocksLive(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
