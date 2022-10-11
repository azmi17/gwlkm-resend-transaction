package helper

import (
	"time"
)

// Format string date
const (
	DDMMYYYYhhmmss = "20060102150405"
)

func GetCurrentDate() string {
	now := time.Now()
	var getCurrentDate = string(now.Format(DDMMYYYYhhmmss))
	return getCurrentDate
}

// Custom Date
var (
	GETMMDD = GetCurrentDate()
)
