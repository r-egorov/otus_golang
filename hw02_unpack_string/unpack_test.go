package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

var testCasesValid = []struct {
	input    string
	expected string
}{
	{input: "a4bc2d5e", expected: "aaaabccddddde"},
	{input: "d\n5abc", expected: "d\n\n\n\n\nabc"},
	{input: "abccd", expected: "abccd"},
	{input: "", expected: ""},
	{input: "aaa0b", expected: "aab"},
	// uncomment if task with asterisk completed
	{input: `qwe\4\5`, expected: `qwe45`},
	{input: `qwe\45`, expected: `qwe44444`},
	{input: `qwe\\5`, expected: `qwe\\\\\`},
	{input: `qwe\\\3`, expected: `qwe\3`},
}

func TestUnpack(t *testing.T) {
	for _, tc := range testCasesValid {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b"}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}

func BenchmarkUnpack(b *testing.B) {
	for _, tc := range testCasesValid {
		for i := 0; i < b.N; i++ {
			if _, err := Unpack(tc.input); err != nil {
				b.Fatalf("Error: %q", err)
			}
		}
	}
}

func TestFindNextSubstr(t *testing.T) {
	cases := []struct {
		input, expected string
	}{
		{"a3b4c5", "a3"},
		{"b4c5", "b4"},
		{"c5", "c5"},
		{"s", "s"},
		{"ab5cc3a", "a"},
		// tasks with asterisk
		{`\4`, `\4`},
		{`\\`, `\\`},
		{`\45a5`, `\45`},
		{`\\5bca`, `\\5`},
		{`\\\5a5`, `\\`},
	}
	for _, tc := range cases {
		t.Run(
			tc.input, func(t *testing.T) {
				result, err := findNextSubstr(tc.input)
				require.NoError(t, err)
				require.Equal(t, tc.expected, result)
			})
	}
}

func TestFindNextSubstrInvalidString(t *testing.T) {
	invalidStrings := []string{
		"4",
		"3a4",
		"45cc3",
		"a10b",
		"a45",
		`\`,
		`\n5`,
	}
	for _, tc := range invalidStrings {
		t.Run(tc, func(t *testing.T) {
			_, err := findNextSubstr(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}

func TestUnpackSubstr(t *testing.T) {
	cases := []struct {
		input, expected string
	}{
		{"a3", "aaa"},
		{"b", "b"},
		{"c5", "ccccc"},
		// tasks with asterisk
		{`\\`, `\`},
		{`\\5`, `\\\\\`},
		{`\4`, "4"},
		{`\42`, "44"},
	}
	for _, tc := range cases {
		t.Run(
			tc.input, func(t *testing.T) {
				result := unpackSubstr(tc.input)
				require.Equal(t, tc.expected, result)
			})
	}
}
