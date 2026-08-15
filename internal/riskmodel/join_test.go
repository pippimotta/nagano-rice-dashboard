package riskmodel

import "testing"

func TestJoin_InnerJoinByYear(t *testing.T) {
	// Only years present in both inputs survive; result is sorted ascending even
	// when the inputs are not.
	features := []YearFeature{
		{Year: 2001, HeatDaysHeading: 5, GrainFillSunshine: 210},
		{Year: 1999, HeatDaysHeading: 2, GrainFillSunshine: 190}, // no yield -> dropped
		{Year: 2000, HeatDaysHeading: 3, GrainFillSunshine: 200},
	}
	yields := []RiceYield{
		{Year: 2002, SakukyoIndex: 99, Yield10a: 606}, // no feature -> dropped
		{Year: 2000, SakukyoIndex: 104, Yield10a: 628},
		{Year: 2001, SakukyoIndex: 105, Yield10a: 633},
	}

	got := Join(features, yields)

	if len(got) != 2 {
		t.Fatalf("expected 2 joined years, got %d", len(got))
	}
	if got[0].Year != 2000 || got[1].Year != 2001 {
		t.Fatalf("years = [%d, %d], want ascending [2000, 2001]", got[0].Year, got[1].Year)
	}
	// Feature fields are preserved (embedded YearFeature).
	if got[0].HeatDaysHeading != 3 || got[0].GrainFillSunshine != 200 {
		t.Errorf("2000 feature = %+v, want HeatDaysHeading 3 / GrainFillSunshine 200", got[0].YearFeature)
	}
	// Yield fields are paired in.
	if got[0].SakukyoIndex != 104 || got[0].Yield10a != 628 {
		t.Errorf("2000 yield = (%v, %v), want (104, 628)", got[0].SakukyoIndex, got[0].Yield10a)
	}
	if got[1].SakukyoIndex != 105 || got[1].Yield10a != 633 {
		t.Errorf("2001 yield = (%v, %v), want (105, 633)", got[1].SakukyoIndex, got[1].Yield10a)
	}
}

func TestJoin_NoOverlap(t *testing.T) {
	features := []YearFeature{{Year: 2000}}
	yields := []RiceYield{{Year: 2001}}
	if got := Join(features, yields); len(got) != 0 {
		t.Errorf("expected empty join with no overlapping years, got %d", len(got))
	}
}
