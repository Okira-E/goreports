package utils

import "time"

// GetTimestamp returns the current timestamp in nanoseconds.
func GetTimestamp() int64 {
	return time.Now().UnixNano()
}
