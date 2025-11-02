// NOTE: This package may only be needed for remarks section

package parse

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/theHousedev/pilot-bar/pkg/types"
)

type parseHelper func(token string, parsed *types.METAR) (bool, error)

var parsed = types.METAR{}

func parseAlt(token string, parsed *types.METAR) (bool, error) {
	fmt.Println(token)
	alt, err := strconv.Atoi(token[1:])
	if err != nil {
		return false, err
	}
	parsed.Altimeter = types.InHg(float64(alt) / 100.0)
	return true, nil
}

// parse all remarks
func BuildMETAR(data *types.METARresponse) (types.METAR, error) {
	raw := data.RawOb
	parserFunctions := []parseHelper{
		parseAlt,
	}

	// send individual tokens to
	for token := range strings.SplitSeq(raw, " ") {
		if token == "" {
			continue
		}
		handled := false
		for _, parser := range parserFunctions {
			ok, err := parser(token, &parsed)
			if err != nil {
				return types.METAR{}, err
			}
			if ok {
				handled = true
				break
			}
		}
		if !handled {
			return types.METAR{}, fmt.Errorf("unhandled token: %s", token)
		}
	}
	return parsed, nil
}
