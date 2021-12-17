package util

import "time"

// GetDayUnix of UTC+8 00:00:00 if the day utc time belong to.
func GetDayUnix(utc time.Time) int64 {
	unix := utc.Unix()
	unix -= unix % int64((time.Hour * 24).Seconds())
	return unix - int64((time.Hour * 8).Seconds())
}
