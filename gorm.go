package gb

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func InitGormDB(dsn string, gormLogger logger.Interface, opt ...func(db *gorm.DB) error) (*gorm.DB, error) {
	db, err := gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{
			Logger:                 gormLogger,
			TranslateError:         true,
			SkipDefaultTransaction: true,
			PrepareStmt:            true,
		},
	)
	if err != nil {
		return nil, err
	}

	for _, fn := range opt {
		if err := fn(db); err != nil {
			return nil, err
		}
	}

	return db, nil
}

func GormDefaultLogger(logLevel ...int) logger.Interface {
	var ll int
	if len(logLevel) > 0 && logLevel[0] >= 1 && logLevel[0] <= 4 {
		ll = logLevel[0]
	} else {
		ll = 4
	}
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Millisecond * 100,
			LogLevel:                  logger.LogLevel(ll),
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
}
