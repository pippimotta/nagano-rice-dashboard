package riskmodel

import (
	"testing"
	"time"
)

func day(year int, m time.Month, d int) time.Time {
	return time.Date(year, m, d, 0, 0, 0, 0, time.UTC)
}

// one puts all daily records under a single year and returns that year's
// feature, failing if the number of years is not exactly one.
func one(t *testing.T, weather []DailyWeather) YearFeature {
	t.Helper()
	got := ComputeYearFeatures(weather)
	if len(got) != 1 {
		t.Fatalf("expected features for exactly 1 year, got %d", len(got))
	}
	return got[0]
}

func TestComputeYearFeatures_HeatDaysHeading(t *testing.T) {
	// Heading window is Aug 1–15; heat day = TMax >= 35°C inside it.
	weather := []DailyWeather{
		{Date: day(2020, time.August, 1), TMax: 36},  // in window, hot -> count
		{Date: day(2020, time.August, 15), TMax: 35}, // boundary, exactly 35 -> count
		{Date: day(2020, time.August, 3), TMax: 34},  // in window, not hot -> no
		{Date: day(2020, time.July, 30), TMax: 40},   // before window -> no
		{Date: day(2020, time.August, 16), TMax: 40}, // after window -> no
	}
	if got := one(t, weather).HeatDaysHeading; got != 2 {
		t.Errorf("HeatDaysHeading = %d, want 2", got)
	}
}

func TestComputeYearFeatures_GrainFillMeanTempAndSunshine(t *testing.T) {
	// Use September dates: only the grain-fill window (Aug 10 – Sep 15) covers
	// them, so cold/heading windows do not interfere.
	weather := []DailyWeather{
		{Date: day(2020, time.September, 1), TMean: 20, Sunshine: 5},
		{Date: day(2020, time.September, 10), TMean: 24, Sunshine: 7},
		{Date: day(2020, time.September, 20), TMean: 99, Sunshine: 99}, // after window -> ignored
	}
	f := one(t, weather)
	if f.GrainFillMeanTemp != 22 {
		t.Errorf("GrainFillMeanTemp = %v, want 22", f.GrainFillMeanTemp)
	}
	if f.GrainFillSunshine != 12 {
		t.Errorf("GrainFillSunshine = %v, want 12", f.GrainFillSunshine)
	}
}

func TestComputeYearFeatures_ColdMeanTemp(t *testing.T) {
	// Use late-July dates: only the cold window (Jul 20 – Aug 20) covers them.
	weather := []DailyWeather{
		{Date: day(2020, time.July, 20), TMean: 18},
		{Date: day(2020, time.July, 25), TMean: 20},
		{Date: day(2020, time.July, 10), TMean: 99}, // before window -> ignored
	}
	if got := one(t, weather).ColdMeanTemp; got != 19 {
		t.Errorf("ColdMeanTemp = %v, want 19", got)
	}
}

func TestComputeYearFeatures_GroupsByYearSorted(t *testing.T) {
	weather := []DailyWeather{
		{Date: day(2021, time.August, 2), TMax: 36},
		{Date: day(2020, time.August, 2), TMax: 36},
	}
	got := ComputeYearFeatures(weather)
	if len(got) != 2 {
		t.Fatalf("expected 2 years, got %d", len(got))
	}
	if got[0].Year != 2020 || got[1].Year != 2021 {
		t.Errorf("years = [%d, %d], want ascending [2020, 2021]", got[0].Year, got[1].Year)
	}
}

func TestComputeYearFeatures_Empty(t *testing.T) {
	if got := ComputeYearFeatures(nil); len(got) != 0 {
		t.Errorf("expected no features for empty input, got %d", len(got))
	}
}
