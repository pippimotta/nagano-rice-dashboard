package dailyweather

import (
	"math"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/pippimotta/nagano-rice-dashboard/internal/riskmodel"
)

// sameFloat compares two float64s, treating NaN as equal to NaN (a missing
// value is expected to be represented as NaN).
func sameFloat(a, b float64) bool {
	if math.IsNaN(a) && math.IsNaN(b) {
		return true
	}
	return a == b
}

func TestParse_Sample(t *testing.T) {
	f, err := os.Open("testdata/sample.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	got, err := Parse(f)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	nan := math.NaN()
	want := []riskmodel.DailyWeather{
		{Date: date(1997, time.August, 1), TMax: 35.2, TMin: 24.0, TMean: 29.0, Precip: 0.0, Sunshine: 10.5},
		{Date: date(1997, time.August, 2), TMax: 28.0, TMin: 20.0, TMean: 24.0, Precip: 12.5, Sunshine: nan}, // sunshine missing
		{Date: date(1997, time.August, 3), TMax: nan, TMin: 19.0, TMean: 23.0, Precip: 3.0, Sunshine: 7.0},   // tmax missing
		{Date: date(1997, time.August, 4), TMax: 30.0, TMin: 21.5, TMean: 25.0, Precip: 0.0, Sunshine: 8.5},
	}
	if len(got) != len(want) {
		t.Fatalf("parsed %d rows, want %d", len(got), len(want))
	}
	for i := range want {
		g, w := got[i], want[i]
		if !g.Date.Equal(w.Date) {
			t.Errorf("row %d Date = %v, want %v", i, g.Date, w.Date)
		}
		if !sameFloat(g.TMax, w.TMax) {
			t.Errorf("row %d TMax = %v, want %v", i, g.TMax, w.TMax)
		}
		if !sameFloat(g.TMin, w.TMin) {
			t.Errorf("row %d TMin = %v, want %v", i, g.TMin, w.TMin)
		}
		if !sameFloat(g.TMean, w.TMean) {
			t.Errorf("row %d TMean = %v, want %v", i, g.TMean, w.TMean)
		}
		if !sameFloat(g.Precip, w.Precip) {
			t.Errorf("row %d Precip = %v, want %v", i, g.Precip, w.Precip)
		}
		if !sameFloat(g.Sunshine, w.Sunshine) {
			t.Errorf("row %d Sunshine = %v, want %v", i, g.Sunshine, w.Sunshine)
		}
	}
}

func TestParse_ShortDataRow(t *testing.T) {
	// A row whose first field is a date but that has too few columns is
	// malformed and must be reported, not silently truncated.
	in := "header line\r\n2020/1/1,10,8,1\r\n"
	if _, err := Parse(strings.NewReader(in)); err == nil {
		t.Fatal("expected error for a short data row, got nil")
	}
}

func date(y int, m time.Month, d int) time.Time {
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}
