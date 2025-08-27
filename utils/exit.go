package utils

import (
	"fmt"
	"os"
)

func ExitOnError(message string) {
	fmt.Fprintf(os.Stderr, "%s\n", message)
	os.Exit(1)
}
