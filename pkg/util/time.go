package util

import "time"

const (
	TimeFormat     = "2006-01-02 15:04:05"
	TimeFormatDay  = "20060102"
	TimeFormatDay2 = "2006-01-02"
	TimeFormatDay3 = "2006/01/02"
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
