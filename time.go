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

// StringToDateTimePtr 将 string转为datetimePtr
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

// MakeDirNameByCurrentTime 根据当前时间生成诸如2024/01/01的目录名
func MakeDirNameByCurrentTime() string {
	return time.Now().Format(DateDirLayout)
}

// GetCurrentTimeAddMinutes 获取当前时间加某分钟后的时间
func GetCurrentTimeAddMinutes(minutes int) time.Time {
	return time.Now().Add(time.Duration(minutes) * time.Minute)
}

// GetCurrentTimeSubMinutes 获取当前时间减某分钟后的时间
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

// FormatRelativeDate 根据输入时间返回相对日期描述,otherTimeStr空则返回2006-01-02格式时间
func FormatRelativeDate(inputTime time.Time) string {
	now := GetCurrentTime()

	// 获取今天的开始时间（00:00:00）
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// 获取昨天的开始时间
	yesterday := today.AddDate(0, 0, -1)

	// 获取本月的开始时间
	thisMonthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// 获取上月的开始时间和结束时间
	lastMonthStart := thisMonthStart.AddDate(0, -1, 0)
	lastMonthEnd := thisMonthStart.AddDate(0, 0, -1)

	// 获取输入时间的日期部分（忽略具体时间）
	inputDate := time.Date(inputTime.Year(), inputTime.Month(), inputTime.Day(), 0, 0, 0, 0, inputTime.Location())

	// 判断是否是今天
	if inputDate.Equal(today) {
		return "今天"
	}

	// 判断是否是昨天
	if inputDate.Equal(yesterday) {
		return "昨天"
	}

	// 判断是否是本月（但不是今天和昨天）
	if inputDate.After(yesterday) && inputDate.Before(today.AddDate(0, 1, 0)) {
		return "本月"
	}

	// 判断是否是上月
	if (inputDate.Equal(lastMonthStart) || inputDate.After(lastMonthStart)) &&
		(inputDate.Equal(lastMonthEnd) || inputDate.Before(lastMonthEnd.AddDate(0, 0, 1))) {
		return "上月"
	}

	// 不符合任何条件，返回空字符串
	return ""
}

// GetTodayInterval 获取今天从开始到结束的时间区间
func GetTodayInterval() (start string, end string) {
	start = fmt.Sprintf("%s 00:00:00", GetCurrentDate())
	end = fmt.Sprintf("%s 00:00:00", GetCurrentTime().AddDate(0, 0, 1).Format("2006-01-02"))
	return
}

// GetYesterdayInterval 获取昨天从开始到结束的时间区间
func GetYesterdayInterval() (start string, end string) {
	start = fmt.Sprintf("%s 00:00:00", GetCurrentTime().AddDate(0, 0, -1).Format("2006-01-02"))
	end = fmt.Sprintf("%s 00:00:00", GetCurrentDate())
	return
}

// GetLastMonthInterval 获取上个月从开始到结束的时间区间
func GetLastMonthInterval() (start string, end string) {
	start = fmt.Sprintf("%s-01 00:00:00", GetCurrentTime().AddDate(0, -1, 0).Format("2006-01"))
	end = fmt.Sprintf("%s-01 00:00:00", GetCurrentTime().AddDate(0, 0, 0).Format("2006-01"))
	return
}

// GetCurrentMonthInterval 获取当前月从开始到结束的时间区间
func GetCurrentMonthInterval() (start string, end string) {
	start = fmt.Sprintf("%s-01 00:00:00", GetCurrentTime().AddDate(0, 0, 0).Format("2006-01"))
	end = fmt.Sprintf("%s-01 00:00:00", GetCurrentTime().AddDate(0, 1, 0).Format("2006-01"))
	return
}

// GetNextMonthInterval 获取下个月从开始到结束的时间区间
func GetNextMonthInterval() (start string, end string) {
	start = fmt.Sprintf("%s-01 00:00:00", GetCurrentTime().AddDate(0, 1, 0).Format("2006-01"))
	end = fmt.Sprintf("%s-01 00:00:00", GetCurrentTime().AddDate(0, 2, 0).Format("2006-01"))
	return
}

// GetLastYearsInterval 获取去年从开始到结束的时间区间
func GetLastYearsInterval() (start string, end string) {
	start = fmt.Sprintf("%s-01-01 00:00:00", GetCurrentTime().AddDate(-1, 0, 0).Format("2006"))
	end = fmt.Sprintf("%s-01-01 00:00:00", GetCurrentTime().AddDate(0, 0, 0).Format("2006"))
	return
}

// GetCurrentYearsInterval 获取今年从开始到结束的时间区间
func GetCurrentYearsInterval() (start string, end string) {
	start = fmt.Sprintf("%s-01-01 00:00:00", GetCurrentTime().AddDate(0, 0, 0).Format("2006"))
	end = fmt.Sprintf("%s-01-01 00:00:00", GetCurrentTime().AddDate(1, 0, 0).Format("2006"))
	return
}

// GetNextYearsInterval 获取明年从开始到结束的时间区间
func GetNextYearsInterval() (start string, end string) {
	start = fmt.Sprintf("%s-01-01 00:00:00", GetCurrentTime().AddDate(1, 0, 0).Format("2006"))
	end = fmt.Sprintf("%s-01-01 00:00:00", GetCurrentTime().AddDate(2, 0, 0).Format("2006"))
	return
}
