package main

import (
	"encoding/json"
	"log/slog"
	"os"
	"time"

	"github.com/theHousedev/pilot-bar/internal/fetch"
	"github.com/theHousedev/pilot-bar/internal/parse"
	"github.com/theHousedev/pilot-bar/pkg/types"
)

const (
	MaxTries  = 5
	CachePath = "./testdata/currentWX.json"
	// TODO: seconds vals below, 1min for testing - must edit for prod
	UpdateInterval = 60
	IntervalMETAR  = 60
	IntervalTAF    = 60
	IntervalAFD    = 60
)

type UpdateData struct {
	cached        types.Airport
	requested     string
	now           int64
	intervalMETAR int64
	intervalTAF   int64
	intervalAFD   int64
}

func (d *UpdateData) AirportChanged() bool {
	return d.cached.ICAO != d.requested
}

func (d *UpdateData) METAREmpty() bool {
	return d.cached.METAR.Reported.Epoch == 0
}

func (d *UpdateData) TimeExpired() bool {
	return d.now-d.cached.LastUpdateEpoch > UpdateInterval
}

func (d *UpdateData) NeedsAnyUpdate(force bool) bool {
	if force || d.AirportChanged() || d.METAREmpty() || d.TimeExpired() {
		slog.Debug("Proceeding with update", "list", map[string]any{
			"New ID":       d.AirportChanged(),
			"Timeout":      d.TimeExpired(),
			"Force update": force,
		})
		return true
	}
	slog.Debug("No update needed")
	return false
}

func Update(flags Flags) error {
	cachedWX, err := readCachedWX(CachePath, flags)
	if err != nil {
		return err
	}

	d := &UpdateData{
		cached:        cachedWX,
		requested:     *flags.Airport,
		now:           time.Now().Unix(),
		intervalMETAR: IntervalMETAR,
		intervalTAF:   IntervalTAF,
		intervalAFD:   IntervalAFD,
	}

	if !d.NeedsAnyUpdate(*flags.Update) {
		return nil
	}

	APImetar, err := fetch.GetMETAR(*flags.Airport, MaxTries)
	if err != nil {
		return err
	}

	if *flags.Verbose {
		displayMETAR(APImetar)
	} else {
		slog.Debug("", "metar", APImetar.RawOb)
	}

	if err := parse.BuildInternalMETAR(&APImetar, &cachedWX.METAR); err != nil {
		return err
	}

	if d.AirportChanged() {
		cachedWX.ICAO = *flags.Airport
		cachedWX.Elevation = types.Feet(float64(APImetar.Elev) * 3.28084)
	}

	cachedWX.LastUpdateEpoch = time.Now().Unix()
	cachedWX.METAR.Reported.Epoch = int64(APImetar.ObsTime)

	return writeCachedWX(CachePath, cachedWX)
}

func ensureCacheExists(jsonPath string, flags Flags) error {
	if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
		return writeCachedWX(jsonPath, types.Airport{
			ICAO:            *flags.Airport,
			LastUpdateEpoch: time.Now().Unix(),
			METAR: types.METAR{
				Reported: types.Timestamp{
					Epoch: 0, // empty initial
				},
			},
		})
	}
	return nil
}

func readCachedWX(jsonPath string, flags Flags) (types.Airport, error) {
	err := ensureCacheExists(jsonPath, flags)
	if err != nil {
		return types.Airport{}, err
	}
	jsonData, err := os.ReadFile(jsonPath)
	if err != nil {
		return types.Airport{}, err
	}
	var cachedWX types.Airport
	err = json.Unmarshal(jsonData, &cachedWX)
	return cachedWX, nil
}

func writeCachedWX(jsonPath string, cachedWX types.Airport) error {
	jsonData, err := json.MarshalIndent(cachedWX, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(jsonPath, jsonData, 0644)
}

func getCachedICAO(cachePath string) (string, error) {
	jsonData, err := os.ReadFile(cachePath)
	if err != nil {
		return "", err
	}
	var cached types.Airport
	if err := json.Unmarshal(jsonData, &cached); err != nil {
		return "", err
	}
	return cached.ICAO, nil
}

func resolveAirport() (string, error) {
	icao, err := getCachedICAO(CachePath)
	if err == nil && icao != "" {
		return icao, nil
	}
	slog.Warn("no cached or provided airport. using default: KCGI")
	return "KCGI", nil
}
