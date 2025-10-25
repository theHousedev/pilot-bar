package main

import (
    "flag"
    "fmt"
    "os"
    "log/slog"

    "github.com/theHousedev/pilot-bar/internal/fetch"
)

func main() {
    airport := flag.String("airport", "KCGI", "the airport targeted by fetch")
    debug := flag.Bool("debug", false, "enable debug logging")
    flag.Parse()

    if *debug {
        slog.Info("Begin fetch.", "Airport ID", *airport)
    }
    metar, err := fetch.FetchMETAR(*airport)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error fetching METAR: %v\n", err)
        os.Exit(1)
    }
	if *debug {
		displayMETAR(metar)
	}
}
