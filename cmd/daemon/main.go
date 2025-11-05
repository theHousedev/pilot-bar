package main

import (
	"log/slog"

	"github.com/spf13/pflag"
)

type Flags struct {
	Airport *string
	Update  *bool
	Debug   *bool
	Info    *bool
}

func setupFlags() Flags {
	info := pflag.BoolP("info", "i", false, "enable info logging")
	debug := pflag.BoolP("debug", "d", false, "enable debug logging")
	update := pflag.BoolP("update", "u", false, "update station data")

	defaultID, err := resolveAirport()
	if err != nil {
		slog.Error("failed to resolve default airport", "error", err)
		defaultID = "KCGI"
	}
	airport := pflag.StringP("airport", "a", defaultID, "target station ID")

	pflag.Parse()
	return Flags{
		Airport: airport,
		Update:  update,
		Debug:   debug,
		Info:    info,
	}
}

func main() {
	flags := setupFlags()
	InitLogger(flags)

	// TODO: when not testing, update should be automated
	if err := Update(flags); err != nil {
		slog.Error("Update", "error", err)
	}
}
