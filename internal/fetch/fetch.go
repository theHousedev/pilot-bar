package fetch

import (
    "fmt"
    "net/http"
	"encoding/json"
    "time"

    "github.com/theHousedev/pilot-bar/pkg/types"
)

const baseURL = "https://aviationweather.gov/api/data"

// FetchMETAR retrieves METAR json
func FetchMETAR(icao string) (*types.METARResponse, error) {
    metarURL := fmt.Sprintf("%s/metar?ids=%s&format=json", baseURL, icao)
    
    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Get(metarURL)
    if err != nil {
        return nil, fmt.Errorf("HTTP request failed: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
    }

    var results []types.METARResponse
    if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
        return nil, fmt.Errorf("failed to decode JSON: %w", err)
    }

    if len(results) == 0 {
        return nil, fmt.Errorf("no METAR data for %s", icao)
    }

    return &results[0], nil
}
