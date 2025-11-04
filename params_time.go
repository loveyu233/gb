package gb

type ReqDateTimeStartEnd struct {
	StartDateTimeStr string `json:"start_date_time" form:"start_date_time"`
	EndDateTimeStr   string `json:"end_date_time" form:"end_date_time"`

	StartDateTime DateTime `json:"-" form:"-"`
	EndDateTime   DateTime `json:"-" form:"-"`

	DateTimeFilter bool `json:"-"`
}

func (req ReqDateTimeStartEnd) Parse() (ReqDateTimeStartEnd, error) {
	var data ReqDateTimeStartEnd
	if req.StartDateTimeStr != "" {
		s, err := StringToGbDateTimeErr(req.StartDateTimeStr)
		if err != nil {
			return data, err
		}
		data.StartDateTime = s
	}

	if req.EndDateTimeStr != "" {
		e, err := StringToGbDateTimeErr(req.EndDateTimeStr)
		if err != nil {
			return data, err
		}
		data.StartDateTime = e
	}

	if req.StartDateTimeStr != "" && req.EndDateTimeStr != "" {
		data.DateTimeFilter = true
	}

	return data, nil
}

type ReqDateTime struct {
	DateTimeStr string   `json:"date_time" form:"date_time"`
	DateTime    DateTime `json:"-" form:"-"`
}

func (req ReqDateTime) Parse() (ReqDateTime, error) {
	var data ReqDateTime
	if req.DateTimeStr != "" {
		s, err := StringToGbDateTimeErr(req.DateTimeStr)
		if err != nil {
			return data, err
		}
		data.DateTime = s
	}

	return data, nil
}

type ReqDateStartEnd struct {
	StartDateStr string `json:"start_date" form:"start_date"`
	EndDateStr   string `json:"end_date" form:"end_date"`

	StartDate DateOnly `json:"-" form:"-"`
	EndDate   DateOnly `json:"-" form:"-"`

	DateFilter bool `json:"-"`
}

func (req ReqDateStartEnd) Parse() (ReqDateStartEnd, error) {
	var data ReqDateStartEnd
	if req.StartDateStr != "" {
		s, err := StringToGBDateOnly(req.StartDateStr)
		if err != nil {
			return data, err
		}
		data.StartDate = s
	}

	if req.EndDateStr != "" {
		e, err := StringToGBDateOnly(req.EndDateStr)
		if err != nil {
			return data, err
		}
		data.EndDate = e
	}

	if req.StartDateStr != "" && req.EndDateStr != "" {
		data.DateFilter = true
	}

	return data, nil
}

type ReqDate struct {
	DateStr string   `json:"date" form:"date"`
	Date    DateOnly `json:"-" form:"-"`
}

func (req ReqDate) Parse() (ReqDate, error) {
	var data ReqDate
	if req.DateStr != "" {
		s, err := StringToGBDateOnly(req.DateStr)
		if err != nil {
			return data, err
		}
		data.Date = s
	}
	return data, nil
}

type ReqTimeStartEnd struct {
	StartTimeStr string `json:"start_time" form:"start_time"`
	EndTimeStr   string `json:"end_time" form:"end_time"`

	StartTime TimeOnly `json:"-" form:"-"`
	EndTime   TimeOnly `json:"-" form:"-"`

	TimeFilter bool `json:"-"`
}

func (req ReqTimeStartEnd) Parse() (ReqTimeStartEnd, error) {
	var data ReqTimeStartEnd
	if req.StartTimeStr != "" {
		s, err := StringToGBTimeOnly(req.StartTimeStr)
		if err != nil {
			return data, err
		}
		data.StartTime = s
	}

	if req.EndTimeStr != "" {
		e, err := StringToGBTimeOnly(req.EndTimeStr)
		if err != nil {
			return data, err
		}
		data.EndTime = e
	}

	if req.StartTimeStr != "" && req.EndTimeStr != "" {
		data.TimeFilter = true
	}

	return data, nil
}

type ReqTime struct {
	TimeStr string   `json:"time" form:"time"`
	Time    TimeOnly `json:"-" form:"-"`
}

func (req ReqTime) Parse() (ReqTime, error) {
	var data ReqTime
	if req.TimeStr != "" {
		s, err := StringToGBTimeOnly(req.TimeStr)
		if err != nil {
			return data, err
		}
		data.Time = s
	}
	return data, nil
}

type ReqTimeHourMinuteStartEnd struct {
	StartTimeHourMinuteStr string `json:"start_time_hour_minute" form:"start_time_hour_minute"`
	EndTimeHourMinuteStr   string `json:"end_time_hour_minute" form:"end_time_hour_minute"`

	StartTimeHourMinute TimeHourMinute `json:"-" form:"-"`
	EndTimeHourMinute   TimeHourMinute `json:"-" form:"-"`

	TimeHourMinuteFilter bool `json:"-"`
}

func (req ReqTimeHourMinuteStartEnd) Parse() (ReqTimeHourMinuteStartEnd, error) {
	var data ReqTimeHourMinuteStartEnd
	if req.StartTimeHourMinuteStr != "" {
		s, err := StringToGBTimeHourMinute(req.StartTimeHourMinuteStr)
		if err != nil {
			return data, err
		}
		data.StartTimeHourMinute = s
	}

	if req.EndTimeHourMinuteStr != "" {
		e, err := StringToGBTimeHourMinute(req.EndTimeHourMinuteStr)
		if err != nil {
			return data, err
		}
		data.EndTimeHourMinute = e
	}

	if req.StartTimeHourMinuteStr != "" && req.EndTimeHourMinuteStr != "" {
		data.TimeHourMinuteFilter = true
	}

	return data, nil
}

type ReqTimeHourMinute struct {
	TimeHourMinuteStr string         `json:"time_hour_minute" form:"time_hour_minute"`
	TimeHourMinute    TimeHourMinute `json:"-" form:"-"`
}

func (req ReqTimeHourMinute) Parse() (ReqTimeHourMinute, error) {
	var data ReqTimeHourMinute
	if req.TimeHourMinuteStr != "" {
		s, err := StringToGBTimeHourMinute(req.TimeHourMinuteStr)
		if err != nil {
			return data, err
		}
		data.TimeHourMinute = s
	}

	return data, nil
}
