package gb

type GenIDData struct {
	ID   uint64 `gorm:"column:id"`
	Data string `gorm:"column:data"`
}

type GenIDDateTime struct {
	ID       uint64   `gorm:"column:id"`
	DateTime DateTime `gorm:"column:date_time"`
}

type GenIDPtrDateTime struct {
	ID       uint64    `gorm:"column:id"`
	DateTime *DateTime `gorm:"column:date_time"`
}

type GenIDDateOnly struct {
	ID       uint64   `gorm:"column:id"`
	DateOnly DateOnly `gorm:"column:date_only"`
}

type GenIDPtrDateOnly struct {
	ID       uint64    `gorm:"column:id"`
	DateOnly *DateOnly `gorm:"column:date_only"`
}

type GenIDTimeOnly struct {
	ID       uint64   `gorm:"column:id"`
	TimeOnly TimeOnly `gorm:"column:time_only"`
}

type GenIDPtrTimeOnly struct {
	ID       uint64    `gorm:"column:id"`
	TimeOnly *TimeOnly `gorm:"column:time_only"`
}

type GenIDTimeHourMinute struct {
	ID             uint64         `gorm:"column:id"`
	TimeHourMinute TimeHourMinute `gorm:"column:time_hour_minute"`
}

type GenIDPtrTimeHourMinute struct {
	ID             uint64          `gorm:"column:id"`
	TimeHourMinute *TimeHourMinute `gorm:"column:time_hour_minute"`
}

type GenIDDataTwo struct {
	ID    uint64 `gorm:"column:id"`
	Data1 string `gorm:"column:data_1"`
	Data2 string `gorm:"column:data_2"`
}
