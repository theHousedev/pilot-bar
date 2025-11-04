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
	airport := pflag.StringP("airport", "a", "KCGI", "target station ID")

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

	// TODO: wrap in update controller when not testing
	if err := UpdateWX(flags); err != nil {
		slog.Error("UpdateWX", "error", err)
	}
}
