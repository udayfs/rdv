package utils

import (
	"fmt"
	"os"
)

func ExitOnError(message string) {
	fmt.Fprintf(os.Stderr, "%s\n", Colorize(Red, "[Error] ")+message)
	os.Exit(1)
}
