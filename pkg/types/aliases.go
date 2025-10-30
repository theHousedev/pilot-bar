// an attempt to make aliases that assist in the code documenting itself
package types

type (
	time      uint16 // 0000-2359
	zTime     uint16 // 0000-2359
	degM      uint16 // 1-360, degrees magnetic
	knots     int
	feet      int
	statuteMi int
	inHg      float64
)

type Timestamp struct {
	Day   uint8 `json:"day"` // validate: 1-31
	Zulu  zTime `json:"zulu"`
	Local time  `json:"local"`
	Epoch int64 `json:"epoch"`
}
