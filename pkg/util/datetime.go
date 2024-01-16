package util

import (
	"strings"
	"time"

	"github.com/ffajarpratama/boiler-api/pkg/constant"
)

var TimeZone *time.Location

func SetTimeZone(location string) {
	TimeZone, _ = time.LoadLocation(location)
}

func TimeNow() time.Time {
	return time.Now().In(TimeZone)
}

func ParseDateToUnix(date string) int {
	res := 0
	gmt7 := 7 * 3600

	if strings.Contains(date, "-") {
		timeParse, err := time.Parse(constant.FormatYYYYMMDDHHMM, date)
		if err == nil {
			return int(timeParse.Unix()) - gmt7
		}
	}

	if !strings.Contains(date, ":") {
		timeParse, err := time.Parse(constant.FormatDate7, date)
		if err == nil {
			return int(timeParse.Unix()) - gmt7
		}
	}

	timeParse, err := time.Parse(constant.FormatDDMMYYYYHHMM, date)
	if err == nil {
		return int(timeParse.Unix()) - gmt7
	}

	timeParse, err = time.Parse(constant.FormatDate1, date)
	if err == nil {
		return int(timeParse.Unix()) - gmt7
	}

	return res
}
