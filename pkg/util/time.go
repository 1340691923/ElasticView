package util

import (
	"log"
	"time"
)

const (
	TimeFormat  = "2006-01-02 15:04:05"
	TimeFormat2 = "2006-01-02 15:04:00"

	TimeFormatDay    = "20060102"
	TimeFormatMonth  = "200601"
	TimeFormatMonth2 = "2006-01"
	TimeFormatDay2   = "2006-01-02"
	TimeFormatHour   = "02号15:04"
	TimeFormatDay3   = "2006/01/02"
	TimeFormatDay4   = "2006.01.02_15"
	TimeFormatDay5   = "2006-01-02 00:00:00"
	TimeFormatDay6   = "2006年01月02日"
)

/**
 * 二个时间戳是否同一天
 * @return true 是 false 不是今天
 */
func IsSameDay(oldDay, anotherDay int64) bool {
	tm := time.Unix(oldDay, 0)
	tmAnother := time.Unix(anotherDay, 0)
	if tm.Format(TimeFormatDay2) == tmAnother.Format(TimeFormatDay2) {
		return true
	}
	return false
}

/**字符串->时间对象*/
func Str2Time(formatTimeStr, timeFormat string) time.Time {
	loc, err := time.LoadLocation("Local")
	theTime, err := time.ParseInLocation(timeFormat, formatTimeStr, loc) //使用模板在对应时区转化为time.time类型
	log.Println(err, formatTimeStr, timeFormat)
	return theTime
}

/**字符串->时间对象*/
func StrToTime(timeFormat, formatTimeStr string) time.Time {

	theTime, err := time.ParseInLocation(timeFormat, formatTimeStr, time.Local) //使用模板在对应时区转化为time.time类型
	if err != nil {
		log.Println(err)
	}
	return theTime
}

/**字符串->时间对象*/
func StrToTimeHaveErr(timeFormat, formatTimeStr string) (time.Time, error) {

	theTime, err := time.ParseInLocation(timeFormat, formatTimeStr, time.Local) //使用模板在对应时区转化为time.time类型
	if err != nil {
		return time.Now(), err
	}
	return theTime, nil
}

type TimeLog struct {
	T time.Time
}

func NewTimeLog() *TimeLog {
	return &TimeLog{T: time.Now()}
}

func (this *TimeLog) Log(tag string) {
	return
	log.Println(tag, time.Now().Sub(this.T).String())
}
