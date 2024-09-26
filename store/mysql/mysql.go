package mysql

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/kataras/golog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"irir-layout/config"
	"irir-layout/internal/model"
)

var (
	mysqlDB *gorm.DB
	once    sync.Once
	logMap  = map[string]logger.LogLevel{
		"info":   logger.Info,
		"error":  logger.Error,
		"warn":   logger.Warn,
		"silent": logger.Silent,
	}
)

func DialToMysql(op *config.MySQL) {
	var dbIns *gorm.DB
	once.Do(func() {
		err := createDB(op)
		if err != nil {
			golog.Errorf("---> [MYSQL] Database %s creation failure", op.Database)
			golog.Errorf("---> [MYSQL] %s", err.Error())
			os.Exit(1)
		}

		newLogger := logger.New(golog.New(), logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logMap[op.LogLevel],
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		})

		dsn := fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=%t&loc=%s`,
			op.Username,
			op.Password,
			op.Host,
			op.Database,
			true,
			"Local")
		dbIns, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})

		var sqlDB *sql.DB

		sqlDB, err = dbIns.DB()
		sqlDB.SetMaxOpenConns(op.MaxOpenConnections)
		sqlDB.SetConnMaxLifetime(op.MaxConnectionLifeTime)
		sqlDB.SetMaxIdleConns(op.MaxIdleConnections)
		err = sqlDB.Ping()
		if err != nil {
			golog.Error("---> [MYSQL] mysql connection failure: %s", err)
			os.Exit(0)
		}

		mysqlDB = dbIns
	})

	if mysqlDB == nil {
		golog.Errorf("---> [MYSQL] failed to get mysql store: %+v", mysqlDB)
		os.Exit(1)
	}
}

// createDB 创建数据库
func createDB(opts *config.MySQL) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/", opts.Username, opts.Password, opts.Host)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(db)

	createSql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci;", opts.Database)
	if err = db.Ping(); err != nil {
		return err
	}
	re, err := db.Exec(createSql)
	_, err = re.RowsAffected()
	return err
}

func GetDB() *gorm.DB {
	return mysqlDB
}

func Close() {
	if mysqlDB != nil {
		db, err := mysqlDB.DB()
		if err != nil {
			golog.Errorf("---> [MYSQL] close db failure, %s", err.Error())
			return
		}
		err = db.Close()
		if err != nil {
			golog.Errorf("---> [MYSQL] close db failure, %s", err.Error())
		}
		golog.Info("---> [MYSQL] close db failure")
	}
}

func CreateSuperUser(db *gorm.DB, cf *config.MySQL) {
	superUser := &model.User{}
	err := db.Where("account = ?", cf.SuperUser).Find(superUser).Error
	if err != nil {
		golog.Errorf("创建超级用户失败：%s", err)
		return
	}

	if superUser.ID == 0 {
		superUser.Account = cf.SuperUser
		superUser.Password = cf.SuperUserPwd
		superUser.IsActive = true
		superUser.Password, _ = superUser.Encrypt()
		result := db.Create(&superUser)
		if result.Error != nil {
			golog.Errorf("创建超级用户失败：%s", err)
			os.Exit(1)
		}
	}
}
