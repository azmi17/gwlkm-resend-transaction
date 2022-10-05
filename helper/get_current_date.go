package helper

import (
	"time"
)

const (
	DDMMYYYYhhmmss = "20060102150405"
)

func GetCurrentDate() string {
	now := time.Now()
	var getCurrentDate = string(now.Format(DDMMYYYYhhmmss))
	return getCurrentDate
}
