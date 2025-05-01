package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/shameel-n/stock/pkg/upstox/config"
	"github.com/shameel-n/stock/pkg/upstox/errors"
	"github.com/shameel-n/stock/pkg/upstox/models"
)

// Client represents the Upstox API client
type Client struct {
	config     *config.Config
	httpClient *http.Client
	cache      *cache
}

type cache struct {
	data  map[string]models.HistoricalDataResponse
	mutex sync.RWMutex
}

// NewClient creates a new Upstox client
func NewClient(cfg *config.Config) *Client {
	if cfg == nil {
		cfg = config.DefaultConfig()
	}

	return &Client{
		config: cfg,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
		cache: &cache{
			data: make(map[string]models.HistoricalDataResponse),
		},
	}
}

// GetHistoricalData fetches historical data for a symbol
func (c *Client) GetHistoricalData(req models.HistoricalDataRequest) (*models.HistoricalDataResponse, error) {
	cacheKey := fmt.Sprintf("%s_%s_%s_%s_%t",
		req.Symbol,
		req.Interval,
		req.StartDate.Format("2006-01-02"),
		req.EndDate.Format("2006-01-02"),
		req.IsIntraday,
	)

	// Check cache first
	if cached, ok := c.getFromCache(cacheKey); ok {
		return &cached, nil
	}

	// Build URL
	url := c.buildHistoricalDataURL(req)

	// Create request
	httpReq, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.New(errors.ErrInvalidRequest, "failed to create request", err)
	}

	// Add headers
	httpReq.Header.Add("Accept", "application/json")
	httpReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))

	// Execute request with retries
	var response *http.Response
	for i := 0; i < c.config.MaxRetries; i++ {
		response, err = c.httpClient.Do(httpReq)
		if err == nil {
			break
		}
		time.Sleep(c.config.RetryDelay)
	}
	if err != nil {
		return nil, errors.New(errors.ErrInternalServer, "failed to execute request", err)
	}
	defer response.Body.Close()

	// Read response
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New(errors.ErrInternalServer, "failed to read response", err)
	}

	// Parse response
	var historicalData models.HistoricalDataResponse
	if err := json.Unmarshal(body, &historicalData); err != nil {
		return nil, errors.New(errors.ErrInternalServer, "failed to parse response", err)
	}

	// Cache the response
	c.addToCache(cacheKey, historicalData)

	return &historicalData, nil
}

func (c *Client) buildHistoricalDataURL(req models.HistoricalDataRequest) string {
	if req.IsIntraday {
		return fmt.Sprintf("%s/historical-candle/intraday/%s/%s",
			c.config.BaseURL,
			req.Symbol,
			req.Interval,
		)
	}

	return fmt.Sprintf("%s/historical-candle/%s/%s/%s/%s",
		c.config.BaseURL,
		req.Symbol,
		req.Interval,
		req.EndDate.Format("2006-01-02"),
		req.StartDate.Format("2006-01-02"),
	)
}

func (c *Client) getFromCache(key string) (models.HistoricalDataResponse, bool) {
	c.cache.mutex.RLock()
	defer c.cache.mutex.RUnlock()
	val, ok := c.cache.data[key]
	return val, ok
}

func (c *Client) addToCache(key string, data models.HistoricalDataResponse) {
	c.cache.mutex.Lock()
	defer c.cache.mutex.Unlock()
	c.cache.data[key] = data
}
