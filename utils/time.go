package utils

import "time"

// CurrentTime returns the current time in Unix format
func CurrentTime() int64 {
	return time.Now().Unix()
}
