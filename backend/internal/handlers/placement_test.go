package handlers

import "testing"

func TestPlacementLevel(t *testing.T) {
	cases := []struct {
		correct int
		expect  string
	}{
		{0, "base"},
		{3, "base"},
		{4, "mid"},
		{7, "mid"},
		{8, "pro"},
		{10, "pro"},
	}
	for _, tc := range cases {
		if got := placementLevel(tc.correct); got != tc.expect {
			t.Fatalf("correct=%d expected %s got %s", tc.correct, tc.expect, got)
		}
	}
}
