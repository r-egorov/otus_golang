package rprint_test

import (
	"bytes"
	"testing"

	"github.com/r-egorov/otus_golang/hw01_hello_otus/rprint"
)

func TestRevPrint(t *testing.T) {
	assertPrinted := func(tb testing.TB, got, want string) {
		tb.Helper()
		if got != want {
			tb.Errorf("got %q want %q", got, want)
		}
	}

	assertNoError := func(tb testing.TB, err error) {
		tb.Helper()
		if err != nil {
			tb.Errorf("got error %v", err)
		}
	}

	getOutput := func(tb testing.TB, message string) string {
		tb.Helper()

		out := &bytes.Buffer{}
		err := rprint.RevPrint(out, message)
		assertNoError(tb, err)
		return out.String()
	}

	cases := []struct {
		name, messageToPrint, want string
	}{
		{"hello", "Hello, world", "dlrow ,olleH"},
		{"actual work", "Hello, OTUS!", "!SUTO ,olleH"},
		{"palindrome", "около м иши м олоко", "около м иши м олоко"},
		{"empty string", "", ""},
	}

	for _, testcase := range cases {
		t.Run(testcase.name, func(t *testing.T) {
			got := getOutput(t, testcase.messageToPrint)
			assertPrinted(t, got, testcase.want)
		})
	}
}
