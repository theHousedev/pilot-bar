package main

import (
	"log/slog"

	"github.com/theHousedev/pilot-bar/internal/fetch"
)

func UpdateWX(flags Flags) error {
	slog.Info("Update triggered", "airport", *flags.Airport)

	maxAttempts := 5

	metar, err := fetch.FetchMETAR(*flags.Airport, maxAttempts)
	if err != nil {
		return err
	}

	if *flags.Debug {
		displayMETAR(metar)
	} else {
		slog.Info("", "METAR", metar.RawOb)
	}

	return nil
}
