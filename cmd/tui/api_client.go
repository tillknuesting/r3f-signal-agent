package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type APIClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type TrendDTO struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Summary     string    `json:"summary"`
	Score       int       `json:"score"`
	Author      string    `json:"author"`
	Source      string    `json:"source"`
	SourceID    string    `json:"source_id"`
	Timestamp   time.Time `json:"timestamp"`
	CollectedAt time.Time `json:"collected_at"`
	Starred     bool      `json:"starred"`
}

type SourceDTO struct {
	ID          string        `yaml:"id" json:"id"`
	Name        string        `yaml:"name" json:"name"`
	Description string        `yaml:"description" json:"description"`
	Type        string        `yaml:"type" json:"type"`
	Enabled     bool          `yaml:"enabled" json:"enabled"`
	Display     DisplayConfig `yaml:"display" json:"display"`
}

type DisplayConfig struct {
	Icon  string `yaml:"icon" json:"icon"`
	Color string `yaml:"color" json:"color"`
}

type TrendsResponse struct {
	Trends []TrendDTO `json:"trends"`
	Total  int        `json:"total"`
}

type SourcesResponse struct {
	Sources []SourceDTO `json:"sources"`
}

type CollectResponse struct {
	Status     string   `json:"status"`
	ItemsCount int      `json:"items_count"`
	Errors     []string `json:"errors"`
}

func (c *APIClient) GetTrends() (*TrendsResponse, error) {
	resp, err := c.httpClient.Get(c.baseURL + "/api/v1/trends")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var trends TrendsResponse
	if err := json.Unmarshal(body, &trends); err != nil {
		return nil, err
	}

	return &trends, nil
}

func (c *APIClient) GetSources() (*SourcesResponse, error) {
	resp, err := c.httpClient.Get(c.baseURL + "/api/v1/sources")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var sources SourcesResponse
	if err := json.Unmarshal(body, &sources); err != nil {
		return nil, err
	}

	return &sources, nil
}

func (c *APIClient) Collect() (*CollectResponse, error) {
	resp, err := c.httpClient.Post(c.baseURL+"/api/v1/collect", "application/json", bytes.NewReader([]byte{}))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result CollectResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *APIClient) StarTrend(id string) error {
	url := fmt.Sprintf("%s/api/v1/trends/%s/star", c.baseURL, id)
	resp, err := c.httpClient.Post(url, "application/json", bytes.NewReader([]byte{}))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("failed to star trend: HTTP %d", resp.StatusCode)
	}

	return nil
}

func (c *APIClient) HealthCheck() error {
	resp, err := c.httpClient.Get(c.baseURL + "/api/v1/health")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("health check failed: HTTP %d", resp.StatusCode)
	}

	return nil
}
