// an attempt to make aliases that assist in the code documenting itself
package types

type (
	LocalTime uint16 // 0000-2359
	ZuluTime  uint16 // 0000-2359
	DegMag    uint16 // 1-360, degrees magnetic
	Knots     int
	Feet      int
	Mi        int
	InHg      float64
)

type Timestamp struct {
	Day   uint8     `json:"day"` // validate: 1-31
	Zulu  ZuluTime  `json:"zulu"`
	Local LocalTime `json:"local"`
	Epoch int64     `json:"epoch"`
}
