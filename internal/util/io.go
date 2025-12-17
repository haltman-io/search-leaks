package util

import (
	"os"
)

func HasStdinData() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	// If stdin is not a character device, it likely has piped data.
	return (fi.Mode() & os.ModeCharDevice) == 0
}
