package parse

import (
	"log/slog"
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

func HandleMETAR(data *types.METARresponse) (int, error) {
	parsers := []parseFunc{
		getAltimeter, getGusts, getWXString, getRemarks,
	}
	ctx := &ParseContext{
		tokens: strings.Split(data.RawOb, " "),
		index:  0,
		output: &types.METAR{},
	}
	slog.Debug("handleMETAR origin",
		"tokens", ctx.tokens,
		"index", ctx.index,
	)

	for _, parser := range parsers {
		used, err := parser(ctx)
		if err != nil {
			return 0, err
		}
		ctx.advance(used)
	}

	return 0, nil
}

func getAltimeter(ctx *ParseContext) (int, error) {
	return 0, nil
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
