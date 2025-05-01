package instrument

import (
	"encoding/json"
	"fmt"
	"os"
)

// Info represents the basic information about a stock
type Info struct {
	CompanyName string `json:"companyName"`
	Symbol      string `json:"symbol"`
	Industry    string `json:"industry"`
	Group       string `json:"group"`
}

// Manager handles instrument information operations
type Manager struct {
	instruments map[string]Info
}

// NewManager creates a new instrument manager
func NewManager() *Manager {
	return &Manager{
		instruments: make(map[string]Info),
	}
}

// LoadFromFile loads instrument information from a JSON file
func (m *Manager) LoadFromFile(filePath string) error {
	jsonFile, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading instrument info: %v", err)
	}

	if err := json.Unmarshal(jsonFile, &m.instruments); err != nil {
		return fmt.Errorf("error unmarshaling instrument info: %v", err)
	}

	return nil
}

// GetInfo returns the information for a given instrument key
func (m *Manager) GetInfo(instrumentKey string) (Info, bool) {
	info, exists := m.instruments[instrumentKey]
	return info, exists
}

// GetAllInstruments returns all instrument information
func (m *Manager) GetAllInstruments() map[string]Info {
	return m.instruments
}
