# devnotes
## status
**working:**
- fetch w/timeout
- JSON unmarshaling
- flags:
    - `--airport` to test varied data
    - `--debug` essentially -V 

**wip:**
- type defs
    - fix alignment with SLP
    - make stand-in conversion for alt
- debug display (HIVAL)
    - basic data + formatted printing

**roadmap:**
- config loading
- caching
    - per-airport dirs
- parsing
    - altimeter
    - remarks
- testing
    - only able to test parsing? unsure what else
- Waybar integrations
    - performance oriented

## architecture
### modules
```
cmd/
  daemon/     - main entry point, fetch w/debug display
  waybar/     - integration
internal/
  fetch/      - HTTP + JSON
  parse/      - remarks, altimeter
  config/     - XDG-compliant uconf
  cache/      - local caching
pkg/
  types/      - shared typedefs
```

### key points
- daemon is long-running, waybar process only called when change detected (perf+)
- altimeter can fallback to simple conversion if needed
- parsing is mostly nonessential - nice-to-haves (tooltip?):
    - weather event timing
    - pressure trends
    - windshear? other low probability phenomena

## add
- save complex/long example METARs/TAFs (json & raw) for eventual full-scale test
- lat/long available: may be able to define geobounds w/ radius, subregion, etc

