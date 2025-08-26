package utils

import (
	"fmt"
	"os"
)

func ExitOnError(message ...string) {
	fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}
