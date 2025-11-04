package parse

import (
	"log/slog"
	"slices"
	"strconv"
	"strings"

	"github.com/theHousedev/pilot-bar/pkg/types"
)

type ParseContext struct {
	tokens []string
	index  int
	output *types.METAR
}

// --------------------------------------------------------------------------------------
type parseFunc func(c *ParseContext) (int, error) // int: tokens used

func (c *ParseContext) advance(n int)   { c.index += n }
func (c *ParseContext) current() string { return c.tokens[c.index] }
func (c *ParseContext) next() string    { c.index++; return c.current() }
func (c *ParseContext) inRange() bool   { return c.index < len(c.tokens) }

func (c *ParseContext) getToken(i int) string {
	target := c.tokens[i]
	c.tokens = slices.Delete(c.tokens, i, i+1)
	return target
}

func HandleMETAR(data *types.METARresponse) error {
	parsers := []parseFunc{
		getAltimeter, getGusts, getWXString, getRemarks,
	}
	c := &ParseContext{
		tokens: strings.Split(data.RawOb, " "),
		index:  0,
		output: &types.METAR{},
	}
	slog.Debug("handleMETAR origin",
		"tokens", c.tokens,
		"index", c.index,
	)

	for _, parser := range parsers {
		used, err := parser(c)
		if err != nil {
			return err
		}
		c.advance(used)
	}
	return nil
}

func getAltimeter(ctx *ParseContext) (int, error) {
	for i, token := range ctx.tokens {
		if len(token) == 5 && strings.HasPrefix(token, "A") {
			altVal, err := strconv.ParseFloat(ctx.getToken(i)[1:], 64)
			if err != nil {
				return 0, err
			}
			ctx.output.Altimeter = types.InHg(altVal / 100.0)
			break
		}
	}
	return 1, nil
}

func getGusts(ctx *ParseContext) (int, error) {
	return 0, nil
}

func getWXString(ctx *ParseContext) (int, error) {
	return 0, nil
}

func getRemarks(ctx *ParseContext) (int, error) {
	return 0, nil
}
