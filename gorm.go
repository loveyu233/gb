package gb

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type GormClient struct {
	*gorm.DB
}

type GormConnConfig struct {
	Username string
	Password string
	Host     string
	Port     int
	Database string
	Params   map[string]interface{} // 连接参数,默认添加charset=utf8和parseTime=true以及loc=Asia%2FShanghai
}

// InitGormDB gormLogger可以使用默认的GormDefaultLogger
func InitGormDB(gcc GormConnConfig, gormLogger logger.Interface, opt ...func(db *gorm.DB) error) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?", gcc.Username, gcc.Password, gcc.Host, gcc.Port, gcc.Database)
	if gcc.Params["charset"] == nil {
		dsn = fmt.Sprintf("%scharset=utf8", dsn)
	}
	if gcc.Params["parseTime"] == nil {
		dsn = fmt.Sprintf("%s&parseTime=true", dsn)
	}
	if gcc.Params["loc"] == nil {
		dsn = fmt.Sprintf("%s&loc=Asia%%2FShanghai", dsn)
	}
	for k, v := range gcc.Params {
		dsn = fmt.Sprintf("%s&%s=%v", dsn, k, v)
	}
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
