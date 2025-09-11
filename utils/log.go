package utils

import (
	"fmt"
	"os"
)

func ExitOnError(message string) {
	fmt.Fprintf(os.Stderr, "%s\n", Colorize(Red, "[Error] ")+message)
	os.Exit(1)
}

func ExitOnSuccess(message string) {
	fmt.Printf("%s\n", Colorize(Green, "[Success] ")+message)
	os.Exit(0)
}

func Info(message string) {
	fmt.Printf("%s\n", Colorize(Gray, "[Info] ")+message)
}

func Warn(message string) {
	fmt.Printf("%s\n", Colorize(Yellow, "[Warn] ")+message)
}
