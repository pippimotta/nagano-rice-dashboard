package riskmodel

import (
	"math"
	"sort"
)

// SimilarYear is one past year ranked by climate similarity to a target year.
type SimilarYear struct {
	Year     int
	Distance float64 // standardized Euclidean distance; smaller is more similar
}

// RankSimilarYears ranks every year except target by climate similarity to
// target, closest first. Feature dimensions are standardized (z-score) across
// all years so each contributes comparably regardless of its natural scale.
// Returns an empty slice if target is absent or there are no other years.
func RankSimilarYears(target int, features []YearFeature) []SimilarYear {
	// Drop years with any missing (NaN) feature: they cannot be standardized and
	// a single one would poison the shared z-score ruler for every other year.
	usable := make([]YearFeature, 0, len(features))
	for _, f := range features {
		if !hasNaN(vector(f)) {
			usable = append(usable, f)
		}
	}
	if len(usable) < 2 {
		return []SimilarYear{}
	}

	vecs := make(map[int][]float64, len(usable))
	var targetVec []float64
	for _, f := range usable {
		v := vector(f)
		vecs[f.Year] = v
		if f.Year == target {
			targetVec = v
		}
	}
	if targetVec == nil {
		return []SimilarYear{}
	}

	means, stds := standardizeParams(usable)
	ts := standardize(targetVec, means, stds)

	out := make([]SimilarYear, 0, len(usable)-1)
	for _, f := range usable {
		if f.Year == target {
			continue
		}
		cs := standardize(vecs[f.Year], means, stds)
		out = append(out, SimilarYear{Year: f.Year, Distance: euclidean(ts, cs)})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Distance != out[j].Distance {
			return out[i].Distance < out[j].Distance
		}
		return out[i].Year < out[j].Year // stable tiebreak on equal distance
	})
	return out
}

// hasNaN reports whether any dimension is NaN, i.e. the year has a missing
// feature and is not usable for ranking.
func hasNaN(v []float64) bool {
	for _, x := range v {
		if math.IsNaN(x) {
			return true
		}
	}
	return false
}

// vector flattens a YearFeature into its ordered numeric dimensions.
func vector(f YearFeature) []float64 {
	return []float64{
		float64(f.HeatDaysHeading),
		f.GrainFillMeanTemp,
		f.ColdMeanTemp,
		f.GrainFillSunshine,
	}
}

// standardizeParams returns the per-dimension mean and (population) standard
// deviation across all years.
func standardizeParams(features []YearFeature) (means, stds []float64) {
	dims := len(vector(features[0]))
	means = make([]float64, dims)
	stds = make([]float64, dims)
	n := float64(len(features))

	for _, f := range features {
		for i, x := range vector(f) {
			means[i] += x
		}
	}
	for i := range means {
		means[i] /= n
	}
	for _, f := range features {
		for i, x := range vector(f) {
			d := x - means[i]
			stds[i] += d * d
		}
	}
	for i := range stds {
		stds[i] = math.Sqrt(stds[i] / n)
	}
	return means, stds
}

// standardize converts a raw vector to z-scores. A constant dimension (std 0)
// carries no information and contributes 0.
func standardize(v, means, stds []float64) []float64 {
	out := make([]float64, len(v))
	for i := range v {
		if stds[i] == 0 {
			continue
		}
		out[i] = (v[i] - means[i]) / stds[i]
	}
	return out
}

func euclidean(a, b []float64) float64 {
	var sum float64
	for i := range a {
		d := a[i] - b[i]
		sum += d * d
	}
	return math.Sqrt(sum)
}
