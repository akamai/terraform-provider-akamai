package txtrecord

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeTarget(t *testing.T) {
	tests := []struct {
		in        string
		expected  string
		withError error
	}{
		{
			in:       "Hel\\lo\"world",
			expected: "\"Hel\\\\lo\\\"world\"",
		},
		{
			in:       "\"Hel\\\\lo\\\"world\"",
			expected: "\"Hel\\\\lo\\\"world\"",
		},
	}

	for _, tc := range tests {
		res, err := NormalizeTarget(tc.in)
		if tc.withError != nil {
			assert.Error(t, err)
			continue
		}

		require.NoError(t, err)
		assert.Equal(t, tc.expected, res)
	}
}

func Test_normalizeTarget(t *testing.T) {
	tests := []struct {
		in  string
		out string
		ok  bool
	}{
		{
			in:  "AnIdentifier \"a quoted \\\" string\"\r\n; this is \"my\"\t(comment)\nanotherIdentifier (\ramultilineIdentifier\n)",
			out: "\"AnIdentifier\" \"a quoted \\\" string\" \"\\013\\010\"",
			ok:  true,
		},
		{
			in:  `"v=DKIM1; k=rsa; p=MQw7+fmMp6is3OPUL9sD/KpsauPk4gra5qsPJGtP6QVjht+Qm3lOzydHEkYE974PaxnZtGGH2wndRhL7KdinrlEhofEeq7uHXTL+yrMuQox3QiZcM+00mOLsToRJ/0i28oBtqQ2LCQCMUPo3bG8JRwFIF1nPGNP5YjCmScgRRWsY+lqY7p1PZ4Pf+/qNM3RJ818tLa5ZcO/Ae2T1gFnRTsy7iQ/xP1GUlAd+09/aSqw" "MQw7+fmMp6is3OPUL9sD/KpsauPk4gra5qsPJGtP6QVjht+Qm3lOzydHEkYE974PaxnZtGGH2wndRhL7KdinrlEhofEeq7uHXTL+yrMuQox3QiZcM+00mOLsToRJ/0i28oBtqQ2LCQCMUPo3bG8JRwFIF1nPGNP5YjCmScgRRWsY+lqY7p1PZ4Pf+/qNM3RJ818tLa5ZcO/Ae2T1gFnRTsy7iQ/xP1GUlAd+09/aSqw\010"`,
			out: "\"v=DKIM1; k=rsa; p=MQw7+fmMp6is3OPUL9sD/KpsauPk4gra5qsPJGtP6QVjht+Qm3lOzydHEkYE974PaxnZtGGH2wndRhL7KdinrlEhofEeq7uHXTL+yrMuQox3QiZcM+00mOLsToRJ/0i28oBtqQ2LCQCMUPo3bG8JRwFIF1nPGNP5YjCmScgRRWsY+lqY7p1PZ4Pf+/qNM3RJ818tLa5ZcO/Ae2T1gFnRTsy7iQ/xP1GUlAd+09/aSqw\" \"MQw7+fmMp6is3OPUL9sD/KpsauPk4gra5qsPJGtP6QVjht+Qm3lOzydHEkYE974PaxnZtGGH2wndRhL7KdinrlEhofEeq7uHXTL+yrMuQox3QiZcM+00mOLsToRJ/0i28oBtqQ2LCQCMUPo3bG8JRwFIF1nPGNP5YjCmScgRRWsY+lqY7p1PZ4Pf+/qNM3RJ818tLa5ZcO/Ae2T1gFnRTsy7iQ/xP1GUlAd+09/aSqw\\010\"",
			ok:  true,
		},
		{
			in:  "onlyOneIdentifier",
			out: "\"onlyOneIdentifier\"",
			ok:  true,
		},
		{
			in:  "identifier ;",
			out: "\"identifier\"",
			ok:  true,
		},
		{
			in:  "identifier \nidentifier2; junk comment",
			out: "\"identifier\" \"\\010identifier2\"",
			ok:  true,
		},
		{
			in:  "onetwo",
			out: "\"onetwo\"",
			ok:  true,
		},
		{
			in:  "\"one\" two",
			out: "\"one\" \"two\"",
			ok:  true,
		},
		{
			in:  "\"one\"two",
			out: "\"one\" \"two\"",
			ok:  true,
		},
		{
			in:  "\"one\" \"two\"",
			out: "\"one\" \"two\"",
			ok:  true,
		},
		{
			in:  "\"one; two\"",
			out: "\"one; two\"",
			ok:  true,
		},
		{
			in:  "one; two",
			out: "\"one\"",
			ok:  true,
		},
		{
			in:  "one\" \"two",
			out: "\"one\" \" \" \"two\"",
			ok:  true,
		},
		{
			in:  "\"one\" \" \" \"two\"",
			out: "\"one\" \" \" \"two\"",
			ok:  true,
		},
		{
			in:  "\"one\"\\\"two",
			out: "\"one\" \"\\\"two\"",
			ok:  true,
		},
		{
			in:  "\"one\" \n",
			out: "\"one\" \"\\010\"",
			ok:  true,
		},
		{
			in:  "\"one\" \"two\\010\"",
			out: "\"one\" \"two\\010\"",
			ok:  true,
		},
		{
			in: "\"bad",
			ok: false,
		},
		{
			in: ")",
			ok: false,
		},
		{
			in: "\\",
			ok: false,
		},
		{
			in: "\"\n",
			ok: false,
		},
		{
			in: "(this ;",
			ok: false,
		},
		{
			in: "Hel\\lo\"world",
			ok: false,
		},
	}

	for _, tc := range tests {
		out, ok := normalizeTarget(tc.in)
		if ok != tc.ok {
			t.Errorf("oops tc.in: %q; ok: %v", tc.in, ok)
		}
		if out != tc.out {
			t.Errorf("oops tc.in: %q; out: %q; tc.out: %q", tc.in, out, tc.out)
		}
	}
}
