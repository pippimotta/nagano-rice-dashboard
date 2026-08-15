package riskmodel

import "time"

// growthWindow is an inclusive calendar date range within a single year.
// Windows are fixed (not variety- or elevation-adjusted) in v0.
type growthWindow struct {
	startMonth time.Month
	startDay   int
	endMonth   time.Month
	endDay     int
}

// contains reports whether t's month/day falls within the window (inclusive).
// Windows never wrap the year end, so comparing a month*100+day ordinal is enough.
func (w growthWindow) contains(t time.Time) bool {
	md := int(t.Month())*100 + t.Day()
	start := int(w.startMonth)*100 + w.startDay
	end := int(w.endMonth)*100 + w.endDay
	return md >= start && md <= end
}

// Growth-stage risk windows for Nagano Koshihikari (fixed calendar, v0).
var (
	// headingWindow covers heading/pollination — high-temperature damage risk.
	headingWindow = growthWindow{time.August, 1, time.August, 15}
	// grainFillWindow covers grain-filling — quality (chalky grain) and
	// sunshine risk.
	grainFillWindow = growthWindow{time.August, 10, time.September, 15}
	// coldWindow covers the cold-damage (障害型冷害) sensitive period.
	coldWindow = growthWindow{time.July, 20, time.August, 20}
)
