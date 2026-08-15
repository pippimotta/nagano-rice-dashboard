package riskmodel

import "sort"

// Join inner-joins per-year features with per-year rice yields on Year, keeping
// only years present in both inputs. Results are sorted ascending by year.
func Join(features []YearFeature, yields []RiceYield) []YearRecord {
	yieldByYear := make(map[int]RiceYield, len(yields))
	for _, y := range yields {
		yieldByYear[y.Year] = y
	}

	out := make([]YearRecord, 0)
	for _, f := range features {
		y, ok := yieldByYear[f.Year]
		if !ok {
			continue
		}
		out = append(out, YearRecord{
			YearFeature:  f,
			SakukyoIndex: y.SakukyoIndex,
			Yield10a:     y.Yield10a,
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Year < out[j].Year })
	return out
}
