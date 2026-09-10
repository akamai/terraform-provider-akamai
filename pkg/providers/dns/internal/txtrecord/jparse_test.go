package txtrecord

import (
	"strings"
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
		in      string
		out     string
		wantErr string
	}{
		{
			in:  "AnIdentifier \"a quoted \\\" string\"\r\n; this is \"my\"\t(comment)\nanotherIdentifier (\ramultilineIdentifier\n)",
			out: "\"AnIdentifier\" \"a quoted \\\" string\" \"\\013\\010\"",
		},
		{
			in:  `"v=DKIM1; k=rsa; p=MQw7+fmMp6is3OPUL9sD/KpsauPk4gra5qsPJGtP6QVjht+Qm3lOzydHEkYE974PaxnZtGGH2wndRhL7KdinrlEhofEeq7uHXTL+yrMuQox3QiZcM+00mOLsToRJ/0i28oBtqQ2LCQCMUPo3bG8JRwFIF1nPGNP5YjCmScgRRWsY+lqY7p1PZ4Pf+/qNM3RJ818tLa5ZcO/Ae2T1gFnRTsy7iQ/xP1GUlAd+09/aSqw" "MQw7+fmMp6is3OPUL9sD/KpsauPk4gra5qsPJGtP6QVjht+Qm3lOzydHEkYE974PaxnZtGGH2wndRhL7KdinrlEhofEeq7uHXTL+yrMuQox3QiZcM+00mOLsToRJ/0i28oBtqQ2LCQCMUPo3bG8JRwFIF1nPGNP5YjCmScgRRWsY+lqY7p1PZ4Pf+/qNM3RJ818tLa5ZcO/Ae2T1gFnRTsy7iQ/xP1GUlAd+09/aSqw\010"`,
			out: "\"v=DKIM1; k=rsa; p=MQw7+fmMp6is3OPUL9sD/KpsauPk4gra5qsPJGtP6QVjht+Qm3lOzydHEkYE974PaxnZtGGH2wndRhL7KdinrlEhofEeq7uHXTL+yrMuQox3QiZcM+00mOLsToRJ/0i28oBtqQ2LCQCMUPo3bG8JRwFIF1nPGNP5YjCmScgRRWsY+lqY7p1PZ4Pf+/qNM3RJ818tLa5ZcO/Ae2T1gFnRTsy7iQ/xP1GUlAd+09/aSqw\" \"MQw7+fmMp6is3OPUL9sD/KpsauPk4gra5qsPJGtP6QVjht+Qm3lOzydHEkYE974PaxnZtGGH2wndRhL7KdinrlEhofEeq7uHXTL+yrMuQox3QiZcM+00mOLsToRJ/0i28oBtqQ2LCQCMUPo3bG8JRwFIF1nPGNP5YjCmScgRRWsY+lqY7p1PZ4Pf+/qNM3RJ818tLa5ZcO/Ae2T1gFnRTsy7iQ/xP1GUlAd+09/aSqw\\010\"",
		},
		{
			in:  "onlyOneIdentifier",
			out: "\"onlyOneIdentifier\"",
		},
		{
			in:  "identifier ;",
			out: "\"identifier\"",
		},
		{
			in:  "identifier \nidentifier2; junk comment",
			out: "\"identifier\" \"\\010identifier2\"",
		},
		{
			in:  "onetwo",
			out: "\"onetwo\"",
		},
		{
			in:  "\"one\" two",
			out: "\"one\" \"two\"",
		},
		{
			in:  "\"one\"two",
			out: "\"one\" \"two\"",
		},
		{
			in:  "\"one\" \"two\"",
			out: "\"one\" \"two\"",
		},
		{
			in:  "\"one; two\"",
			out: "\"one; two\"",
		},
		{
			in:  "one; two",
			out: "\"one\"",
		},
		{
			in:  "one\" \"two",
			out: "\"one\" \" \" \"two\"",
		},
		{
			in:  "\"one\" \" \" \"two\"",
			out: "\"one\" \" \" \"two\"",
		},
		{
			in:  "\"one\"\\\"two",
			out: "\"one\" \"\\\"two\"",
		},
		{
			in:  "\"one\" \n",
			out: "\"one\" \"\\010\"",
		},
		{
			in:  "\"one\" \"two\\010\"",
			out: "\"one\" \"two\\010\"",
		},
		{
			in:      "\"bad",
			wantErr: "EOF in quoted string",
		},
		{
			in:      ")",
			wantErr: "invalid close parenthesis",
		},
		{
			in:      "\\",
			wantErr: "unterminated escape sequence",
		},
		{
			in:      "\"\n",
			wantErr: "EOF in quoted string",
		},
		{
			in:      "(this ;",
			wantErr: "unbalanced parentheses",
		},
		{
			in:      "Hel\\lo\"world",
			wantErr: "EOF in quoted string",
		},
		{
			in:      strings.Repeat("a", 256),
			wantErr: "tokenizer exception: text string longer than 255 characters",
		},
	}

	for _, tc := range tests {
		out, err := normalizeTarget(tc.in)
		if tc.wantErr != "" {
			require.EqualError(t, err, tc.wantErr)
		} else {
			require.NoError(t, err)
		}
		if out != tc.out {
			t.Errorf("oops tc.in: %q; out: %q; tc.out: %q", tc.in, out, tc.out)
		}
	}
}
