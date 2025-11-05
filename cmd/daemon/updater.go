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
	MaxTries   = 5
	CachePath  = "./testdata/currentWX.json"
	UpdateFreq = 60 * 1 // TODO: adjust for prod
)

func Update(flags Flags) error {
	slog.Info("Update triggered", "airport", *flags.Airport)

	if _, err := os.Stat(CachePath); os.IsNotExist(err) {
		// create initial cache file with basic info if it isn't found
		err := writeCachedWX(CachePath, types.Airport{
			ICAO:            *flags.Airport,
			LastUpdateEpoch: time.Now().Unix(),
			METAR: types.METAR{
				Reported: types.Timestamp{
					Epoch: 0, // empty entry to satisfy conditional below
				},
			},
		})
		if err != nil {
			return err
		}
	}
	cachedWX, err := readCachedWX(CachePath)
	if err != nil {
		return err
	}

	if cachedWX.LastUpdateEpoch < time.Now().Unix()-UpdateFreq {
		if cachedWX.ICAO == *flags.Airport {
			metarAPIResp, err := fetch.GetMETAR(*flags.Airport, MaxTries)
			if err != nil {
				return err
			}

			if cachedWX.METAR.Reported.Epoch == metarAPIResp.ObsTime {
				slog.Debug("fetched ObsTime matches Reported.Epoch")
			} else {
				slog.Debug("fetched ObsTime does not match Reported.Epoch")
				err = parse.BuildInternalMETAR(&metarAPIResp, &cachedWX.METAR)
				if err != nil {
					return err
				}

				cachedWX.LastUpdateEpoch = time.Now().Unix()
			}
			err = writeCachedWX(CachePath, cachedWX)
			if err != nil {
				return err
			}
		} else { // airport has changed
			slog.Info("new airport detected,", "icao", *flags.Airport)
			// TODO: full package update
			// 			- ICAO
			// 			- LastUpdateEpoch
			// 			- Elevation
			// 			- METAR
		}
	}
	return nil
}

func readCachedWX(jsonPath string) (types.Airport, error) {
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
