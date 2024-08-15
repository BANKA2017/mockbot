package share

import "time"

// for GMT+8
var LocalTime, _ = time.LoadLocation("Asia/Shanghai")

var Now = time.Now().In(LocalTime)

func UpdateNow() {
	Now = time.Now().In(LocalTime)
}

func LocaleTimeDiff(hour int64) int64 {
	targetTime := time.Date(Now.Year(), Now.Month(), Now.Day(), int(hour), 0, 0, 0, LocalTime)

	if targetTime.After(Now) {
		targetTime = targetTime.Add(-24 * time.Hour)
	}

	return targetTime.Unix()
}
