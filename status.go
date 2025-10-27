package gb

type StructGormIDData struct {
	ID   uint64 `gorm:"column:id"`
	Data string `gorm:"column:data"`
}

type StructGormIDDateTime struct {
	ID       uint64   `gorm:"column:id"`
	DateTime DateTime `gorm:"column:date_time"`
}

type StructGormIDPtrDateTime struct {
	ID       uint64    `gorm:"column:id"`
	DateTime *DateTime `gorm:"column:date_time"`
}

type StructGormIDDateOnly struct {
	ID       uint64   `gorm:"column:id"`
	DateOnly DateOnly `gorm:"column:date_only"`
}

type StructGormIDPtrDateOnly struct {
	ID       uint64    `gorm:"column:id"`
	DateOnly *DateOnly `gorm:"column:date_only"`
}

type StructGormIDTimeOnly struct {
	ID       uint64   `gorm:"column:id"`
	TimeOnly TimeOnly `gorm:"column:time_only"`
}

type StructGormIDPtrTimeOnly struct {
	ID       uint64    `gorm:"column:id"`
	TimeOnly *TimeOnly `gorm:"column:time_only"`
}

type StructGormIDTimeHourMinute struct {
	ID             uint64         `gorm:"column:id"`
	TimeHourMinute TimeHourMinute `gorm:"column:time_hour_minute"`
}

type StructGormIDPtrTimeHourMinute struct {
	ID             uint64          `gorm:"column:id"`
	TimeHourMinute *TimeHourMinute `gorm:"column:time_hour_minute"`
}

type StructGormIDDataTwo struct {
	ID    uint64 `gorm:"column:id"`
	Data1 string `gorm:"column:data_1"`
	Data2 string `gorm:"column:data_2"`
}
