package utils

import "time"

var ProcessStartTime = CurrentTime()

// CurrentTime returns the current time in Unix format
func CurrentTime() int64 {
	return time.Now().Unix()
}
