package parse

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/theHousedev/pilot-bar/pkg/types"
)

type ParseContext struct {
	tokens []string
	index  int
	input  *types.METARresponse
	output *types.METAR
}

// --------------------------------------------------------------------------------------
type parseFunc func(c *ParseContext) error

func BuildInternalMETAR(data *types.METARresponse, output *types.METAR) error {
	output.Temp.AmbientExact = float64(data.Temp)
	output.Temp.DewpointExact = float64(data.Dewp)
	output.Temp.Ambient = int(data.Temp)
	output.Temp.Dewpoint = int(data.Dewp)
	output.Reported.Age = int(time.Since(time.Unix(data.ObsTime, 0)).Minutes())

	output.Clouds = make([]types.CloudData, 0)
	for _, layer := range data.Clouds {
		output.Clouds = append(output.Clouds, types.CloudData{
			Base:     types.Feet(layer.Base),
			Coverage: provideCloudCover(layer.Cover),
		})
	}

	// "receiptTime": "2025-10-28T14:56:19.211Z",
	output.Reported.Zulu = provideTimeData(data.ReportTime, "zulu")
	output.Reported.Local = provideTimeData(data.ReportTime, "local")

	parsers := []parseFunc{
		getAltimeter, getWind, getWXString, getRemarks,
	}
	c := &ParseContext{
		tokens: strings.Split(data.RawOb, " "),
		index:  0,
		output: output,
	}
	slog.Debug("handleMETAR origin",
		"tokens", c.tokens,
		"index", c.index,
	)

	for _, parser := range parsers {
		err := parser(c)
		if err != nil {
			return err
		}
	}
	return nil
}

func provideTimeData(timeString string, request string) types.Time {
	timeValue, err := time.Parse("2006-01-02T15:04:05Z", timeString)
	if err != nil {
		slog.Error("provideTimeData failed", "error", err)
		return types.Time{}
	}
	if request == "zulu" {
		return types.Time{
			Day:  uint8(timeValue.In(time.UTC).Day()),
			Hour: uint8(timeValue.In(time.UTC).Hour()),
		}
	} else {
		return types.Time{
			Day:  uint8(timeValue.In(time.Local).Day()),
			Hour: uint8(timeValue.In(time.Local).Hour()),
		}
	}
}

func provideCloudCover(coverage string) string {
	switch coverage {
	case "FEW":
		return "few"
	case "SCT":
		return "scattered"
	case "BKN":
		return "broken"
	case "OVC":
		return "overcast"
	default:
		return coverage
	}
}

func getAltimeter(ctx *ParseContext) error {
	for _, token := range ctx.tokens {
		if len(token) == 5 && strings.HasPrefix(token, "A") {
			altVal, err := strconv.ParseFloat(token[1:], 64)
			if err != nil {
				return fmt.Errorf("getAltimeter failed: %w", err)
			}
			ctx.output.Altimeter = types.InHg(altVal / 100.0)
			break
		}
	}
	return nil
}

func getWind(ctx *ParseContext) error {
	for _, token := range ctx.tokens {
		if strings.HasSuffix(token, "KT") {
			if token[:3] == "VRB" {
				ctx.output.Wind.Direction = types.DegMag(0)
				ctx.output.Wind.Variable = true
				speed, err := strconv.Atoi(token[3:5])
				if err != nil {
					return fmt.Errorf("getWind(0) failed: %w", err)
				}
				ctx.output.Wind.Speed = types.Knots(speed)
			} else {
				ctx.output.Wind.Variable = false
				direction, err := strconv.Atoi(token[0:3])
				if err != nil {
					return fmt.Errorf("getWind(1) failed: %w", err)
				}
				ctx.output.Wind.Direction = types.DegMag(direction)
				speed, err := strconv.Atoi(token[3:5])
				if err != nil {
					return fmt.Errorf("getWind(2) failed: %w", err)
				}
				ctx.output.Wind.Speed = types.Knots(speed)
			}

			if strings.Contains(token, "G") {
				gustsIdx := strings.Index(token, "G")
				gusts, err := strconv.Atoi(token[gustsIdx+1 : gustsIdx+3])
				if err != nil {
					return fmt.Errorf("getWind(3) failed: %w", err)
				}
				gustsValue := types.Knots(gusts)
				ctx.output.Wind.Gusts = &gustsValue
			} else {
				ctx.output.Wind.Gusts = nil
			}

			if ctx.output.Wind.Speed == 0 && ctx.output.Wind.Direction == 0 {
				ctx.output.Wind.Calm = true
			} else {
				ctx.output.Wind.Calm = false
			}
		}
	}
	return nil
}

func getWXString(ctx *ParseContext) error {
	return nil
}

func getRemarks(ctx *ParseContext) error {
	return nil
}
