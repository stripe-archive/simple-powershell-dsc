package urls

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertRegexMatch(t *testing.T, re, str string, expected bool) {
	r := regexp.MustCompile(re)
	assert.Equal(t, expected, r.MatchString(str))
}

func TestTypes(t *testing.T) {
	tcs := []struct {
		Re       string
		Str      string
		Expected bool
	}{
		{ModuleVersion, `1.0`, true},
		{ModuleVersion, `1.2.3`, true},
		{ModuleVersion, `1.2.3.4`, true},
		{ModuleVersion, ``, true},
		{ModuleVersion, `1`, false},

		{ConfigurationId, `B1F28971-2CEB-46D5-9DCB-79C044395F81`, true},
		{ConfigurationId, `B1F28971-2CEB-46D5-9DCB-79C044395F8`, false},

		{ModuleName, `random_m0dule`, true},
		{ModuleName, `invalid!module`, false},
	}

	for i, tc := range tcs {
		r, err := regexp.Compile(`^` + tc.Re + `$`)
		if assert.NoError(t, err) {
			match := r.MatchString(tc.Str)
			assert.Equal(t, tc.Expected, match,
				"test case %d: expected %v, got %v",
				i, tc.Expected, match)
		}
	}
}
