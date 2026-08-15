// Package riceyield parses the JSON returned by the e-Stat getStatsData API for
// the MAFF "収穫量累年統計 水稲" table (a single prefecture, year x metric) into
// domain RiceYield records. A non-numeric cell (e.g. e-Stat's "…" for missing)
// is represented as NaN.
package riceyield

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"regexp"
	"strconv"

	"github.com/pippimotta/nagano-rice-dashboard/internal/riskmodel"
)

// cat02 metric codes in the 収穫量累年統計 table.
const (
	metricYield10a = "1016" // 収穫期_10a当たり収量
	metricSakukyo  = "1017" // 収穫期_作況指数
)

// yearInName pulls the Western year out of an e-Stat time label, e.g.
// "平.28(2016)" -> 2016.
var yearInName = regexp.MustCompile(`\((\d{4})\)`)

// statsData mirrors the subset of the getStatsData JSON that this parser reads.
type statsData struct {
	GetStatsData struct {
		StatisticalData struct {
			ClassInf struct {
				ClassObj []struct {
					ID    string `json:"@id"`
					Class []struct {
						Code string `json:"@code"`
						Name string `json:"@name"`
					} `json:"CLASS"`
				} `json:"CLASS_OBJ"`
			} `json:"CLASS_INF"`
			DataInf struct {
				Value []struct {
					Cat01 string `json:"@cat01"`
					Cat02 string `json:"@cat02"`
					Val   string `json:"$"`
				} `json:"VALUE"`
			} `json:"DATA_INF"`
		} `json:"STATISTICAL_DATA"`
	} `json:"GET_STATS_DATA"`
}

// Parse reads an e-Stat getStatsData JSON response and returns one RiceYield per
// year in the table, sorted ascending. A non-numeric measurement becomes NaN.
func Parse(r io.Reader) ([]riskmodel.RiceYield, error) {
	var sd statsData
	if err := json.NewDecoder(r).Decode(&sd); err != nil {
		return nil, err
	}

	// Map each year code (cat01) to its Western year, in table order.
	var yearCodes []string
	yearByCode := map[string]int{}
	for _, obj := range sd.GetStatsData.StatisticalData.ClassInf.ClassObj {
		if obj.ID != "cat01" {
			continue
		}
		for _, c := range obj.Class {
			m := yearInName.FindStringSubmatch(c.Name)
			if m == nil {
				return nil, fmt.Errorf("cat01 %q: no year in name %q", c.Code, c.Name)
			}
			y, err := strconv.Atoi(m[1])
			if err != nil {
				return nil, fmt.Errorf("cat01 %q: %w", c.Code, err)
			}
			yearCodes = append(yearCodes, c.Code)
			yearByCode[c.Code] = y
		}
	}

	// Index the metric values we want by year code.
	byYear := map[string]map[string]string{}
	for _, v := range sd.GetStatsData.StatisticalData.DataInf.Value {
		if v.Cat02 != metricYield10a && v.Cat02 != metricSakukyo {
			continue
		}
		if byYear[v.Cat01] == nil {
			byYear[v.Cat01] = map[string]string{}
		}
		byYear[v.Cat01][v.Cat02] = v.Val
	}

	out := make([]riskmodel.RiceYield, 0, len(yearCodes))
	for _, code := range yearCodes {
		m := byYear[code]
		out = append(out, riskmodel.RiceYield{
			Year:         yearByCode[code],
			SakukyoIndex: parseValue(m[metricSakukyo]),
			Yield10a:     parseValue(m[metricYield10a]),
		})
	}
	return out, nil
}

// parseValue converts an e-Stat cell to a float. Anything non-numeric (missing
// markers like "…", or an absent cell) becomes NaN.
func parseValue(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return math.NaN()
	}
	return f
}
