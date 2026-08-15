package riskmodel

import "testing"

func TestRankSimilarYears_ClosestFirstAndExcludesTarget(t *testing.T) {
	// Only HeatDaysHeading varies; the other dimensions are constant, so the
	// ranking is driven by distance in that one dimension and is unambiguous.
	features := []YearFeature{
		{Year: 2000, HeatDaysHeading: 5, GrainFillMeanTemp: 25, ColdMeanTemp: 22, GrainFillSunshine: 600},
		{Year: 2001, HeatDaysHeading: 5, GrainFillMeanTemp: 25, ColdMeanTemp: 22, GrainFillSunshine: 600}, // identical to target
		{Year: 2002, HeatDaysHeading: 10, GrainFillMeanTemp: 25, ColdMeanTemp: 22, GrainFillSunshine: 600},
		{Year: 2003, HeatDaysHeading: 15, GrainFillMeanTemp: 25, ColdMeanTemp: 22, GrainFillSunshine: 600},
	}
	got := RankSimilarYears(2000, features)

	if len(got) != 3 {
		t.Fatalf("expected 3 ranked years (target excluded), got %d", len(got))
	}
	want := []int{2001, 2002, 2003}
	for i, w := range want {
		if got[i].Year != w {
			t.Errorf("rank %d = year %d, want %d", i, got[i].Year, w)
		}
	}
	for _, s := range got {
		if s.Year == 2000 {
			t.Errorf("target year 2000 must not appear in its own ranking")
		}
	}
	if got[0].Distance != 0 {
		t.Errorf("identical year distance = %v, want 0", got[0].Distance)
	}
}

func TestRankSimilarYears_TargetNotPresent(t *testing.T) {
	features := []YearFeature{
		{Year: 2000, HeatDaysHeading: 5},
		{Year: 2001, HeatDaysHeading: 6},
	}
	if got := RankSimilarYears(1999, features); len(got) != 0 {
		t.Errorf("expected empty ranking when target absent, got %d", len(got))
	}
}

func TestRankSimilarYears_SingleYear(t *testing.T) {
	features := []YearFeature{{Year: 2000, HeatDaysHeading: 5}}
	if got := RankSimilarYears(2000, features); len(got) != 0 {
		t.Errorf("expected empty ranking with no other years, got %d", len(got))
	}
}
