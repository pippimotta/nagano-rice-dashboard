package riskmodel

import "sort"

// heatThresholdC is the daily maximum temperature at or above which a heading-window
// day counts as a high-temperature (pollination-failure) risk day.
const heatThresholdC = 35.0

// ComputeYearFeatures groups daily weather by calendar year and computes each
// year's growth-stage risk features. Years are returned sorted ascending.
func ComputeYearFeatures(weather []DailyWeather) []YearFeature {
	type acc struct {
		heatDays          int
		grainFillTempSum  float64
		grainFillTempN    int
		coldTempSum       float64
		coldTempN         int
		grainFillSunshine float64
	}
	byYear := map[int]*acc{}
	for _, d := range weather {
		a := byYear[d.Date.Year()]
		if a == nil {
			a = &acc{}
			byYear[d.Date.Year()] = a
		}
		if headingWindow.contains(d.Date) && d.TMax >= heatThresholdC {
			a.heatDays++
		}
		if grainFillWindow.contains(d.Date) {
			a.grainFillTempSum += d.TMean
			a.grainFillTempN++
			a.grainFillSunshine += d.Sunshine
		}
		if coldWindow.contains(d.Date) {
			a.coldTempSum += d.TMean
			a.coldTempN++
		}
	}

	years := make([]int, 0, len(byYear))
	for y := range byYear {
		years = append(years, y)
	}
	sort.Ints(years)

	out := make([]YearFeature, 0, len(years))
	for _, y := range years {
		a := byYear[y]
		out = append(out, YearFeature{
			Year:              y,
			HeatDaysHeading:   a.heatDays,
			GrainFillMeanTemp: mean(a.grainFillTempSum, a.grainFillTempN),
			ColdMeanTemp:      mean(a.coldTempSum, a.coldTempN),
			GrainFillSunshine: a.grainFillSunshine,
		})
	}
	return out
}

// mean returns sum/n, or 0 when n is 0 (a window with no observed days).
func mean(sum float64, n int) float64 {
	if n == 0 {
		return 0
	}
	return sum / float64(n)
}
