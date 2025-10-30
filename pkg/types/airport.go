package types

type Airport struct {
	ICAO            string `json:"icao"`
	LastUpdateEpoch int64  `json:"last_update"`
	Elevation       feet   `json:"elevation"`
	METAR           METAR  `json:"metar"`
}
