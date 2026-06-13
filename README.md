# pilot-bar

Aviation weather data in your Waybar. Fetches METAR, TAF, and Area Forecast Discussion (AFD) for a selected airport and displays a configurable summary.

Daemon updates weather on a schedule; the waybar companion binary reads the cached result and formats it for display.

I am a novice coder. So far it's working pretty well for me!

## Installation

### Prerequisites

- Waybar
- Go 1.25+
- [wofi](https://sr.ht/~scoopta/wofi/) (for interactive airport switching)
- A Nerd font to support the sky cover symbols

*(At time of writing, `wofi` is a working Arch package.)*

### Building from source

```sh
git clone https://github.com/house-holder/pilot-bar
cd pilot-bar
chmod +x install.sh
./install.sh
```

### Create a Waybar module

Add a custom module to your `.config/waybar/config.jsonc`:

```json
"custom/wx": {
    "exec": "pilot-bar-daemon && pilot-bar",
    "return-type": "json",
    "signal": 8,
    "tooltip": true,
    "on-click-right": "pilot-bar-daemon change",
    "exec-on-event": true
}
```

### Styling the output

These default colors work well for flight category at a glance, but it's obviously configurable in your `.config/waybar/style.css`:

```css
#custom-wx.vfr {
	color: #9ece6a;
}

#custom-wx.mvfr {
	color: #0080ff;
}

#custom-wx.ifr {
	color: #ff0000;
}

#custom-wx.lifr {
	color: #ff00ff;
}

```

### Switch airports

The Waybar module above switches on right click, or from your shell:

```sh
pilot-bar-daemon switch KSTL
```

## Configuration

Create your config file: `~/.config/pilot-bar/config.json` (or `$XDG_CONFIG_HOME/pilot-bar/config.json`)

### Example

```json
{
    "airport": "KLAX",
    "format": "{temps} {vis} {cloud-icon} {clouds} {wx}",
    "tempUnit": "C",
    "modules": {
        "metar": false,
        "taf": true,
        "discussion": true,
    }
}
```

### Config Options

| Field | Description | Default |
|---|---|---|
| `airport` | ICAO airport code | — |
| `format` | Display format string (see tokens below) | `{temps} {vis} {cloud-icon} {clouds} {wx}` |
| `tempUnit` | Temperature unit: `"C"` or `"F"` | `"C"` |
| `modules` | Which data types to fetch | `{"metar": true}` |

### Format tokens

| Token | Description | Example |
|---|---|---|---|
| `{temps}` | Ambient / Dewpoint (in configured unit) | `12.8/10.2` |
| `{temp}` | Ambient only (in configured unit) | `12.8` |
| `{dewpoint}` | Dewpoint only (in configured unit) | `10.2` |
| `{winds}` | Wind direction/speed/gusts | `120/10G15` |
| `{cloud-icon}` | Nerd Font icon for ceiling layer | 󰪣 |
| `{clouds}` | Ceiling altitude (hundreds of ft) | `030` |
| `{vis}` | Visibility (shown if < 6SM) | `3SM` |
| `{wx}` | Weather string | `-RA BR` |
| `{stationID}` | ICAO code | `KCGI` |
| `{age}` | Minutes since observation | `13` |
| `{fltcat}` | Flight category | `IFR` |
| `{altimeter}` | Altimeter in inHg | `29.92` |

### Cloud icons (Nerd Fonts)

| Coverage | Icon | Codepoint |
|---|---|---|
| FEW | 󰪟 | `U+F0A9F` |
| SCT | 󰪡 | `U+F0AA1` |
| BKN | 󰪣 | `U+F0AA3` |
| OVC | 󰪥 | `U+F0AA5` |


## Data sources

- **METAR / TAF**: [Aviation Weather Center API](https://aviationweather.gov/data/example)
- **AFD**: [Aviation Weather Center](https://aviationweather.gov/) via CWA lookup
- **CWA resolution**: [National Weather Service API](https://api.weather.gov/)

## Project structure

```
cmd/
  daemon/       # Background fetcher daemon
    main.go       Entry point, flag parsing, subcommand routing
    updater.go    Weather data fetch orchestration
    switch.go     Airport switching with wofi integration
    logger.go     Colored logging
    debug.go      Verbose METAR display
  waybar/
    main.go       Waybar JSON output formatter
internal/
  cache/          Disk cache (~/.cache/pilot-bar/currentWX.json)
  config/         Config file loading
  fetch/          HTTP clients for weather APIs
  parse/          METAR string parsing
pkg/
  types/          Shared data types (METAR, Airport, etc.)
cfg/
  config.example.json
```
