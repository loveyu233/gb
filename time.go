package gb

import (
	"fmt"
	"time"

	"github.com/araddon/dateparse"
)

var (
	ShangHaiTimeLocation *time.Location
)

// CSTLayout China Standard Time Layout
const (
	CSTLayout                       = "2006-01-02 15:04:05"
	CSTLayoutChinese                = "2006年01月02日 15:04:05"
	CSTLayoutPoint                  = "2006.01.02 15:04:05"
	CSTLayoutDate                   = "2006-01-02"
	CSTLayoutDateChinese            = "2006年01月02日"
	CSTLayoutDatePoint              = "2006.01.02"
	CSTLayoutTime                   = "15:04:05"
	CSTLayoutDateHourMinutes        = "2006-01-02 15:04"
	CSTLayoutDateHourMinutesChinese = "2006年01月02日 15:04"
	CSTLayoutDateHourMinutesPoint   = "2006.01.02 15:04"
	CSTLayoutYearMonth              = "2006-01"
	CSTLayoutYearMonthChinese       = "2006年01月"
	CSTLayoutYearMonthPoint         = "2006.01.02"
	CSTLayoutSecond                 = "20060102150405"
	DateDirLayout                   = "2006/0101"

	DayStartTimeStr = "00:00:00"
	DayEndTimeStr   = "23:59:59"
)

func init() {
	var err error
	if ShangHaiTimeLocation, err = time.LoadLocation("Asia/Shanghai"); err != nil {
		panic(err)
	}

	// 默认设置为中国时区
	time.Local = ShangHaiTimeLocation
}

// Now 获取当前时间
func Now() time.Time {
	return time.Now().In(ShangHaiTimeLocation)
}

func NowGBTimeOnly() TimeOnly {
	return TimeOnly(time.Date(0, 0, 0, Now().Hour(), Now().Minute(), Now().Second(), 0, ShangHaiTimeLocation))
}

func NowGBDateOnly() DateOnly {
	return DateOnly(time.Date(Now().Year(), Now().Month(), Now().Day(), 0, 0, 0, 0, ShangHaiTimeLocation))
}

func NowGBDateTime() DateTime {
	return DateTime(time.Date(Now().Year(), Now().Month(), Now().Day(), Now().Hour(), Now().Minute(), Now().Second(), 0, ShangHaiTimeLocation))
}

func NowGBTimeHourMinute() TimeHourMinute {
	return TimeHourMinute(time.Date(0, 0, 0, Now().Hour(), Now().Minute(), 0, 0, ShangHaiTimeLocation))
}

// NowPtr 获取当前时间指针
func NowPtr() *time.Time {
	now := Now()
	return &now
}

// NowString 格式化当前时间
func NowString() string {
	return Now().Format(CSTLayout)
}

// NowDateString 获取当前日期字符串
func NowDateString() string {
	return Now().Format(CSTLayoutDate)
}

// NowTimeString 获取当前时间字符串
func NowTimeString() string {
	return Now().Format(CSTLayoutTime)
}

// TimeToGBDateTime 把时间类型转为gb库的DateTime类型
func TimeToGBDateTime(t time.Time) DateTime {
	return DateTime(t)
}

// TimeToGBDateOnly 把时间类型转为gb库的DateOnly类型
func TimeToGBDateOnly(t time.Time) DateOnly {
	return DateOnly(t)
}

// TimeToGBTimeOnly 把时间类型转为gb库的TimeOnly类型
func TimeToGBTimeOnly(t time.Time) TimeOnly {
	return TimeOnly(t)
}

// TimeToGBTimeOnlyNoSec 把时间类型转为gb库的TimeOnly类型秒设置为00
func TimeToGBTimeOnlyNoSec(t time.Time) TimeOnly {
	return TimeOnly(time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, ShangHaiTimeLocation))
}

// StringToDateTime 时间字符串转为time.time
func StringToDateTime(dateTime string) (time.Time, error) {
	return time.ParseInLocation(CSTLayout, dateTime, ShangHaiTimeLocation)
}

// StringToGbDateTime 将日期字符串转为gb库的DateTime类型
func StringToGbDateTime(dateTime string) DateTime {
	location, _ := time.ParseInLocation(CSTLayout, dateTime, ShangHaiTimeLocation)
	return DateTime(location)
}

// StringToGbDateTimeErr 将日期字符串转为gb库的DateTime类型返回错误的
func StringToGbDateTimeErr(dateTime string) (DateTime, error) {
	location, err := time.ParseInLocation(CSTLayout, dateTime, ShangHaiTimeLocation)
	if err != nil {
		return DateTime{}, err
	}
	return DateTime(location), nil
}

// StringToDate 日期字符串转为time Date
func StringToDate(date string) (time.Time, error) {
	return time.ParseInLocation(CSTLayoutDate, date, ShangHaiTimeLocation)
}

// StringToDateNoErr 将日期字符串转为time无错误返回的
func StringToDateNoErr(dateTime string) time.Time {
	t, _ := StringToDate(dateTime)
	return t.In(ShangHaiTimeLocation)
}

// StringToGBDateOnly 将时间转为gb库的DateOnly
func StringToGBDateOnly(dateTime string) (DateOnly, error) {
	t, err := StringToDate(dateTime)
	if err != nil {
		return DateOnly{}, err
	}
	return DateOnly(t), nil
}

// StringToTime 日期字符串转为time Date
func StringToTime(date string) (time.Time, error) {
	return time.ParseInLocation(CSTLayoutTime, date, ShangHaiTimeLocation)
}

func StringToTimeNoErr(date string) time.Time {
	toTime, _ := StringToTime(date)
	return toTime
}

func StringToGBTimeOnly(time string) (TimeOnly, error) {
	toTime, err := StringToTime(time)
	if err != nil {
		return TimeOnly{}, err
	}
	return TimeOnly(toTime), nil
}

func StringToGBTimeHourMinute(t string) (TimeHourMinute, error) {
	toTime, err := StringToTime(t)
	if err != nil {
		return TimeHourMinute{}, err
	}

	return TimeHourMinute(time.Date(toTime.Year(), toTime.Month(), toTime.Day(), toTime.Hour(), toTime.Minute(), 0, 0, ShangHaiTimeLocation)), nil
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
	return t.In(ShangHaiTimeLocation).Format(CSTLayout)
}

// FuzzParseTimeString 模糊解析时间
func FuzzParseTimeString(timeString string) (time.Time, error) {
	return dateparse.ParseIn(timeString, ShangHaiTimeLocation)
}

// DateTimePtrToString *time.time类型转换为CSTLayout格式
func DateTimePtrToString(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.In(ShangHaiTimeLocation).Format(CSTLayout)
}

// DatetimePtrToDateString *time.time类型转换为CSTLayoutDate格式
func DatetimePtrToDateString(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.In(ShangHaiTimeLocation).Format(CSTLayoutDate)
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

	t, err := time.ParseInLocation(time.RFC3339, rfc3339, ShangHaiTimeLocation)
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
	return time.Unix(unix, 0).In(ShangHaiTimeLocation).Format(CSTLayout)
}

// GetCurrentTimeUnix 获取当前时间的 Unix 时间戳
func GetCurrentTimeUnix() int {
	return TimeToUnix(Now())
}

// MakeDirNameByCurrentTime 根据当前时间生成诸如2024/01/01的目录名
func MakeDirNameByCurrentTime() string {
	return Now().Format(DateDirLayout)
}

// GetCurrentTimeAddMinutes 获取当前时间加某分钟后的时间
func GetCurrentTimeAddMinutes(minutes int) time.Time {
	return Now().Add(time.Duration(minutes) * time.Minute)
}

// GetCurrentTimeSubMinutes 获取当前时间减某分钟后的时间
func GetCurrentTimeSubMinutes(minutes int) time.Time {
	return Now().Add(-time.Duration(minutes) * time.Minute)
}

// AfterMinutes 判断时间是否在某个时间之后
func AfterMinutes(t time.Time, minutes int) bool {
	return t.After(Now().Add(-time.Duration(minutes) * time.Minute))
}

// GetCurrentTimeSubHours 获取当前时间减去指定小时数后的时间
func GetCurrentTimeSubHours(hours int) time.Time {
	return Now().Add(-time.Duration(hours) * time.Hour)
}

// FormatDateRelativeDate 根据输入时间返回相对日期描述,otherTimeStr空则返回2006-01-02格式时间
func FormatDateRelativeDate(inputTime time.Time) string {
	now := Now()

	year, month, day := inputTime.Date()
	if year == now.Year() && month == now.Month() && day == now.Day() {
		return "今天"
	}
	if year == now.Year() && month == now.Month() && day == now.AddDate(0, 0, 1).Day() {
		return "明天"
	}
	if year == now.Year() && month == now.Month() && day == now.AddDate(0, 0, -1).Day() {
		return "昨天"
	}
	if year == now.Year() && month == now.Month() {
		return "本月"
	}
	if year == now.Year() && month == now.AddDate(0, -1, 0).Month() {
		return "上月"
	}

	// 不符合任何条件，返回空字符串
	return ""
}

func FormatTimeRelativeDate(t time.Time) string {
	hour := t.Hour()
	switch {
	case hour >= 0 && hour < 6:
		return "凌晨"
	case hour >= 6 && hour < 12:
		return "上午"
	case hour == 12:
		return "中午"
	case hour >= 13 && hour < 18:
		return "下午"
	case hour >= 18 && hour < 24:
		return "晚上"
	default:
		return ""
	}
}

// GetTodayInterval 获取今天从开始到结束的时间区间
func GetTodayInterval() (start DateTime, end DateTime) {
	start = StringToGbDateTime(fmt.Sprintf("%s 00:00:00", NowDateString()))
	end = StringToGbDateTime(fmt.Sprintf("%s 00:00:00", Now().AddDate(0, 0, 1).Format("2006-01-02")))
	return
}

// GetYesterdayInterval 获取昨天从开始到结束的时间区间
func GetYesterdayInterval() (start DateTime, end DateTime) {
	start = StringToGbDateTime(fmt.Sprintf("%s 00:00:00", Now().AddDate(0, 0, -1).Format("2006-01-02")))
	end = StringToGbDateTime(fmt.Sprintf("%s 00:00:00", NowDateString()))
	return
}

// GetLastMonthInterval 获取上个月从开始到结束的时间区间
func GetLastMonthInterval() (start DateTime, end DateTime) {
	start = StringToGbDateTime(fmt.Sprintf("%s-01 00:00:00", Now().AddDate(0, -1, 0).Format("2006-01")))
	end = StringToGbDateTime(fmt.Sprintf("%s-01 00:00:00", Now().AddDate(0, 0, 0).Format("2006-01")))
	return
}

// GetCurrentMonthInterval 获取当前月从开始到结束的时间区间
func GetCurrentMonthInterval() (start DateTime, end DateTime) {
	start = StringToGbDateTime(fmt.Sprintf("%s-01 00:00:00", Now().AddDate(0, 0, 0).Format("2006-01")))
	end = StringToGbDateTime(fmt.Sprintf("%s-01 00:00:00", Now().AddDate(0, 1, 0).Format("2006-01")))
	return
}

// GetNextMonthInterval 获取下个月从开始到结束的时间区间
func GetNextMonthInterval() (start DateTime, end DateTime) {
	start = StringToGbDateTime(fmt.Sprintf("%s-01 00:00:00", Now().AddDate(0, 1, 0).Format("2006-01")))
	end = StringToGbDateTime(fmt.Sprintf("%s-01 00:00:00", Now().AddDate(0, 2, 0).Format("2006-01")))
	return
}

// GetLastYearsInterval 获取去年从开始到结束的时间区间
func GetLastYearsInterval() (start DateTime, end DateTime) {
	start = StringToGbDateTime(fmt.Sprintf("%s-01-01 00:00:00", Now().AddDate(-1, 0, 0).Format("2006")))
	end = StringToGbDateTime(fmt.Sprintf("%s-01-01 00:00:00", Now().AddDate(0, 0, 0).Format("2006")))
	return
}

// GetCurrentYearsInterval 获取今年从开始到结束的时间区间
func GetCurrentYearsInterval() (start DateTime, end DateTime) {
	start = StringToGbDateTime(fmt.Sprintf("%s-01-01 00:00:00", Now().AddDate(0, 0, 0).Format("2006")))
	end = StringToGbDateTime(fmt.Sprintf("%s-01-01 00:00:00", Now().AddDate(1, 0, 0).Format("2006")))
	return
}

// GetNextYearsInterval 获取明年从开始到结束的时间区间
func GetNextYearsInterval() (start DateTime, end DateTime) {
	start = StringToGbDateTime(fmt.Sprintf("%s-01-01 00:00:00", Now().AddDate(1, 0, 0).Format("2006")))
	end = StringToGbDateTime(fmt.Sprintf("%s-01-01 00:00:00", Now().AddDate(2, 0, 0).Format("2006")))
	return
}

// TimeChineseWeekday 获取日期的星期几（中文）
func TimeChineseWeekday(t time.Time) string {
	weekdays := []string{"日", "一", "二", "三", "四", "五", "六"}
	return "星期" + weekdays[t.Weekday()]
}

// TimeEnglishWeekday 获取日期的星期几（英文）
func TimeEnglishWeekday(t time.Time) string {
	return t.Weekday().String()
}

// TimeIsWeekend 判断是否是周末
func TimeIsWeekend(t time.Time) bool {
	return t.Weekday() == time.Saturday || t.Weekday() == time.Sunday
}

// TimeRange 表示一个时间段
type TimeRange struct {
	ID    uint64
	Start time.Time
	End   time.Time
}

// IsValid 检查时间段是否有效
func (tr TimeRange) IsValid() bool {
	return !tr.Start.After(tr.End)
}

// HasConflictWith 检查当前时间段是否与另一个时间段冲突
func (tr TimeRange) HasConflictWith(other TimeRange) bool {
	if !tr.IsValid() || !other.IsValid() {
		return false
	}
	return tr.Start.Before(other.End) && tr.End.After(other.Start)
}

// HasTimeConflict 检查多个时间段之间是否有冲突
// 参数: 可变数量的TimeRange，每个TimeRange包含开始时间和结束时间
// 返回 true 表示存在冲突，false 表示无冲突
func HasTimeConflict(timeRanges ...TimeRange) bool {
	// 如果时间段数量少于2个，不可能有冲突
	if len(timeRanges) < 2 {
		return false
	}

	// 检查每一对时间段是否冲突
	for i := 0; i < len(timeRanges); i++ {
		for j := i + 1; j < len(timeRanges); j++ {
			if timeRanges[i].HasConflictWith(timeRanges[j]) {
				return true
			}
		}
	}

	return false
}

func HasTimeConflictReturnIDS(ranges ...TimeRange) []uint64 {
	overlappingIDs := make(map[uint64]bool)

	for i := 0; i < len(ranges); i++ {
		for j := i + 1; j < len(ranges); j++ {
			// 检查时间重合条件
			if ranges[i].Start.Before(ranges[j].End) && ranges[i].End.After(ranges[j].Start) {
				overlappingIDs[ranges[i].ID] = true
				overlappingIDs[ranges[j].ID] = true
			}
		}
	}

	// 转换为切片
	var result []uint64
	for id := range overlappingIDs {
		result = append(result, id)
	}

	return result
}
