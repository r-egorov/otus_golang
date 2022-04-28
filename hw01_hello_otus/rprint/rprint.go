package rprint

import (
	"fmt"
	"io"

	"golang.org/x/example/stringutil"
)

func RevPrint(out io.Writer, message string) error {
	_, err := fmt.Fprint(out, stringutil.Reverse(message))
	if err != nil {
		return err
	}
	return nil
}
