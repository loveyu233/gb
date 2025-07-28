package gb

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"time"
)

// DateOnly 表示只有日期的类型
type DateOnly time.Time

// NewDateString 根据日期字符串创建一个 DateOnly 实例
func NewDateString(dateString string) (*DateOnly, error) {
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		return nil, err
	}
	return (*DateOnly)(&date), nil
}

// Scan 实现 sql.Scanner 接口，处理从数据库读取的值
func (d *DateOnly) Scan(v interface{}) error {
	if v == nil {
		*d = DateOnly(time.Time{})
		return nil
	}

	switch value := v.(type) {
	case []byte:
		dateStr := string(value)
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return err
		}

		// 确保时间部分为零
		fixedDate := time.Date(
			parsedDate.Year(), parsedDate.Month(), parsedDate.Day(),
			0, 0, 0, 0,
			time.UTC,
		)

		*d = DateOnly(fixedDate)
		return nil

	case string:
		dateStr := value
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return err
		}

		// 确保时间部分为零
		fixedDate := time.Date(
			parsedDate.Year(), parsedDate.Month(), parsedDate.Day(),
			0, 0, 0, 0,
			time.UTC,
		)

		*d = DateOnly(fixedDate)
		return nil

	case time.Time:
		// 确保时间部分为零
		fixedDate := time.Date(
			value.Year(), value.Month(), value.Day(),
			0, 0, 0, 0,
			time.UTC,
		)

		*d = DateOnly(fixedDate)
		return nil
	}

	return errors.New("类型转换错误：不支持的日期格式")
}

// Value 实现 driver.Valuer 接口，准备写入数据库的值
func (d DateOnly) Value() (driver.Value, error) {
	tm := time.Time(d)
	if tm.IsZero() {
		return nil, nil
	}
	return tm.Format("2006-01-02"), nil
}

// 添加辅助方法以方便使用
func (d DateOnly) String() string {
	return time.Time(d).Format("2006-01-02")
}

// Format 允许自定义格式化输出
func (d DateOnly) Format(layout string) string {
	return time.Time(d).Format(layout)
}

// Time 返回对应的 time.Time 值
func (d DateOnly) Time() time.Time {
	return time.Time(d)
}

// MarshalJSON 实现 JSON 序列化
func (d DateOnly) MarshalJSON() ([]byte, error) {
	formatted := time.Time(d).Format("2006-01-02")
	return json.Marshal(formatted)
}

// UnmarshalJSON 实现 JSON 反序列化
func (d *DateOnly) UnmarshalJSON(data []byte) error {
	var dateStr string
	if err := json.Unmarshal(data, &dateStr); err != nil {
		return err
	}
	parsed, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return err
	}
	*d = DateOnly(parsed)
	return nil
}

func (d DateOnly) FormatRelativeDate() string {
	return FormatDateRelativeDate(d.Time())
}

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

type TimeOnly time.Time

func (t *TimeOnly) Scan(v interface{}) error {
	if v == nil {
		*t = TimeOnly(time.Time{})
		return nil
	}

	switch value := v.(type) {
	case []byte:
		timeStr := string(value)
		layout := "15:04:05"
		if len(timeStr) > 8 {
			layout = "15:04"
		}
		parsedTime, err := time.Parse(layout, timeStr)
		if err != nil {
			parsedTime, err = time.Parse("2006-01-02 15:04:05", "1970-01-01 "+timeStr)
			if err != nil {
				return err
			}
		}
		fixedTime := time.Date(
			1970, 1, 1,
			parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(), parsedTime.Nanosecond(),
			time.UTC,
		)

		*t = TimeOnly(fixedTime)
		return nil

	case string:
		timeStr := value
		layout := "15:04:05"
		if len(timeStr) > 8 {
			layout = "15:04"
		}

		parsedTime, err := time.Parse(layout, timeStr)
		if err != nil {
			parsedTime, err = time.Parse("2006-01-02 15:04:05", "1970-01-01 "+timeStr)
			if err != nil {
				return err
			}
		}

		fixedTime := time.Date(
			1970, 1, 1,
			parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(), parsedTime.Nanosecond(),
			time.UTC,
		)

		*t = TimeOnly(fixedTime)
		return nil

	case time.Time:
		fixedTime := time.Date(
			1970, 1, 1,
			value.Hour(), value.Minute(), value.Second(), value.Nanosecond(),
			time.UTC,
		)
		*t = TimeOnly(fixedTime)
		return nil
	}

	return errors.New("类型转换错误：不支持的时间格式")
}

func (t TimeOnly) Value() (driver.Value, error) {
	tm := time.Time(t)
	if tm.IsZero() {
		return nil, nil
	}
	return tm.Format("15:04:05"), nil
}

func (t TimeOnly) String() string {
	return time.Time(t).Format("15:04:05")
}

func (t TimeOnly) Format(layout string) string {
	return time.Time(t).Format(layout)
}

func (t TimeOnly) Time() time.Time {
	return time.Time(t)
}
func (t TimeOnly) MarshalJSON() ([]byte, error) {
	formatted := time.Time(t).Format("15:04:05")
	return json.Marshal(formatted)
}

func (t *TimeOnly) UnmarshalJSON(data []byte) error {
	var timeStr string
	if err := json.Unmarshal(data, &timeStr); err != nil {
		return err
	}
	parsed, err := time.Parse("15:04", timeStr)
	if err != nil {
		return err
	}
	*t = TimeOnly(parsed)
	return nil
}

func (t TimeOnly) AddTime(hours, minutes, seconds int) TimeOnly {
	tm := time.Time(t)

	// 添加指定的时间
	newTime := tm.Add(
		time.Duration(hours)*time.Hour +
			time.Duration(minutes)*time.Minute +
			time.Duration(seconds)*time.Second,
	)

	// 确保日期部分保持为 1970-01-01，只保留时间部分
	fixedTime := time.Date(
		1970, 1, 1,
		newTime.Hour(), newTime.Minute(), newTime.Second(), newTime.Nanosecond(),
		time.UTC,
	)

	return TimeOnly(fixedTime)
}

// Before 方法 - 判断当前时间是否早于另一个时间
func (t TimeOnly) Before(other TimeOnly) bool {
	return t.Time().Before(other.Time())
}

// After 方法 - 判断当前时间是否晚于另一个时间
func (t TimeOnly) After(other TimeOnly) bool {
	return t.Time().Before(other.Time())
}

// Equal 方法 - 判断两个时间是否相等
func (t TimeOnly) Equal(other TimeOnly) bool {
	return t.Time().Equal(other.Time())
}

// Sub 方法 - 计算两个时间的差值，返回 Duration
func (t TimeOnly) Sub(other TimeOnly) time.Duration {
	return t.Time().Sub(other.Time())
}

func (t TimeOnly) FormatRelativeDate() string {
	return FormatTimeRelativeDate(t.Time())
}

type DateTime time.Time

func (t *DateTime) Scan(v interface{}) error {
	if v == nil {
		*t = DateTime(time.Time{})
		return nil
	}

	*t = DateTime(cast.ToTime(v).In(cst))
	return nil
}

func (t DateTime) Value() (driver.Value, error) {
	tm := time.Time(t)
	if tm.IsZero() {
		return nil, nil
	}
	return tm.In(cst).Format(CSTLayout), nil
}

func (t DateTime) String() string {
	return time.Time(t).In(cst).Format(CSTLayout)
}

func (t DateTime) Format(layout string) string {
	return time.Time(t).In(cst).Format(layout)
}

func (t DateTime) Time() time.Time {
	return time.Time(t).In(cst)
}
func (t DateTime) MarshalJSON() ([]byte, error) {
	formatted := time.Time(t).In(cst).Format(CSTLayout)
	return json.Marshal(formatted)
}

func (t *DateTime) UnmarshalJSON(data []byte) error {
	var timeStr string
	if err := json.Unmarshal(data, &timeStr); err != nil {
		return err
	}
	parsed, err := time.ParseInLocation(CSTLayout, timeStr, cst)
	if err != nil {
		return err
	}
	*t = DateTime(parsed)
	return nil
}

type BoolType bool

func (t *BoolType) Scan(v interface{}) error {
	if v == nil {
		*t = false
		return nil
	}
	*t = BoolType(cast.ToBool(v))
	return nil
}

func (t BoolType) Value() (driver.Value, error) {
	return cast.ToBool(t), nil
}

func (t BoolType) String() string {
	return fmt.Sprintf("%v", cast.ToBool(t))
}

func (t BoolType) MarshalJSON() ([]byte, error) {
	return json.Marshal(cast.ToBool(t))
}

func (t *BoolType) UnmarshalJSON(data []byte) error {
	*t = BoolType(cast.ToBool(string(data)))
	return nil
}
