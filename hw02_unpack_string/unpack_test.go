package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		// uncomment if task with asterisk completed
		// {input: `qwe\4\5`, expected: `qwe45`},
		// {input: `qwe\45`, expected: `qwe44444`},
		// {input: `qwe\\5`, expected: `qwe\\\\\`},
		// {input: `qwe\\\3`, expected: `qwe\3`},
	}

	for _, tc := range tests {
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

func TestFindNextSubstr(t *testing.T) {
	cases := []struct {
		inputStr, want string
	}{
		{"a3b4c5", "a3"},
		{"b4c5", "b4"},
		{"c5", "c5"},
		{"s", "s"},
		{"ab5cc3a", "a"},
	}
	for _, tc := range cases {
		t.Run(
			tc.inputStr, func(t *testing.T) {
				got, err := findNextSubstr(tc.inputStr)
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
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
		inputStr, want string
	}{
		{"a3", "aaa"},
		{"b", "b"},
		{"c5", "ccccc"},
	}
	for _, tc := range cases {
		t.Run(
			tc.inputStr, func(t *testing.T) {
				got := unpackSubstr(tc.inputStr)
				require.Equal(t, tc.want, got)
			})
	}
}
