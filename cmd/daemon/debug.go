package main

import (
	"fmt"

	"github.com/house-holder/pilot-bar/pkg/types"
)

func displayMETAR(data types.METARresponse) error {
	fmt.Println("  IcaoID...........", data.IcaoID)
	fmt.Println("  ReceiptTime......", data.ReceiptTime)
	fmt.Println("  ObsTime..........", data.ObsTime)
	fmt.Println("  ReportTime.......", data.ReportTime)
	fmt.Println("  Metar Type.......", data.MetarType)
	fmt.Println("  Temp.............", data.Temp)
	fmt.Println("  Dewp.............", data.Dewp)
	fmt.Println("  Wind dir.........", data.Wdir)
	fmt.Println("  Wind speed.......", data.Wspd)
	fmt.Println("  Visib............", data.Visib)
	fmt.Println("  Altimeter........", fmt.Sprintf("%.1f", data.Altim))
	if data.WxString != "" {
		fmt.Println("  + Wx string......", data.WxString)
	}
	if data.Slp != 0 {
		fmt.Println("  + SLP............", data.Slp)
	}
	if data.PresTend != nil {
		fmt.Println("  + Pres Tend......", *data.PresTend)
	}
	if data.MaxT != nil {
		fmt.Println("  + Max Temp.......", *data.MaxT)
	}
	if data.MinT != nil {
		fmt.Println("  + Min Temp.......", *data.MinT)
	}
	if data.MaxT24 != nil {
		fmt.Println("  + 24h Max Temp...", *data.MaxT24)
	}
	if data.MinT24 != nil {
		fmt.Println("  + 24h Min Temp...", *data.MinT24)
	}
	if data.Precip != nil {
		fmt.Println("  + Precip.........", *data.Precip)
	}
	if data.Pcp3hr != nil {
		fmt.Println("  + Precip 3hr.....", *data.Pcp3hr)
	}
	if data.Pcp6hr != nil {
		fmt.Println("  + Precip 6hr.....", *data.Pcp6hr)
	}
	if data.Pcp24hr != nil {
		fmt.Println("  + Precip 24hr....", *data.Pcp24hr)
	}
	if data.Snow != nil {
		fmt.Println("  + Snow...........", *data.Snow)
	}
	if data.VertVis != nil {
		fmt.Println("  + Vertical vis...", *data.VertVis)
	}
	fmt.Println("  QC Field.........", data.QcField)
	fmt.Println("  Raw METAR........", data.RawOb)
	fmt.Println("  Latitude.........", data.Lat)
	fmt.Println("  Longitude........", data.Long)
	fmt.Println("  Elevation........", data.Elev)
	fmt.Println("  Airport name.....", data.Name)
	fmt.Println("  Cover............", data.Cover)
	if len(data.Clouds) > 0 {
		for i, cloud := range data.Clouds {
			fmt.Printf("    + Layer %d: %s %5d\n", i+1, cloud.Cover, cloud.Base)
		}
	}
	fmt.Println("  FltCat...........", data.FltCat)
	return nil
}
