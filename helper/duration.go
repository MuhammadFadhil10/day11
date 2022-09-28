package helper

import (
	"strconv"
	"time"
)

func GetDuration(startDate string, endDate string) string {
	// fmt.Println(startDate)
	// fmt.Println(endDate)
	var duration string
	layout := "2006-01-02"
	parsedStartDate, _ := time.Parse(layout,startDate)
	parsedEndDate, _ := time.Parse(layout,endDate)

	var startMs = parsedStartDate.UnixMicro()
	var endMs = parsedEndDate.UnixMicro()

	margin := ((endMs - startMs) / (1000 * 60 * 60 * 24) / 1000)

	// fmt.Println(margin)

	if margin < 30 {
		if margin == 0 {
			duration = "a few hours";
		} else {
			duration = strconv.Itoa(int(margin)) + " Day"
		}
	}  else {
		if margin < 365 {
			if margin % 30 == 0 {
				duration = strconv.Itoa(int(margin / 30)) + " Month"
			} else {
				duration = strconv.Itoa(int(margin / 30)) + " Month " + strconv.Itoa(int(margin % 30)) + " Day"
			}
		} else {
			duration = strconv.Itoa(int(margin / 365)) + " Year"
		}
	}

	return duration
}