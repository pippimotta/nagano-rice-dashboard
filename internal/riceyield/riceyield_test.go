package riceyield

import (
	"math"
	"os"
	"strings"
	"testing"

	"github.com/pippimotta/nagano-rice-dashboard/internal/riskmodel"
)

func sameFloat(a, b float64) bool {
	if math.IsNaN(a) && math.IsNaN(b) {
		return true
	}
	return a == b
}

func TestParse_Sample(t *testing.T) {
	f, err := os.Open("testdata/sample.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	got, err := Parse(f)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	want := []riskmodel.RiceYield{
		{Year: 2016, SakukyoIndex: 99, Yield10a: 596},
		{Year: 2017, SakukyoIndex: math.NaN(), Yield10a: 620}, // sakukyo "…" -> NaN
	}
	if len(got) != len(want) {
		t.Fatalf("parsed %d records, want %d", len(got), len(want))
	}
	for i := range want {
		g, w := got[i], want[i]
		if g.Year != w.Year {
			t.Errorf("record %d Year = %d, want %d", i, g.Year, w.Year)
		}
		if !sameFloat(g.SakukyoIndex, w.SakukyoIndex) {
			t.Errorf("record %d SakukyoIndex = %v, want %v", i, g.SakukyoIndex, w.SakukyoIndex)
		}
		if !sameFloat(g.Yield10a, w.Yield10a) {
			t.Errorf("record %d Yield10a = %v, want %v", i, g.Yield10a, w.Yield10a)
		}
	}
}

func TestParse_InvalidJSON(t *testing.T) {
	if _, err := Parse(strings.NewReader("{not json")); err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}
