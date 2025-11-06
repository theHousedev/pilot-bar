package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

const (
	colorReset  = "\x1b[0m"
	colorDim    = "\x1b[2m"
	colorCyan   = "\x1b[36m"
	colorGreen  = "\x1b[32m"
	colorYellow = "\x1b[33m"
	colorRed    = "\x1b[31m"
)

type colorHandler struct {
	level  slog.Level
	writer io.Writer
}

func newColorHandler(w io.Writer, level slog.Level) slog.Handler {
	return &colorHandler{level: level, writer: w}
}

func (h *colorHandler) Enabled(_ context.Context, lvl slog.Level) bool {
	return lvl >= h.level
}

func (h *colorHandler) Handle(ctx context.Context, r slog.Record) error {
	if !h.Enabled(ctx, r.Level) {
		return nil
	}

	var b strings.Builder
	fmt.Fprintf(&b, "[%s%s%s%s] %s %s",
		colorDim, colorCyan, r.Time.Format("01/02 1504:05.000"), colorReset,
		levelLabel(r.Level),
		r.Message,
	)

	r.Attrs(func(a slog.Attr) bool {
		key := strings.ToLower(a.Key)
		switch key {
		case "metar":
			fmt.Fprintf(&b, "\n%s", a.Value.String())
			return false
		case "airport":
			fmt.Fprintf(&b, " [%s]", a.Value.String())
			return true
		case "took":
			fmt.Fprintf(&b, " (%s)", a.Value.String())
			return true
		case "error":
			if r.Level >= slog.LevelError {
				fmt.Fprintf(&b, " [%s%s%s]", colorRed, a.Value.String(), colorReset)
			} else {
				fmt.Fprintf(&b, "%s", a.Value.String())
			}
			return false
		case "list":
			m, ok := a.Value.Any().(map[string]any)
			if !ok {
				return true
			}
			if len(m) == 0 {
				return true
			}

			maxLen := 0
			for key := range m {
				if len(key) > maxLen {
					maxLen = len(key)
				}
			}

			for key, value := range m {
				var valStr string
				switch value.(type) {
				case string:
					valStr = fmt.Sprintf("%q", value)
				case bool:
					if value == true {
						valStr = fmt.Sprintf("%s%t%s", colorGreen, value, colorReset)
					} else {
						valStr = fmt.Sprintf("%s%t%s", colorRed, value, colorReset)
					}
				case int, int8, int16, int32, int64, uint, uint8, uint16,
					uint32, uint64, float32, float64:
					valStr = fmt.Sprintf("%s%v%s", colorYellow, value, colorReset)
				}

				if maxLen > 0 {
					fmt.Fprintf(
						&b, "\n%s%s    %-*s: %s%s%s%s",
						colorGreen, colorDim, maxLen, key, colorReset,
						colorGreen, valStr, colorReset)
				} else {
					fmt.Fprintf(&b, "\n%s    %s: %s%s%s%s%s",
						colorGreen, colorDim, key, colorReset,
						colorGreen, valStr, colorReset)
				}
			}
			fmt.Fprintf(&b, colorReset)
			return true
		case "tokens":
			fmt.Fprintf(&b, "\n%s", a.Value.String())
			return true
		default:
			fmt.Fprintf(&b, " %s=%s", a.Key, a.Value.String())
			return true
		}
	})

	b.WriteByte('\n')
	_, err := io.WriteString(h.writer, b.String())
	return err
}

func (h *colorHandler) WithAttrs(_ []slog.Attr) slog.Handler { return h }
func (h *colorHandler) WithGroup(_ string) slog.Handler      { return h }

func levelLabel(level slog.Level) string {
	switch {
	case level >= slog.LevelError:
		return fmt.Sprintf("%sERROR%s", colorRed, colorReset)
	case level >= slog.LevelWarn:
		return fmt.Sprintf("%sWARN:%s", colorYellow, colorReset)
	case level >= slog.LevelInfo:
		return fmt.Sprintf("%s%sINFO:%s", colorDim, colorGreen, colorReset)
	default:
		return fmt.Sprintf("%s%sDEBUG%s", colorDim, colorYellow, colorReset)
	}
}

func InitLogger(flags Flags) {
	var level slog.Level
	switch {
	case *flags.Debug:
		level = slog.LevelDebug
	case *flags.Info:
		level = slog.LevelInfo
	default:
		level = slog.LevelWarn
	}

	handler := newColorHandler(os.Stderr, level)
	slog.SetDefault(slog.New(handler))
}
