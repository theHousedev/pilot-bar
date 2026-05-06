package main

import (
	"fmt"
	"log/slog"
	"os/exec"
	"strings"

	"github.com/house-holder/pilot-bar/internal/cache"
	"github.com/house-holder/pilot-bar/pkg/types"
)

const waybarSignal = 8

func switchAirport(icao string, flags Flags) error {
	icao = strings.ToUpper(icao)
	if len(icao) != 4 {
		return fmt.Errorf("invalid ICAO identifier: %q (expected 4 characters)", icao)
	}

	slog.Info("Switching airport", "icao", icao)

	if err := cache.Write(types.Airport{
		ICAO: icao,
		METAR: types.METAR{
			Reported: types.Timestamp{Epoch: 0},
		},
	}); err != nil {
		return fmt.Errorf("cache reset failed: %w", err)
	}

	*flags.Airport = icao
	if err := Update(flags); err != nil {
		return err
	}

	wx, err := cache.Read()
	if err == nil {
		fmt.Println(wx.METAR.RawOb)
	}

	signalWaybar()
	return nil
}

func changeAirport(flags Flags) error {
	icao, err := promptUser()
	if err != nil {
		return err
	}
	return switchAirport(icao, flags)
}

func promptUser() (string, error) {
	cmd := exec.Command("wofi", "--dmenu", "--prompt", "Airport ICAO:")
	output, err := cmd.Output()
	return strings.TrimSpace(string(output)), err
}

func signalWaybar() {
	if err := exec.Command("pkill", fmt.Sprintf("-RTMIN+%d", waybarSignal), "waybar").Run(); err != nil {
		slog.Debug("waybar signal failed (waybar may not be running)", "error", err)
	}
}
