package hw02unpackstring

import (
	"errors"
	"reflect"
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
	// Unicode tests
	{input: "Привет!", expected: "Привет!"},
	{input: "П2р3ивет4!", expected: "ППррриветттт!"},
	{input: "П", expected: "П"},
	{input: `\\0Прив4ет\35`, expected: `Приввввет33333`},
	{input: `😅5😅`, expected: `😅😅😅😅😅😅`},
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
	invalidStrings := []string{
		"3abc",
		"45",
		"aaa10b",
		"п19ривет",
		`\Привет`,
		`\123ы`,
	}
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
		input, expected []rune
	}{
		{[]rune("a3b4c5"), []rune("a3")},
		{[]rune("b4c5"), []rune("b4")},
		{[]rune("c5"), []rune("c5")},
		{[]rune("s"), []rune("s")},
		{[]rune("ab5cc3a"), []rune("a")},
		// tasks with asterisk
		{[]rune(`\4`), []rune(`\4`)},
		{[]rune(`\\`), []rune(`\\`)},
		{[]rune(`\45a5`), []rune(`\45`)},
		{[]rune(`\\5bca`), []rune(`\\5`)},
		{[]rune(`\\\5a5`), []rune(`\\`)},
		// Unicode tests
		{[]rune(`пр2иве5т`), []rune(`п`)},
		{[]rune(`р2иве5т`), []rune(`р2`)},
		{[]rune(`иве5т`), []rune(`и`)},
		{[]rune(`е5т`), []rune(`е5`)},
	}
	for _, tc := range cases {
		t.Run(
			string(tc.input), func(t *testing.T) {
				result, err := findNextSubstr(tc.input)
				require.NoError(t, err)
				require.True(t, reflect.DeepEqual(tc.expected, result))
			})
	}
}

func TestFindNextSubstrInvalidString(t *testing.T) {
	invalidStrings := [][]rune{
		[]rune("4"),
		[]rune("3a4"),
		[]rune("45cc3"),
		[]rune("a10b"),
		[]rune("a45"),
		[]rune(`\`),
		[]rune(`\n5`),
		[]rune(`\П5`),
		[]rune("П55"),
	}
	for _, tc := range invalidStrings {
		t.Run(string(tc), func(t *testing.T) {
			_, err := findNextSubstr(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}

func TestUnpackSubstr(t *testing.T) {
	cases := []struct {
		input    []rune
		expected string
	}{
		{[]rune("a3"), "aaa"},
		{[]rune("b"), "b"},
		{[]rune("c5"), "ccccc"},
		// tasks with asterisk
		{[]rune(`\\`), `\`},
		{[]rune(`\\5`), `\\\\\`},
		{[]rune(`\4`), "4"},
		{[]rune(`\42`), "44"},
		{[]rune(`ё5`), "ёёёёё"},
		{[]rune(`ё`), "ё"},
	}
	for _, tc := range cases {
		t.Run(
			string(tc.input), func(t *testing.T) {
				result := unpackSubstr(tc.input)
				require.Equal(t, tc.expected, result)
			})
	}
}
