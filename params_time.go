package gb

type ReqDateTimeStartEnd struct {
	StartDateTimeStr string `json:"start_date_time" form:"start_date_time"`
	EndDateTimeStr   string `json:"end_date_time" form:"end_date_time"`

	StartDateTime DateTime `json:"-" form:"-"`
	EndDateTime   DateTime `json:"-" form:"-"`

	DateTimeFilter bool `json:"-"`
}

func (req *ReqDateTimeStartEnd) Parse() error {
	if req.StartDateTimeStr != "" {
		s, err := StringToGbDateTimeErr(req.StartDateTimeStr)
		if err != nil {
			return err
		}
		req.StartDateTime = s
	}

	if req.EndDateTimeStr != "" {
		e, err := StringToGbDateTimeErr(req.EndDateTimeStr)
		if err != nil {
			return err
		}
		req.EndDateTime = e
	}

	if req.StartDateTimeStr != "" && req.EndDateTimeStr != "" {
		req.DateTimeFilter = true
	}

	return nil
}

type ReqDateTime struct {
	DateTimeStr string   `json:"date_time" form:"date_time"`
	DateTime    DateTime `json:"-" form:"-"`
}

func (req *ReqDateTime) Parse() error {
	if req.DateTimeStr != "" {
		s, err := StringToGbDateTimeErr(req.DateTimeStr)
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

func (req *ReqDateStartEnd) Parse() error {
	if req.StartDateStr != "" {
		s, err := StringToGBDateOnly(req.StartDateStr)
		if err != nil {
			return err
		}
		req.StartDate = s
	}

	if req.EndDateStr != "" {
		e, err := StringToGBDateOnly(req.EndDateStr)
		if err != nil {
			return err
		}
		req.EndDate = e
	}

	if req.StartDateStr != "" && req.EndDateStr != "" {
		req.DateFilter = true
	}

	return nil
}

type ReqDate struct {
	DateStr string   `json:"date" form:"date"`
	Date    DateOnly `json:"-" form:"-"`
}

func (req *ReqDate) Parse() error {
	if req.DateStr != "" {
		s, err := StringToGBDateOnly(req.DateStr)
		if err != nil {
			return err
		}
		req.Date = s
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

func (req *ReqTimeStartEnd) Parse() error {
	if req.StartTimeStr != "" {
		s, err := StringToGBTimeOnly(req.StartTimeStr)
		if err != nil {
			return err
		}
		req.StartTime = s
	}

	if req.EndTimeStr != "" {
		e, err := StringToGBTimeOnly(req.EndTimeStr)
		if err != nil {
			return err
		}
		req.EndTime = e
	}

	if req.StartTimeStr != "" && req.EndTimeStr != "" {
		req.TimeFilter = true
	}

	return nil
}

type ReqTime struct {
	TimeStr string   `json:"time" form:"time"`
	Time    TimeOnly `json:"-" form:"-"`
}

func (req *ReqTime) Parse() error {
	if req.TimeStr != "" {
		s, err := StringToGBTimeOnly(req.TimeStr)
		if err != nil {
			return err
		}
		req.Time = s
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

func (req *ReqTimeHourMinuteStartEnd) Parse() error {
	if req.StartTimeHourMinuteStr != "" {
		s, err := StringToGBTimeHourMinute(req.StartTimeHourMinuteStr)
		if err != nil {
			return err
		}
		req.StartTimeHourMinute = s
	}

	if req.EndTimeHourMinuteStr != "" {
		e, err := StringToGBTimeHourMinute(req.EndTimeHourMinuteStr)
		if err != nil {
			return err
		}
		req.EndTimeHourMinute = e
	}

	if req.StartTimeHourMinuteStr != "" && req.EndTimeHourMinuteStr != "" {
		req.TimeHourMinuteFilter = true
	}

	return nil
}

type ReqTimeHourMinute struct {
	TimeHourMinuteStr string         `json:"time_hour_minute" form:"time_hour_minute"`
	TimeHourMinute    TimeHourMinute `json:"-" form:"-"`
}

func (req *ReqTimeHourMinute) Parse() error {
	if req.TimeHourMinuteStr != "" {
		s, err := StringToGBTimeHourMinute(req.TimeHourMinuteStr)
		if err != nil {
			return err
		}
		req.TimeHourMinute = s
	}

	return nil
}
