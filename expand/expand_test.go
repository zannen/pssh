package expand_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zannen/pssh/expand"
)

func TestExpand(t *testing.T) {

	tests := []struct {
		given    string
		expected []string
	}{
		{
			given:    "a",
			expected: []string{
				"a",
			},
		},
		{
			given:    "a;b",
			expected: []string{
				"a",
				"b",
			},
		},
		{
			given:    "pre{a,b,c}post",
			expected: []string{
				"preapost",
				"prebpost",
				"precpost",
			},
		},
		{
			given:    "pre[1-3,5,7-9]post",
			expected: []string{
				"pre1post",
				"pre2post",
				"pre3post",
				"pre5post",
				"pre7post",
				"pre8post",
				"pre9post",
			},
		},
		{
			given: "a[1-2]b[3-4]c",
			expected: []string{
				"a1b3c",
				"a1b4c",
				"a2b3c",
				"a2b4c",
			},
		},
		{
			given: "a;pre{bbb[1-3],ccc[5,8]}post;z",
			expected: []string{
				"a",
				"prebbb1post",
				"prebbb2post",
				"prebbb3post",
				"preccc5post",
				"preccc8post",
				"z",
			},
		},
	}

	for _, test := range tests {
		got, err := expand.Expand(test.given)
		assert.NoError(t, err)
		assert.Equal(t, test.expected, got)
	}
}

func BenchmarkExpand(b *testing.B) {
	for n := 0; n < b.N; n++ {
		expand.Expand("a;pre{bbb[1-3],ccc[5,8]}post;z")
	}
}
