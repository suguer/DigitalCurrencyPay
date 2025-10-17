package util

import "time"

func Now() *time.Time {
	timer := time.Now()
	return &timer
}
func LastHour(hour int) string {
	timeLayout := "2006-01-02 15:04:05" //转化所需模板
	t := time.Now()
	t = t.Add(-time.Hour * time.Duration(hour))
	return t.Format(timeLayout)
}

func GetDateByTimestamp[T int64 | uint64](t T, layout ...string) time.Time {
	nt := int64(t)
	return time.UnixMilli(nt)
}

func GetDateByString(str string, layout ...string) time.Time {
	timeLayout := "2006-01-02 15:04:05" //转化所需模板
	if len(layout) > 0 {
		timeLayout = layout[0]
	}
	loc, err := time.LoadLocation("Local") //获取时区
	if err != nil {
		return time.Time{}
	}
	theTime, err := time.ParseInLocation(timeLayout, str, loc) //使用模板在对应时区转化为time.time类型
	if err != nil {
		return time.Time{}
	}
	return theTime
}
