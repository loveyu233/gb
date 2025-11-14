package gb

import (
	"fmt"
	"time"

	"github.com/araddon/dateparse"
)

var (
	ShangHaiTimeLocation *time.Location
)

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

// init 初始化默认的上海时区配置。
func init() {
	var err error
	if ShangHaiTimeLocation, err = time.LoadLocation("Asia/Shanghai"); err != nil {
		panic(err)
	}
	time.Local = ShangHaiTimeLocation
}

// Now 返回当前上海时区的时间。
func Now() time.Time {
	return time.Now().In(ShangHaiTimeLocation)
}

// NowPointer 返回当前时间的指针副本，便于与可选参数兼容。
func NowPointer() *time.Time {
	now := Now()
	return &now
}

// NowDateTimeString 以标准格式返回当前日期时间字符串。
func NowDateTimeString() string {
	return FormatDateTime(Now())
}

// NowDateString 返回当前日期字符串（YYYY-MM-DD）。
func NowDateString() string {
	return Now().Format(CSTLayoutDate)
}

// NowTimeString 返回当前时间字符串（HH:MM:SS）。
func NowTimeString() string {
	return Now().Format(CSTLayoutTime)
}

// NowAsDateTime 返回当前时间的 DateTime 封装类型。
func NowAsDateTime() DateTime {
	return DateTime(Now())
}

// NowAsDateOnly 返回当天零点的 DateOnly 封装类型。
func NowAsDateOnly() DateOnly {
	now := Now()
	return DateOnly(time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, ShangHaiTimeLocation))
}

// NowAsTimeOnly 返回当前时间的 TimeOnly 封装类型。
func NowAsTimeOnly() TimeOnly {
	now := Now()
	return TimeOnly(time.Date(0, 0, 0, now.Hour(), now.Minute(), now.Second(), 0, ShangHaiTimeLocation))
}

// NowAsHourMinute 返回当前时间的时分部分。
func NowAsHourMinute() TimeHourMinute {
	now := Now()
	return TimeHourMinute(time.Date(0, 0, 0, now.Hour(), now.Minute(), 0, 0, ShangHaiTimeLocation))
}

// ToDateTime 将标准 time.Time 转为 DateTime 类型。
func ToDateTime(t time.Time) DateTime {
	return DateTime(t)
}

// ToDateOnly 将 time.Time 转为 DateOnly，只保留日期部分。
func ToDateOnly(t time.Time) DateOnly {
	return DateOnly(t)
}

// ToTimeOnly 将 time.Time 转为 TimeOnly，只保留时间部分。
func ToTimeOnly(t time.Time) TimeOnly {
	return TimeOnly(t)
}

// ToTimeOnlyTrimSeconds 将 time.Time 转为无秒的 TimeOnly。
func ToTimeOnlyTrimSeconds(t time.Time) TimeOnly {
	return TimeOnly(time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, ShangHaiTimeLocation))
}

// ParseDateTime 解析标准日期时间字符串为 time.Time。
func ParseDateTime(value string) (time.Time, error) {
	return time.ParseInLocation(CSTLayout, value, ShangHaiTimeLocation)
}

// ParseDateTimeValue 解析字符串并返回 DateTime 类型。
func ParseDateTimeValue(value string) (DateTime, error) {
	parsed, err := ParseDateTime(value)
	if err != nil {
		return DateTime{}, err
	}
	return DateTime(parsed), nil
}

// MustParseDateTimeValue 解析字符串为 DateTime，失败返回零值。
func MustParseDateTimeValue(value string) DateTime {
	parsed, err := ParseDateTime(value)
	if err != nil {
		return DateTime{}
	}
	return DateTime(parsed)
}

// ParseDate 解析日期字符串为 time.Time。
func ParseDate(value string) (time.Time, error) {
	return time.ParseInLocation(CSTLayoutDate, value, ShangHaiTimeLocation)
}

// MustParseDate 解析日期字符串，失败返回零值。
func MustParseDate(value string) time.Time {
	parsed, err := ParseDate(value)
	if err != nil {
		return time.Time{}
	}
	return parsed
}

// ParseDateOnly 解析日期字符串为 DateOnly。
func ParseDateOnly(value string) (DateOnly, error) {
	parsed, err := ParseDate(value)
	if err != nil {
		return DateOnly{}, err
	}
	return DateOnly(parsed), nil
}

// ParseClock 解析时间字符串为 time.Time。
func ParseClock(value string) (time.Time, error) {
	return time.ParseInLocation(CSTLayoutTime, value, ShangHaiTimeLocation)
}

// MustParseClock 解析时间字符串，失败返回零值。
func MustParseClock(value string) time.Time {
	parsed, err := ParseClock(value)
	if err != nil {
		return time.Time{}
	}
	return parsed
}

// ParseTimeOnly 解析时间字符串为 TimeOnly。
func ParseTimeOnly(value string) (TimeOnly, error) {
	parsed, err := ParseClock(value)
	if err != nil {
		return TimeOnly{}, err
	}
	return TimeOnly(parsed), nil
}

// ParseHourMinute 解析时间字符串并返回时分结构。
func ParseHourMinute(value string) (TimeHourMinute, error) {
	parsed, err := ParseClock(value)
	if err != nil {
		return TimeHourMinute{}, err
	}
	return TimeHourMinute(time.Date(parsed.Year(), parsed.Month(), parsed.Day(), parsed.Hour(), parsed.Minute(), 0, 0, ShangHaiTimeLocation)), nil
}

// ParseDateAndTimePointer 将日期与时间字符串组合为可选的 time.Time 指针。
func ParseDateAndTimePointer(date string, hourMinuteSecond string) (*time.Time, error) {
	if date == "" {
		return nil, nil
	}
	if hourMinuteSecond != "" {
		date = fmt.Sprintf("%s %s", date, hourMinuteSecond)
	}
	parsed, err := ParseDateTime(date)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

// FormatDateTime 以标准格式输出 time.Time。
func FormatDateTime(t time.Time) string {
	return t.In(ShangHaiTimeLocation).Format(CSTLayout)
}

// ParseFuzzyTime 使用 dateparse 模糊解析时间。
func ParseFuzzyTime(value string) (time.Time, error) {
	return dateparse.ParseIn(value, ShangHaiTimeLocation)
}

// FormatDateTimePointer 安全地格式化 time.Time 指针。
func FormatDateTimePointer(t *time.Time) string {
	if t == nil {
		return ""
	}
	return FormatDateTime(*t)
}

// FormatDatePointer 将 time.Time 指针格式化为日期字符串。
func FormatDatePointer(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.In(ShangHaiTimeLocation).Format(CSTLayoutDate)
}

// ParseDateTimePointer 解析日期时间字符串并返回指针。
func ParseDateTimePointer(value string) *time.Time {
	if value == "" {
		return nil
	}
	parsed, err := ParseDateTime(value)
	if err != nil {
		return nil
	}
	return &parsed
}

// ParseRFC3339Pointer 解析 RFC3339 时间并返回指针。
func ParseRFC3339Pointer(value string) *time.Time {
	if value == "" {
		return nil
	}
	parsed, err := time.ParseInLocation(time.RFC3339, value, ShangHaiTimeLocation)
	if err != nil {
		return nil
	}
	return &parsed
}

// SecondsToDuration 将秒数转换为 duration。
func SecondsToDuration(seconds int) time.Duration {
	return time.Duration(seconds) * time.Second
}

// UnixFromTime 返回时间的 Unix 秒。
func UnixFromTime(t time.Time) int {
	return int(t.Unix())
}

// FormatUnixDateTime 将 Unix 秒格式化为日期字符串。
func FormatUnixDateTime(unix int64) string {
	return time.Unix(unix, 0).In(ShangHaiTimeLocation).Format(CSTLayout)
}

// NowUnix 返回当前时间的 Unix 秒。
func NowUnix() int {
	return UnixFromTime(Now())
}

// NowDateDirectory 生成按日期分层的目录名。
func NowDateDirectory() string {
	return Now().Format(DateDirLayout)
}

// NowAddMinutes 返回当前时间增加指定分钟后的时间。
func NowAddMinutes(minutes int) time.Time {
	return Now().Add(time.Duration(minutes) * time.Minute)
}

// NowSubMinutes 返回当前时间减少指定分钟后的时间。
func NowSubMinutes(minutes int) time.Time {
	return Now().Add(-time.Duration(minutes) * time.Minute)
}

// IsAfterMinutesAgo 判断时间是否晚于若干分钟前。
func IsAfterMinutesAgo(t time.Time, minutes int) bool {
	return t.After(Now().Add(-time.Duration(minutes) * time.Minute))
}

// NowSubHours 返回当前时间减少指定小时后的时间。
func NowSubHours(hours int) time.Time {
	return Now().Add(-time.Duration(hours) * time.Hour)
}

// DescribeRelativeDate 返回日期相对于今天的中文描述。
func DescribeRelativeDate(input time.Time) string {
	now := Now()
	y, m, d := input.Date()
	switch {
	case y == now.Year() && m == now.Month() && d == now.Day():
		return "今天"
	case y == now.Year() && m == now.Month() && d == now.AddDate(0, 0, 1).Day():
		return "明天"
	case y == now.Year() && m == now.Month() && d == now.AddDate(0, 0, -1).Day():
		return "昨天"
	case y == now.Year() && m == now.Month():
		return "本月"
	case y == now.Year() && m == now.AddDate(0, -1, 0).Month():
		return "上月"
	default:
		return ""
	}
}

// DescribeRelativeTimeOfDay 根据小时返回一天中的描述。
func DescribeRelativeTimeOfDay(t time.Time) string {
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

// TodayRange 返回今天的起止时间范围。
func TodayRange() (DateTime, DateTime) {
	start := beginningOfDay(Now())
	end := start.AddDate(0, 0, 1)
	return buildRange(start, end)
}

// YesterdayRange 返回昨天的起止时间范围。
func YesterdayRange() (DateTime, DateTime) {
	start := beginningOfDay(Now().AddDate(0, 0, -1))
	end := start.AddDate(0, 0, 1)
	return buildRange(start, end)
}

// LastMonthRange 返回上个月的时间范围。
func LastMonthRange() (DateTime, DateTime) {
	start := beginningOfMonth(Now().AddDate(0, -1, 0))
	end := beginningOfMonth(Now())
	return buildRange(start, end)
}

// CurrentMonthRange 返回本月的时间范围。
func CurrentMonthRange() (DateTime, DateTime) {
	start := beginningOfMonth(Now())
	end := beginningOfMonth(Now().AddDate(0, 1, 0))
	return buildRange(start, end)
}

// NextMonthRange 返回下个月的时间范围。
func NextMonthRange() (DateTime, DateTime) {
	start := beginningOfMonth(Now().AddDate(0, 1, 0))
	end := beginningOfMonth(Now().AddDate(0, 2, 0))
	return buildRange(start, end)
}

// LastYearRange 返回上一年的时间范围。
func LastYearRange() (DateTime, DateTime) {
	start := beginningOfYear(Now().AddDate(-1, 0, 0))
	end := beginningOfYear(Now())
	return buildRange(start, end)
}

// CurrentYearRange 返回今年的时间范围。
func CurrentYearRange() (DateTime, DateTime) {
	start := beginningOfYear(Now())
	end := beginningOfYear(Now().AddDate(1, 0, 0))
	return buildRange(start, end)
}

// NextYearRange 返回明年的时间范围。
func NextYearRange() (DateTime, DateTime) {
	start := beginningOfYear(Now().AddDate(1, 0, 0))
	end := beginningOfYear(Now().AddDate(2, 0, 0))
	return buildRange(start, end)
}

// ChineseWeekday 返回中文星期表示。
func ChineseWeekday(t time.Time) string {
	weekdays := []string{"日", "一", "二", "三", "四", "五", "六"}
	return "星期" + weekdays[t.Weekday()]
}

// EnglishWeekday 返回英文星期字符串。
func EnglishWeekday(t time.Time) string {
	return t.Weekday().String()
}

// IsWeekend 判断是否为周末。
func IsWeekend(t time.Time) bool {
	weekday := t.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

// TimeRange 表示一个具备唯一 ID 的时间段。
type TimeRange struct {
	ID    uint64
	Start time.Time
	End   time.Time
}

// IsValid 判断时间段的开始是否早于结束。
func (tr TimeRange) IsValid() bool {
	return !tr.Start.After(tr.End)
}

// Overlaps 判断两个时间段是否交叉。
func (tr TimeRange) Overlaps(other TimeRange) bool {
	if !tr.IsValid() || !other.IsValid() {
		return false
	}
	return tr.Start.Before(other.End) && tr.End.After(other.Start)
}

// TimeRangesConflict 检查多个时间段是否存在冲突。
func TimeRangesConflict(ranges ...TimeRange) bool {
	if len(ranges) < 2 {
		return false
	}
	for i := 0; i < len(ranges); i++ {
		for j := i + 1; j < len(ranges); j++ {
			if ranges[i].Overlaps(ranges[j]) {
				return true
			}
		}
	}
	return false
}

// ConflictTimeRangeIDs 返回存在冲突的时间段 ID 集合。
func ConflictTimeRangeIDs(ranges ...TimeRange) []uint64 {
	overlapping := make(map[uint64]struct{})
	for i := 0; i < len(ranges); i++ {
		for j := i + 1; j < len(ranges); j++ {
			if ranges[i].Overlaps(ranges[j]) {
				overlapping[ranges[i].ID] = struct{}{}
				overlapping[ranges[j].ID] = struct{}{}
			}
		}
	}
	var ids []uint64
	for id := range overlapping {
		ids = append(ids, id)
	}
	return ids
}

func beginningOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, ShangHaiTimeLocation)
}

func beginningOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, ShangHaiTimeLocation)
}

func beginningOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, ShangHaiTimeLocation)
}

func buildRange(start, end time.Time) (DateTime, DateTime) {
	return DateTime(start), DateTime(end)
}
