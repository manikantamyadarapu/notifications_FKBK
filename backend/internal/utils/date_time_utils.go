package utils

import "time"

const istLocationName = "Asia/Kolkata"
const istDateTimeLayout = "2006-01-02 15:04:05 MST"

func ToIST(t time.Time) time.Time {
	loc, err := time.LoadLocation(istLocationName)
	if err != nil {
		// Fallback to fixed IST offset in case timezone data is unavailable.
		loc = time.FixedZone("IST", 5*60*60+30*60)
	}
	return t.In(loc)
}

func FormatIST(t time.Time) string {
	return ToIST(t).Format(istDateTimeLayout)
}

