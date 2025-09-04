package utils

import (
	"fmt"
	"time"
)

func Spinner[T any](f func() (T, error), message string) (T, error) {
	done := make(chan any)

	go func() {
		frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0
		for {
			select {
			case <-done:
				fmt.Println()
				return
			default:
				fmt.Printf("\r\033[K%s%s ... %s", Colorize(Gray, "[Info] "), message, frames[i%len(frames)])
				time.Sleep(100 * time.Millisecond)
				i++
			}
		}
	}()

	res, err := f()
	done <- res
	return res, err
}
