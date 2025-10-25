package main

import (
	"fmt"
	"log/slog"

	"github.com/theHousedev/pilot-bar/pkg/types"
)

func displayMETAR(data *types.METARResponse) {
	slog.Info("METAR", "Airport ID", data.IcaoID)
	fmt.Println("Raw...........", data.RawMETAR)
	fmt.Println("ReportTime....", data.ReportTime)
	fmt.Println("Temp..........", data.Temp, "C")
	fmt.Println("Dewpoint......", data.Dewpoint, "C"	)
	fmt.Println("WindDir.......", data.WindDir)
	fmt.Println("WindSpeed.....", data.WindSpeed, "Kt")
	fmt.Println("Visibility....", data.Visibility, "SM")
	fmt.Println("Altimeter.....", data.Altimeter, "hPa")
	fmt.Println("FlightCat.....", data.FlightCat)
	fmt.Println("Clouds........", data.Clouds)
	fmt.Println("Name..........", data.Name)
}