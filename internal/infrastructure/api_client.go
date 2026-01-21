package infrastructure

import (
	"context"
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
			Timeout: 10 * time.Second,
		},
	}
}

func (c *APIClient) FetchLogEntry(ctx context.Context) ([]byte, error) {
	var lastErr error
	maxRetries := 3

	for i := 0; i < maxRetries; i++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/logs", nil)
		if err != nil {
			return nil, err
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			time.Sleep(100 * time.Millisecond)
			continue
		}

		if resp.StatusCode >= 500 {
			lastErr = fmt.Errorf("unexpected status code: %d", resp.StatusCode)
			resp.Body.Close()
			time.Sleep(100 * time.Millisecond)
			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		return io.ReadAll(resp.Body)
	}

	return nil, fmt.Errorf("failed to fetch log entry after %d attempts: %w", maxRetries, lastErr)
}
