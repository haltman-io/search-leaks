package util

import "errors"

var ErrAPIConsecutiveFailures = errors.New("API returned errors 3 times in a row")
