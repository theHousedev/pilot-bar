package fetch

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/theHousedev/pilot-bar/pkg/types"
)

const baseURL = "https://aviationweather.gov/api/data"

// FetchMETAR loads full report into a default-shaped struct
func FetchMETAR(icao string, maxAttempts int) (types.METARresponse, error) {
	if maxAttempts < 1 {
		maxAttempts = 1
	}

	metarURL := fmt.Sprintf("%s/metar?ids=%s&format=json", baseURL, icao)
	client := &http.Client{Timeout: 10 * time.Second}
	startTime := time.Now()

	var payload []types.METARresponse
	err := doWithRetry(maxAttempts, func(attempt int) (bool, error) {
		if attempt > 1 {
			slog.Info(fmt.Sprintf("Fetch METAR retry (%d of %d)", attempt, maxAttempts))
		} else {
			slog.Info("Fetching METAR")
		}

		resp, err := client.Get(metarURL)
		if err != nil {
			var netErr net.Error
			if errors.As(err, &netErr) && netErr.Timeout() {
				slog.Warn("Fetch timeout", "attempt", attempt, "max", maxAttempts)
				return true, err
			}
			return false, fmt.Errorf("HTTP request failed: %w", err)
		}
		defer resp.Body.Close()

		if statusRetryOK(resp.StatusCode) {
			slog.Warn("OK to retry", "status", resp.Status, "attempt", attempt)
			return true, fmt.Errorf("status %d: %s", resp.StatusCode, resp.Status)
		}

		if resp.StatusCode != http.StatusOK {
			return false, fmt.Errorf("status: %s", resp.Status)
		}

		var decoded []types.METARresponse
		if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
			return false, fmt.Errorf("decode failed: %w", err)
		}
		if len(decoded) == 0 {
			return false, fmt.Errorf("no METAR data for %s", icao)
		}

		payload = decoded
		return false, nil
	})

	if err != nil {
		return types.METARresponse{}, err
	}

	fetchDuration := time.Since(startTime).Seconds()
	slog.Info("Fetch OK", "took", fmt.Sprintf("%.3fs", fetchDuration))
	return payload[0], nil
}

func doWithRetry(maxAttempts int, op func(attempt int) (bool, error)) error {
	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		retry, err := op(attempt)
		if err == nil {
			return nil
		}

		lastErr = err
		if !retry || attempt == maxAttempts {
			break
		}
		time.Sleep(2 * time.Second) // backoff delay
	}

	return lastErr
}

func statusRetryOK(code int) bool {
	switch code {
	case http.StatusRequestTimeout, http.StatusTooManyRequests,
		http.StatusInternalServerError, http.StatusBadGateway,
		http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}
