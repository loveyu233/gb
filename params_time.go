package gb

import (
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

type ReqDateTimeStartEnd struct {
	StartDateTimeStr string `json:"start_date_time" form:"start_date_time"`
	EndDateTimeStr   string `json:"end_date_time" form:"end_date_time"`

	StartDateTime DateTime `json:"-" form:"-"`
	EndDateTime   DateTime `json:"-" form:"-"`

	DateTimeFilter bool `json:"-"`
}

// Parse 方法用于处理Parse相关逻辑。
func (req *ReqDateTimeStartEnd) Parse() error {
	start, hasStart, err := parseOptional(req.StartDateTimeStr, ParseDateTimeValue)
	if err != nil {
		return err
	}
	end, hasEnd, err := parseOptional(req.EndDateTimeStr, ParseDateTimeValue)
	if err != nil {
		return err
	}
	if hasStart {
		req.StartDateTime = start
	}
	if hasEnd {
		req.EndDateTime = end
	}
	req.DateTimeFilter = hasStart && hasEnd
	return nil
}

// Enabled 判断是否具备完整的时间范围过滤条件。
func (req ReqDateTimeStartEnd) Enabled() bool {
	return req.DateTimeFilter
}

// Scope 根据字段名生成对应的 GORM Scope。
func (req ReqDateTimeStartEnd) Scope(column field.IColumnName) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if !req.Enabled() {
			return db
		}
		return db.Where(column.ColumnName()+" BETWEEN ? AND ?", req.StartDateTime.Time(), req.EndDateTime.Time())
	}
}

type ReqDateTime struct {
	DateTimeStr string   `json:"date_time" form:"date_time"`
	DateTime    DateTime `json:"-" form:"-"`
}

// Parse 方法用于处理Parse相关逻辑。
func (req *ReqDateTime) Parse() error {
	if req.DateTimeStr != "" {
		s, err := ParseDateTimeValue(req.DateTimeStr)
		if err != nil {
			return err
		}
		req.DateTime = s
	}

	return nil
}

type ReqDateStartEnd struct {
	StartDateStr string `json:"start_date" form:"start_date"`
	EndDateStr   string `json:"end_date" form:"end_date"`

	StartDate DateOnly `json:"-" form:"-"`
	EndDate   DateOnly `json:"-" form:"-"`

	DateFilter bool `json:"-"`
}

// Parse 方法用于处理Parse相关逻辑。
func (req *ReqDateStartEnd) Parse() error {
	start, hasStart, err := parseOptional(req.StartDateStr, ParseDateOnly)
	if err != nil {
		return err
	}
	end, hasEnd, err := parseOptional(req.EndDateStr, ParseDateOnly)
	if err != nil {
		return err
	}
	if hasStart {
		req.StartDate = start
	}
	if hasEnd {
		req.EndDate = end
	}
	req.DateFilter = hasStart && hasEnd
	return nil
}

// Enabled 判断是否启用日期范围过滤。
func (req ReqDateStartEnd) Enabled() bool {
	return req.DateFilter
}

// Scope 返回针对指定字段的日期范围 Scope。
func (req ReqDateStartEnd) Scope(column field.IColumnName) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if !req.Enabled() {
			return db
		}
		return db.Where(column.ColumnName()+" BETWEEN ? AND ?", req.StartDate.Time(), req.EndDate.Time())
	}
}

type ReqDate struct {
	DateStr string   `json:"date" form:"date"`
	Date    DateOnly `json:"-" form:"-"`
}

// Parse 方法用于处理Parse相关逻辑。
func (req *ReqDate) Parse() error {
	if value, ok, err := parseOptional(req.DateStr, ParseDateOnly); err != nil {
		return err
	} else if ok {
		req.Date = value
	}
	return nil
}

type ReqTimeStartEnd struct {
	StartTimeStr string `json:"start_time" form:"start_time"`
	EndTimeStr   string `json:"end_time" form:"end_time"`

	StartTime TimeOnly `json:"-" form:"-"`
	EndTime   TimeOnly `json:"-" form:"-"`

	TimeFilter bool `json:"-"`
}

// Parse 方法用于处理Parse相关逻辑。
func (req *ReqTimeStartEnd) Parse() error {
	start, hasStart, err := parseOptional(req.StartTimeStr, ParseTimeOnly)
	if err != nil {
		return err
	}
	end, hasEnd, err := parseOptional(req.EndTimeStr, ParseTimeOnly)
	if err != nil {
		return err
	}
	if hasStart {
		req.StartTime = start
	}
	if hasEnd {
		req.EndTime = end
	}
	req.TimeFilter = hasStart && hasEnd
	return nil
}

// Enabled 判断是否启用具体时间的范围过滤。
func (req ReqTimeStartEnd) Enabled() bool {
	return req.TimeFilter
}

// Scope 返回针对时分秒字段的范围查询 Scope。
func (req ReqTimeStartEnd) Scope(column field.IColumnName) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if !req.Enabled() {
			return db
		}
		return db.Where(column.ColumnName()+" BETWEEN ? AND ?", req.StartTime.Time(), req.EndTime.Time())
	}
}

type ReqTime struct {
	TimeStr string   `json:"time" form:"time"`
	Time    TimeOnly `json:"-" form:"-"`
}

// Parse 方法用于处理Parse相关逻辑。
func (req *ReqTime) Parse() error {
	if value, ok, err := parseOptional(req.TimeStr, ParseTimeOnly); err != nil {
		return err
	} else if ok {
		req.Time = value
	}
	return nil
}

type ReqTimeHourMinuteStartEnd struct {
	StartTimeHourMinuteStr string `json:"start_time_hour_minute" form:"start_time_hour_minute"`
	EndTimeHourMinuteStr   string `json:"end_time_hour_minute" form:"end_time_hour_minute"`

	StartTimeHourMinute TimeHourMinute `json:"-" form:"-"`
	EndTimeHourMinute   TimeHourMinute `json:"-" form:"-"`

	TimeHourMinuteFilter bool `json:"-"`
}

// Parse 方法用于处理Parse相关逻辑。
func (req *ReqTimeHourMinuteStartEnd) Parse() error {
	start, hasStart, err := parseOptional(req.StartTimeHourMinuteStr, ParseHourMinute)
	if err != nil {
		return err
	}
	end, hasEnd, err := parseOptional(req.EndTimeHourMinuteStr, ParseHourMinute)
	if err != nil {
		return err
	}
	if hasStart {
		req.StartTimeHourMinute = start
	}
	if hasEnd {
		req.EndTimeHourMinute = end
	}
	req.TimeHourMinuteFilter = hasStart && hasEnd
	return nil
}

// Enabled 判断时分范围过滤是否生效。
func (req ReqTimeHourMinuteStartEnd) Enabled() bool {
	return req.TimeHourMinuteFilter
}

// Scope 返回针对时分字段的范围查询 Scope。
func (req ReqTimeHourMinuteStartEnd) Scope(column field.IColumnName) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if !req.Enabled() {
			return db
		}
		return db.Where(column.ColumnName()+" BETWEEN ? AND ?", req.StartTimeHourMinute.Time(), req.EndTimeHourMinute.Time())
	}
}

type ReqTimeHourMinute struct {
	TimeHourMinuteStr string         `json:"time_hour_minute" form:"time_hour_minute"`
	TimeHourMinute    TimeHourMinute `json:"-" form:"-"`
}

// Parse 方法用于处理Parse相关逻辑。
func (req *ReqTimeHourMinute) Parse() error {
	if value, ok, err := parseOptional(req.TimeHourMinuteStr, ParseHourMinute); err != nil {
		return err
	} else if ok {
		req.TimeHourMinute = value
	}

	return nil
}

func parseOptional[T any](raw string, parse func(string) (T, error)) (value T, ok bool, err error) {
	if raw == "" {
		return value, false, nil
	}
	parsed, err := parse(raw)
	if err != nil {
		var zero T
		return zero, false, err
	}
	return parsed, true, nil
}
