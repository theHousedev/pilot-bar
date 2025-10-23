package types

// Alias typing pretending to be enums
type Sky string

const (
	SkyCLR Sky = "clear"
	SkyFEW Sky = "few"
	SkySCT Sky = "scattered"
	SkyBKN Sky = "broken"
	SkyOVC Sky = "overcast"
)

// Subcomponent types =========================================================
type Timestamp struct {
	Day  int `json:"day"`
	Time int `json:"time"`
}

type SiteWX struct {
	Station string    `json:"station"`
	Issued  Timestamp `json:"issued"`
}

type SurfaceWind struct {
	From  int `json:"from"`
	Speed int `json:"speed"`
}

type CloudLayer struct {
	Coverage  string `json:"coverage"`
	HeightAGL int    `json:"height"`
}

type TempData struct {
	Ambient       int     `json:"ambient"`
	Dewpoint      int     `json:"dewpoint"`
	ExactAmbient  float64 `json:"exact_ambient"`
	ExactDewpoint float64 `json:"exact_dewpoint"`
}

// Component types =========================================================
type METAR struct {
	RawData string       `json:"raw"`
	Wind    SurfaceWind  `json:"winds"`
	Clouds  []CloudLayer `json:"clouds"`
	Vis     int          `json:"vis"`
	Temps   TempData     `json:"temps"`
}

type TAF struct {
	Raw string
}
