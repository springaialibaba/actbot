package retest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRetestCommentBodyMatch(t *testing.T) {
	cases := []struct {
		caseName string
		comment  string
		expect   bool
	}{
		{
			caseName: "Match the retest instruction",
			comment:  "/retest",
			expect:   true,
		},
		{
			caseName: "Match the instructions that show multiple spaces after retest",
			comment:  "/retest    ",
			expect:   true,
		},
		{
			caseName: "unmatched instructions",
			comment:  "/redo",
			expect:   false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.caseName, func(t *testing.T) {
			assert.Equal(t, tc.expect, retestRegexp.MatchString(tc.comment))
		})
	}
}
