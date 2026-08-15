// Package riskmodel is the domain core: it turns daily weather into per-year
// growth-stage risk features and ranks past years by climate similarity. It is
// pure — no GCP, network, or filesystem access — so it can be unit-tested in
// isolation.
package riskmodel

import "time"

// DailyWeather is one day of observations from a single JMA station.
type DailyWeather struct {
	Date     time.Time
	TMax     float64 // daily maximum temperature, °C
	TMin     float64 // daily minimum temperature, °C
	TMean    float64 // daily mean temperature, °C
	Precip   float64 // precipitation, mm
	Sunshine float64 // sunshine duration, hours
}

// RiceYield is one crop year's paddy-rice outcome for the station's prefecture,
// from the MAFF crop statistics. SakukyoIndex (作況指数) is normalized against the
// expected "normal-year" yield, so it isolates the weather signal; Yield10a is
// the absolute yield in kg per 10 ares.
type RiceYield struct {
	Year         int
	SakukyoIndex float64 // 作況指数, crop situation index (平年 = 100)
	Yield10a     float64 // 収穫期 10a当たり収量, kg per 10a
}

// YearFeature is the growth-stage risk feature vector for one crop year,
// derived from that year's daily weather.
type YearFeature struct {
	Year int

	// HeatDaysHeading counts days with TMax >= 35°C in the heading window
	// (pollination-failure risk).
	HeatDaysHeading int
	// GrainFillMeanTemp is the mean temperature over the grain-fill window
	// (chalky-grain / quality risk).
	GrainFillMeanTemp float64
	// ColdMeanTemp is the mean temperature over the cold-damage window.
	ColdMeanTemp float64
	// GrainFillSunshine is the total sunshine hours over the grain-fill window.
	GrainFillSunshine float64
}
