package instrument

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/shameel-n/stock/pkg/constant"
)

// Limit represents the trading limits for a stock
type Limit struct {
	UpperLimit float64 `json:"UpperLimit"`
	LowerLimit float64 `json:"LowerLimit"`
	Symbol     string  `json:"Symbol"`
}

// LimitManager handles stock limit operations
type LimitManager struct {
	limits map[string]Limit
}

// NewLimitManager creates a new limit manager
func NewLimitManager() *LimitManager {
	return &LimitManager{
		limits: make(map[string]Limit),
	}
}

// LoadFromFile loads stock limits from a JSON file
func (m *LimitManager) LoadFromFile(filePath string) error {
	jsonFile, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading limit file: %v", err)
	}

	if err := json.Unmarshal(jsonFile, &m.limits); err != nil {
		return fmt.Errorf("error unmarshaling limit info: %v", err)
	}

	return nil
}

// GetLimit returns the limits for a given instrument key
func (m *LimitManager) GetLimit(instrumentKey string) (Limit, bool) {
	limit, exists := m.limits[instrumentKey]

	limit.UpperLimit = limit.UpperLimit * constant.UpperLimitMultiplier
	limit.LowerLimit = limit.LowerLimit * constant.LowerLimitMultiplier

	return limit, exists
}
