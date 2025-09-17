package gb

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cast"
)

// DateTime 完整的日期时间类型 (YYYY-MM-DD HH:MM:SS)
type DateTime time.Time

// DateOnly 只有日期的类型 (YYYY-MM-DD)
type DateOnly time.Time

// TimeOnly 只有时间的类型，包含秒 (HH:MM:SS)
type TimeOnly time.Time

// TimeHourMinute 只有小时分钟的时间类型 (HH:MM)
type TimeHourMinute time.Time

// ========== DateTime 完整日期时间类型 ==========

// NewDateTime 创建新的 DateTime 实例
func NewDateTime(t time.Time) DateTime {
	return DateTime(t.In(ShangHaiTimeLocation))
}

// NewDateTimeString 根据字符串创建 DateTime 实例
func NewDateTimeString(dateTimeString string) (*DateTime, error) {
	dt, err := time.ParseInLocation(CSTLayout, dateTimeString, ShangHaiTimeLocation)
	if err != nil {
		return nil, err
	}
	result := DateTime(dt)
	return &result, nil
}

func (dt *DateTime) Scan(v interface{}) error {
	if v == nil {
		*dt = DateTime(time.Time{})
		return nil
	}
	*dt = DateTime(cast.ToTime(v).In(ShangHaiTimeLocation))
	return nil
}

func (dt DateTime) Value() (driver.Value, error) {
	tm := time.Time(dt)
	if tm.IsZero() {
		return nil, nil
	}
	return tm.In(ShangHaiTimeLocation).Format(CSTLayout), nil
}

func (dt DateTime) String() string {
	return time.Time(dt).In(ShangHaiTimeLocation).Format(CSTLayout)
}

func (dt DateTime) Format(layout string) string {
	return time.Time(dt).In(ShangHaiTimeLocation).Format(layout)
}

func (dt DateTime) Time() time.Time {
	return time.Time(dt).In(ShangHaiTimeLocation)
}

func (dt DateTime) MarshalJSON() ([]byte, error) {
	formatted := time.Time(dt).In(ShangHaiTimeLocation).Format(CSTLayout)
	return json.Marshal(formatted)
}

func (dt *DateTime) UnmarshalJSON(data []byte) error {
	var timeStr string
	if err := json.Unmarshal(data, &timeStr); err != nil {
		return err
	}
	parsed, err := time.ParseInLocation(CSTLayout, timeStr, ShangHaiTimeLocation)
	if err != nil {
		return err
	}
	*dt = DateTime(parsed)
	return nil
}

func (dt DateTime) IsZero() bool {
	return time.Time(dt).IsZero()
}

func (dt DateTime) FormatRelativeDate() string {
	return FormatDateRelativeDate(dt.Time())
}

// ========== DateOnly 只有日期类型 ==========

// NewDateOnly 创建新的 DateOnly 实例
func NewDateOnly(year int, month time.Month, day int) DateOnly {
	t := time.Date(year, month, day, 0, 0, 0, 0, ShangHaiTimeLocation)
	return DateOnly(t)
}

// NewDateOnlyString 根据日期字符串创建 DateOnly 实例
func NewDateOnlyString(dateString string) (*DateOnly, error) {
	date, err := time.ParseInLocation("2006-01-02", dateString, ShangHaiTimeLocation)
	if err != nil {
		return nil, err
	}
	// 确保时间部分为零
	fixedDate := time.Date(
		date.Year(), date.Month(), date.Day(),
		0, 0, 0, 0,
		ShangHaiTimeLocation,
	)
	result := DateOnly(fixedDate)
	return &result, nil
}

func (d *DateOnly) Scan(v interface{}) error {
	if v == nil {
		*d = DateOnly(time.Time{})
		return nil
	}

	switch value := v.(type) {
	case []byte:
		dateStr := string(value)
		return d.parseAndSet(dateStr)
	case string:
		return d.parseAndSet(value)
	case time.Time:
		// 确保时间部分为零
		fixedDate := time.Date(
			value.Year(), value.Month(), value.Day(),
			0, 0, 0, 0,
			ShangHaiTimeLocation,
		)
		*d = DateOnly(fixedDate)
		return nil
	}

	return errors.New("类型转换错误：不支持的日期格式")
}

func (d *DateOnly) parseAndSet(dateStr string) error {
	parsedDate, err := time.ParseInLocation("2006-01-02", dateStr, ShangHaiTimeLocation)
	if err != nil {
		return err
	}

	// 确保时间部分为零
	fixedDate := time.Date(
		parsedDate.Year(), parsedDate.Month(), parsedDate.Day(),
		0, 0, 0, 0,
		ShangHaiTimeLocation,
	)

	*d = DateOnly(fixedDate)
	return nil
}

func (d DateOnly) Value() (driver.Value, error) {
	tm := time.Time(d)
	if tm.IsZero() {
		return nil, nil
	}
	return tm.In(ShangHaiTimeLocation).Format("2006-01-02"), nil
}

func (d DateOnly) String() string {
	return time.Time(d).In(ShangHaiTimeLocation).Format("2006-01-02")
}

func (d DateOnly) Format(layout string) string {
	return time.Time(d).In(ShangHaiTimeLocation).Format(layout)
}

func (d DateOnly) Time() time.Time {
	return time.Time(d).In(ShangHaiTimeLocation)
}

func (d DateOnly) MarshalJSON() ([]byte, error) {
	formatted := time.Time(d).In(ShangHaiTimeLocation).Format("2006-01-02")
	return json.Marshal(formatted)
}

func (d *DateOnly) UnmarshalJSON(data []byte) error {
	var dateStr string
	if err := json.Unmarshal(data, &dateStr); err != nil {
		return err
	}
	parsed, err := time.ParseInLocation("2006-01-02", dateStr, ShangHaiTimeLocation)
	if err != nil {
		return err
	}
	*d = DateOnly(parsed)
	return nil
}

func (d DateOnly) IsZero() bool {
	return time.Time(d).IsZero()
}

func (d DateOnly) FormatRelativeDate() string {
	return FormatDateRelativeDate(d.Time())
}

// AddDays 添加天数
func (d DateOnly) AddDays(days int) DateOnly {
	tm := time.Time(d).AddDate(0, 0, days)
	return DateOnly(tm)
}

// ========== TimeOnly 包含秒的时间类型 ==========

// NewTimeOnly 创建新的 TimeOnly 实例
func NewTimeOnly(hour, minute, second int) TimeOnly {
	t := time.Date(1970, 1, 1, hour, minute, second, 0, ShangHaiTimeLocation)
	return TimeOnly(t)
}

// NewTimeOnlyString 根据时间字符串创建 TimeOnly 实例
func NewTimeOnlyString(timeString string) (*TimeOnly, error) {
	t := &TimeOnly{}
	parsedTime, err := t.parseTimeString(timeString)
	if err != nil {
		return nil, err
	}
	result := TimeOnly(parsedTime)
	return &result, nil
}

func (t *TimeOnly) Scan(v interface{}) error {
	if v == nil {
		*t = TimeOnly(time.Time{})
		return nil
	}

	switch value := v.(type) {
	case []byte:
		timeStr := string(value)
		parsedTime, err := t.parseTimeString(timeStr)
		if err != nil {
			return err
		}
		*t = TimeOnly(parsedTime)
		return nil

	case string:
		parsedTime, err := t.parseTimeString(value)
		if err != nil {
			return err
		}
		*t = TimeOnly(parsedTime)
		return nil

	case time.Time:
		fixedTime := time.Date(
			1970, 1, 1,
			value.Hour(), value.Minute(), value.Second(), value.Nanosecond(),
			ShangHaiTimeLocation,
		)
		*t = TimeOnly(fixedTime)
		return nil
	}

	return errors.New("类型转换错误：不支持的时间格式")
}

func (t *TimeOnly) parseTimeString(timeStr string) (time.Time, error) {
	layouts := []string{
		"15:04:05", // HH:MM:SS
		"15:04",    // HH:MM
	}

	for _, layout := range layouts {
		parsedTime, err := time.ParseInLocation(layout, timeStr, ShangHaiTimeLocation)
		if err == nil {
			fixedTime := time.Date(
				1970, 1, 1,
				parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(), parsedTime.Nanosecond(),
				ShangHaiTimeLocation,
			)
			return fixedTime, nil
		}
	}

	return time.Time{}, fmt.Errorf("无法解析时间格式: %s", timeStr)
}

func (t TimeOnly) Value() (driver.Value, error) {
	tm := time.Time(t)
	if tm.IsZero() {
		return nil, nil
	}
	return tm.In(ShangHaiTimeLocation).Format("15:04:05"), nil
}

func (t TimeOnly) String() string {
	return time.Time(t).In(ShangHaiTimeLocation).Format("15:04:05")
}

func (t TimeOnly) Format(layout string) string {
	return time.Time(t).In(ShangHaiTimeLocation).Format(layout)
}

func (t TimeOnly) Time() time.Time {
	return time.Time(t).In(ShangHaiTimeLocation)
}

func (t TimeOnly) MarshalJSON() ([]byte, error) {
	formatted := time.Time(t).In(ShangHaiTimeLocation).Format("15:04:05")
	return json.Marshal(formatted)
}

func (t *TimeOnly) UnmarshalJSON(data []byte) error {
	var timeStr string
	if err := json.Unmarshal(data, &timeStr); err != nil {
		return err
	}

	parsed, err := t.parseTimeString(timeStr)
	if err != nil {
		return err
	}

	*t = TimeOnly(parsed)
	return nil
}

func (t TimeOnly) IsZero() bool {
	return time.Time(t).IsZero()
}

// AddTime 添加时间
func (t TimeOnly) AddTime(hours, minutes, seconds int) TimeOnly {
	tm := time.Time(t).Add(
		time.Duration(hours)*time.Hour +
			time.Duration(minutes)*time.Minute +
			time.Duration(seconds)*time.Second,
	)

	fixedTime := time.Date(
		1970, 1, 1,
		tm.Hour(), tm.Minute(), tm.Second(), tm.Nanosecond(),
		ShangHaiTimeLocation,
	)

	return TimeOnly(fixedTime)
}

// Before 判断当前时间是否早于另一个时间
func (t TimeOnly) Before(other TimeOnly) bool {
	return t.Time().Before(other.Time())
}

// After 判断当前时间是否晚于另一个时间
func (t TimeOnly) After(other TimeOnly) bool {
	return t.Time().After(other.Time())
}

// Equal 判断两个时间是否相等
func (t TimeOnly) Equal(other TimeOnly) bool {
	return t.Time().Equal(other.Time())
}

// Sub 计算两个时间的差值
func (t TimeOnly) Sub(other TimeOnly) time.Duration {
	return t.Time().Sub(other.Time())
}

// ========== TimeHourMinute 不包含秒的时间类型 ==========

// NewTimeHourMinute 创建新的 TimeHourMinute 实例
func NewTimeHourMinute(hour, minute int) TimeHourMinute {
	t := time.Date(1970, 1, 1, hour, minute, 0, 0, ShangHaiTimeLocation)
	return TimeHourMinute(t)
}

// NewTimeHourMinuteString 根据时间字符串创建 TimeHourMinute 实例
func NewTimeHourMinuteString(timeString string) (*TimeHourMinute, error) {
	t := &TimeHourMinute{}
	parsedTime, err := t.parseTimeString(timeString)
	if err != nil {
		return nil, err
	}
	result := TimeHourMinute(parsedTime)
	return &result, nil
}

func (t *TimeHourMinute) Scan(v interface{}) error {
	if v == nil {
		*t = TimeHourMinute(time.Time{})
		return nil
	}

	switch value := v.(type) {
	case []byte:
		timeStr := string(value)
		parsedTime, err := t.parseTimeString(timeStr)
		if err != nil {
			return err
		}
		*t = TimeHourMinute(parsedTime)
		return nil

	case string:
		parsedTime, err := t.parseTimeString(value)
		if err != nil {
			return err
		}
		*t = TimeHourMinute(parsedTime)
		return nil

	case time.Time:
		// 忽略秒和纳秒部分
		fixedTime := time.Date(
			1970, 1, 1,
			value.Hour(), value.Minute(), 0, 0,
			ShangHaiTimeLocation,
		)
		*t = TimeHourMinute(fixedTime)
		return nil
	}

	return errors.New("类型转换错误：不支持的时间格式")
}

func (t *TimeHourMinute) parseTimeString(timeStr string) (time.Time, error) {
	layouts := []string{
		"15:04",    // HH:MM (首选)
		"15:04:05", // HH:MM:SS (忽略秒部分)
	}

	for _, layout := range layouts {
		parsedTime, err := time.ParseInLocation(layout, timeStr, ShangHaiTimeLocation)
		if err == nil {
			// 始终忽略秒和纳秒部分
			fixedTime := time.Date(
				1970, 1, 1,
				parsedTime.Hour(), parsedTime.Minute(), 0, 0,
				ShangHaiTimeLocation,
			)
			return fixedTime, nil
		}
	}

	return time.Time{}, fmt.Errorf("无法解析时间格式: %s", timeStr)
}

func (t TimeHourMinute) Value() (driver.Value, error) {
	tm := time.Time(t)
	if tm.IsZero() {
		return nil, nil
	}
	return tm.In(ShangHaiTimeLocation).Format("15:04"), nil
}

func (t TimeHourMinute) String() string {
	return time.Time(t).In(ShangHaiTimeLocation).Format("15:04")
}

func (t TimeHourMinute) Format(layout string) string {
	return time.Time(t).In(ShangHaiTimeLocation).Format(layout)
}

func (t TimeHourMinute) Time() time.Time {
	return time.Time(t).In(ShangHaiTimeLocation)
}

func (t TimeHourMinute) MarshalJSON() ([]byte, error) {
	formatted := time.Time(t).In(ShangHaiTimeLocation).Format("15:04")
	return json.Marshal(formatted)
}

func (t *TimeHourMinute) UnmarshalJSON(data []byte) error {
	var timeStr string
	if err := json.Unmarshal(data, &timeStr); err != nil {
		return err
	}

	parsed, err := t.parseTimeString(timeStr)
	if err != nil {
		return err
	}

	*t = TimeHourMinute(parsed)
	return nil
}

func (t TimeHourMinute) IsZero() bool {
	return time.Time(t).IsZero()
}

// AddTime 添加时间（只支持小时和分钟）
func (t TimeHourMinute) AddTime(hours, minutes int) TimeHourMinute {
	tm := time.Time(t).Add(
		time.Duration(hours)*time.Hour +
			time.Duration(minutes)*time.Minute,
	)

	fixedTime := time.Date(
		1970, 1, 1,
		tm.Hour(), tm.Minute(), 0, 0,
		ShangHaiTimeLocation,
	)

	return TimeHourMinute(fixedTime)
}

// Before 判断当前时间是否早于另一个时间
func (t TimeHourMinute) Before(other TimeHourMinute) bool {
	return t.Time().Before(other.Time())
}

// After 判断当前时间是否晚于另一个时间
func (t TimeHourMinute) After(other TimeHourMinute) bool {
	return t.Time().After(other.Time())
}

// Equal 判断两个时间是否相等
func (t TimeHourMinute) Equal(other TimeHourMinute) bool {
	return t.Time().Equal(other.Time())
}

// Sub 计算两个时间的差值
func (t TimeHourMinute) Sub(other TimeHourMinute) time.Duration {
	return t.Time().Sub(other.Time())
}

// ========== 通用的 Slice 类型 ==========

type Slice[T any] []T

func (Slice[T]) GormDataType() string {
	return "json"
}

func (s Slice[T]) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

func (s *Slice[T]) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}

	if len(bytes) == 0 {
		return nil
	}

	return json.Unmarshal(bytes, s)
}

func (s *Slice[T]) MarshalJSON() ([]byte, error) {
	if s == nil || *s == nil {
		return []byte("null"), nil
	}
	return json.Marshal([]T(*s))
}

func (s *Slice[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*s = nil
		return nil
	}
	var tmp []T
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	*s = tmp
	return nil
}
