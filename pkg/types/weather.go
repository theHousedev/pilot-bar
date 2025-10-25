package types

type Sky string // Alias typing (FIAE)

const (
	SkyClear     Sky = "clear"
	SkyFew       Sky = "few"
	SkyScattered Sky = "scattered"
	SkyBroken    Sky = "broken"
	SkyOvercast  Sky = "overcast"
)

type Timestamp struct {
	Day  int `json:"day"`
	Time int `json:"time"`
}

type WeatherData struct {
	Station string         `json:"station"`
	METAR   ComponentMETAR `json:"metar"`
	TAF     ComponentTAF   `json:"taf"`
}

type SurfaceWind struct {
	From  int `json:"from"`
	Speed int `json:"speed"`
	Gusts int `json:"gusts"`
}

type CloudLayer struct {
    Coverage string `json:"cover"`
    Base     int    `json:"base"`
}

type TempData struct {
	Ambient       int     `json:"ambient"`
	Dewpoint      int     `json:"dewpoint"`
	ExactAmbient  float64 `json:"exact_ambient"`
	ExactDewpoint float64 `json:"exact_dewpoint"`
}

type Remarks struct {
	// TODO: diverse possible contents
}

type ComponentMETAR struct {
	RawData string       `json:"raw_data"`
	Special bool         `json:"special"`
	Issued  Timestamp    `json:"issued"`
	Wind    SurfaceWind  `json:"wind"`
	Clouds  []CloudLayer `json:"clouds"`
	Vis     int          `json:"vis"`
	Temps   TempData     `json:"temps"`
	Baro    float64      `json:"baro"`
}

type TAFHeader struct {
	Issued Timestamp
}

type METARResponse struct {
    IcaoID      string        `json:"icaoId"`
    ReportTime  string        `json:"reportTime"`
    Temp        float64       `json:"temp"`
    Dewpoint    float64       `json:"dewp"`
    WindDir     int           `json:"wdir"`
    WindSpeed   int           `json:"wspd"`
    Visibility  string        `json:"visib"`
    Altimeter   float64       `json:"altim"`
    FlightCat   string        `json:"fltCat"`
    RawMETAR    string        `json:"rawOb"`
    Clouds      []CloudLayer  `json:"clouds"`
    Name        string        `json:"name"`
}

type TAFLine struct {
	From Timestamp   `json:"from"`
	Wind SurfaceWind `json:"wind"`
}

type ComponentTAF struct {
	RawData string    `json:"raw_data"`
	Issued  Timestamp `json:"issued"`
	// header line:
	//     TAF KCGI 241120Z 2412/2512 00000KT P6SM SCT250

	// part lines
	//     FM241400 10005KT P6SM BKN250
	//     FM251500 22012KT P6SM VCSH SCT020 OVC030
	//     FM241400 17010G20KT P6SM -RA OVC035 WS020/19040KT
}
