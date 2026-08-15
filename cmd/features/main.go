// Command features builds the derived year_features dataset: it reads the raw
// JMA weather and MAFF yield fixtures, computes each year's growth-stage risk
// features, joins them with the rice outcome, and writes a static JSON file for
// the dashboard. Run from the repository root.
package main

import (
	"encoding/json"
	"log"
	"math"
	"os"
	"path/filepath"

	"github.com/pippimotta/nagano-rice-dashboard/internal/dailyweather"
	"github.com/pippimotta/nagano-rice-dashboard/internal/riceyield"
	"github.com/pippimotta/nagano-rice-dashboard/internal/riskmodel"
)

const (
	weatherPath = "data/raw/nagano_daily_1996-2026.csv"
	yieldPath   = "data/raw/estat_nagano_yield.json"
	outPath     = "web/data/year_features.json"
)

// yearJSON is the export shape for one year. Floats that may be missing use a
// pointer so a NaN feature is written as JSON null rather than failing to encode.
type yearJSON struct {
	Year              int      `json:"year"`
	HeatDaysHeading   int      `json:"heatDaysHeading"`
	GrainFillMeanTemp *float64 `json:"grainFillMeanTemp"`
	ColdMeanTemp      *float64 `json:"coldMeanTemp"`
	GrainFillSunshine *float64 `json:"grainFillSunshine"`
	SakukyoIndex      *float64 `json:"sakukyoIndex"`
	Yield10a          *float64 `json:"yield10a"`
}

func main() {
	weather, err := parseWeather(weatherPath)
	if err != nil {
		log.Fatalf("weather: %v", err)
	}
	yields, err := parseYields(yieldPath)
	if err != nil {
		log.Fatalf("yield: %v", err)
	}

	records := riskmodel.Join(riskmodel.ComputeYearFeatures(weather), yields)

	out := make([]yearJSON, 0, len(records))
	for _, r := range records {
		out = append(out, yearJSON{
			Year:              r.Year,
			HeatDaysHeading:   r.HeatDaysHeading,
			GrainFillMeanTemp: nullable(r.GrainFillMeanTemp),
			ColdMeanTemp:      nullable(r.ColdMeanTemp),
			GrainFillSunshine: nullable(r.GrainFillSunshine),
			SakukyoIndex:      nullable(r.SakukyoIndex),
			Yield10a:          nullable(r.Yield10a),
		})
	}

	if err := writeJSON(outPath, out); err != nil {
		log.Fatalf("write: %v", err)
	}
	log.Printf("wrote %d years (%d..%d) to %s", len(out), out[0].Year, out[len(out)-1].Year, outPath)
}

func parseWeather(path string) ([]riskmodel.DailyWeather, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return dailyweather.Parse(f)
}

func parseYields(path string) ([]riskmodel.RiceYield, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return riceyield.Parse(f)
}

func writeJSON(path string, v any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0o644)
}

// nullable maps a NaN (missing) value to nil so it encodes as JSON null.
func nullable(f float64) *float64 {
	if math.IsNaN(f) {
		return nil
	}
	return &f
}
