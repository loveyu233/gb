package gb

import (
	"fmt"
	"github.com/araddon/dateparse"
	"time"
)

var (
	cst *time.Location
)

// CSTLayout China Standard Time Layout
const (
	CSTLayout                = "2006-01-02 15:04:05"
	CSTLayoutDate            = "2006-01-02"
	CSTLayoutDateHourMinutes = "2006-01-02 15:04"
	CSTLayoutYearMonth       = "2006-01"
	CSTLayoutSecond          = "20060102150405"
	DateDirLayout            = "2006/0101"

	DayStartTimeStr = "00:00:00"
	DayEndTimeStr   = "23:59:59"
)

func init() {
	var err error
	if cst, err = time.LoadLocation("Asia/Shanghai"); err != nil {
		panic(err)
	}

	// 默认设置为中国时区
	time.Local = cst
}

// GetCurrentTime 获取当前时间
func GetCurrentTime() time.Time {
	return time.Now()
}

// GetCurrentTimePtr 获取当前时间指针
func GetCurrentTimePtr() *time.Time {
	now := time.Now()
	return &now
}

// FormatCurrentTime 格式化当前时间
func FormatCurrentTime() string {
	ts := time.Now()
	return ts.In(cst).Format(CSTLayout)
}

// GetCurrentDate 获取当前日期
func GetCurrentDate() string {
	return time.Now().In(cst).Format(CSTLayoutDate)
}

// StringToDateTime 时间字符串转为time.time
func StringToDateTime(dateTime string) (time.Time, error) {
	return time.ParseInLocation(CSTLayout, dateTime, cst)
}

// StringToDate 日期字符串转为time Date
func StringToDate(date string) (time.Time, error) {
	return time.ParseInLocation(CSTLayoutDate, date, cst)
}

// StringDateToDateTimePtr 日期字符串转为time.time
func StringDateToDateTimePtr(date string, hourMinuteSecond string) (parsed *time.Time, err error) {
	if date == "" {
		return
	}

	if hourMinuteSecond != "" {
		date = fmt.Sprintf("%s %s", date, hourMinuteSecond)
	}

	_parsed, err := StringToDateTime(date)
	if err != nil {
		return nil, err
	}

	return &_parsed, nil
}

// DateTimeToString time.time类型转换为String类型
func DateTimeToString(t time.Time) string {
	return t.In(cst).Format(CSTLayout)
}

// FuzzParseTimeString 模糊解析时间
func FuzzParseTimeString(timeString string) (time.Time, error) {
	return dateparse.ParseIn(timeString, cst)
}

// DateTimePtrToString *time.time类型转换为CSTLayout格式
func DateTimePtrToString(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.In(cst).Format(CSTLayout)
}

// DatetimePtrToDateString *time.time类型转换为CSTLayoutDate格式
func DatetimePtrToDateString(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.In(cst).Format(CSTLayoutDate)
}

// 将 string转为datetimePtr
func StringToDateTimePtr(dateTime string) *time.Time {
	if dateTime == "" {
		return nil
	}

	t, err := StringToDateTime(dateTime)
	if err != nil {
		return nil
	}

	return &t
}

// Rfc3339StringToDateTimePtr 将rfc3339时间字符串转为time.time Ptr
func Rfc3339StringToDateTimePtr(rfc3339 string) *time.Time {
	if rfc3339 == "" {
		return nil
	}

	t, err := time.ParseInLocation(time.RFC3339, rfc3339, cst)
	if err != nil {
		return nil
	}

	return &t
}

// DurationFromSeconds 将秒数转换为 time.Duration
func DurationFromSeconds(seconds int) time.Duration {
	return time.Duration(seconds) * time.Second
}

// TimeToUnix 将时间转换为 Unix 时间戳
func TimeToUnix(t time.Time) int {
	return int(t.Unix())
}

// UnixToDateTimeString 将 Unix 时间戳转换为 CSTLayout 格式的时间字符串
func UnixToDateTimeString(unix int64) string {
	return time.Unix(unix, 0).In(cst).Format(CSTLayout)
}

// GetCurrentTimeUnix 获取当前时间的 Unix 时间戳
func GetCurrentTimeUnix() int {
	return TimeToUnix(time.Now())
}

// 根据当前时间生成诸如2024/01/01的目录名
func MakeDirNameByCurrentTime() string {
	return time.Now().Format(DateDirLayout)
}

// 获取当前时间加某分钟后的时间
func GetCurrentTimeAddMinutes(minutes int) time.Time {
	return time.Now().Add(time.Duration(minutes) * time.Minute)
}

// 获取当前时间减某分钟后的时间
func GetCurrentTimeSubMinutes(minutes int) time.Time {
	return time.Now().Add(-time.Duration(minutes) * time.Minute)
}

// AfterMinutes 判断时间是否在某个时间之后
func AfterMinutes(t time.Time, minutes int) bool {
	return t.After(time.Now().Add(-time.Duration(minutes) * time.Minute))
}

// GetCurrentTimeSubHours 获取当前时间减去指定小时数后的时间
func GetCurrentTimeSubHours(hours int) time.Time {
	return time.Now().Add(-time.Duration(hours) * time.Hour)
}
