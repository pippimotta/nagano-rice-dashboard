// Package dailyweather parses the raw CSV exported from the JMA "past weather
// data download" tool into domain DailyWeather records. The export is Shift-JIS
// encoded and begins with a fixed multi-line header block; empty measurement
// cells (missing observations) are represented as NaN.
package dailyweather

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"strconv"
	"time"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	"github.com/pippimotta/nagano-rice-dashboard/internal/riskmodel"
)

// Column indexes of the measurement values in a JMA data row. Each measurement
// is followed by quality/homogeneity flag columns, which v0 ignores.
const (
	colDate     = 0
	colTMax     = 1
	colTMin     = 4
	colTMean    = 7
	colPrecip   = 10
	colSunshine = 14
	minFields   = colSunshine + 1
)

const dateLayout = "2006/1/2"

// Parse reads the Shift-JIS JMA CSV from r and returns one DailyWeather per
// data row, in file order. Preamble and header lines (any line whose first
// field is not a date) are skipped. An empty measurement cell becomes NaN.
func Parse(r io.Reader) ([]riskmodel.DailyWeather, error) {
	cr := csv.NewReader(transform.NewReader(r, japanese.ShiftJIS.NewDecoder()))
	cr.FieldsPerRecord = -1 // header and data rows have different widths

	var out []riskmodel.DailyWeather
	for {
		rec, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		d, err := time.Parse(dateLayout, rec[colDate])
		if err != nil {
			continue // preamble/header line, not a data row
		}
		if len(rec) < minFields {
			return nil, fmt.Errorf("row %s: got %d fields, need at least %d", rec[colDate], len(rec), minFields)
		}

		var perr error
		val := func(col int, name string) float64 {
			v, e := parseValue(rec[col])
			if e != nil && perr == nil {
				perr = fmt.Errorf("row %s %s: %w", rec[colDate], name, e)
			}
			return v
		}
		dw := riskmodel.DailyWeather{
			Date:     d,
			TMax:     val(colTMax, "tmax"),
			TMin:     val(colTMin, "tmin"),
			TMean:    val(colTMean, "tmean"),
			Precip:   val(colPrecip, "precip"),
			Sunshine: val(colSunshine, "sunshine"),
		}
		if perr != nil {
			return nil, perr
		}
		out = append(out, dw)
	}
	return out, nil
}

// parseValue converts a measurement cell to a float. An empty cell is a missing
// observation and becomes NaN; a non-numeric cell is an error.
func parseValue(s string) (float64, error) {
	if s == "" {
		return math.NaN(), nil
	}
	return strconv.ParseFloat(s, 64)
}
