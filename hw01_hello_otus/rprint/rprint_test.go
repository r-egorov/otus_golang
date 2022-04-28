package rprint_test

import (
	"bytes"
	"github.com/r-egorov/otus_golang/hw01_hello_otus/rprint"
	"testing"
)

func TestRevPrint(t *testing.T) {
	assertPrinted := func(t testing.TB, got, want string) {
		t.Helper()
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}

	assertNoError := func(t testing.TB, err error) {
		t.Helper()
		if err != nil {
			t.Errorf("got error %v", err)
		}
	}

	getOutput := func(t testing.TB, message string) string {
		t.Helper()

		out := &bytes.Buffer{}
		err := rprint.RevPrint(out, message)
		assertNoError(t, err)
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
