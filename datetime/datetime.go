package datetime

import "time"

//ParesLoc 解析时间字符串到 本地时区的 time
func ParesLoc(layout, val string) (date time.Time) {

	loc, _ := time.LoadLocation("Local")

	date, _ = time.ParseInLocation(layout, val, loc)

	return
}
