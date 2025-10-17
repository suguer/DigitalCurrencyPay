package dao

import (
	"DigitalCurrency/internal/config"
	"DigitalCurrency/internal/model/mdb"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/gookit/color"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var Mdb *gorm.DB

func InitDatabase(conf config.DatabaseConf) {
	switch conf.Driver {
	case "mysql":
		MysqlInit(&conf.Mysql)
	case "sqllite":
		SqlliteInit(&conf.Sqllite.Path)
	default:
		panic("database driver not support")
	}
}

// MysqlInit 数据库初始化
func SqlliteInit(path *string) {
	db, err := gorm.Open(sqlite.Open(*path), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
	Mdb = db
	// migration()
}

func MysqlInit(conf *config.MysqlConf) {

	var err error
	host := conf.Host
	port := strconv.Itoa(conf.Port)
	database := conf.Database
	username := conf.Username
	password := conf.Password
	charset := conf.Charset
	dsn := strings.Join([]string{username, ":", password, "@tcp(", host, ":", port, ")/", database, "?charset=" + charset + "&parseTime=true", "&loc=Local"}, "")
	// ormLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
	// 	SlowThreshold:             20 * time.Second,
	// 	LogLevel:                  logger.Error,
	// 	IgnoreRecordNotFoundError: true,
	// 	Colorful:                  true,
	// })
	Mdb, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         128,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}), &gorm.Config{

		SkipDefaultTransaction: true,
		// Logger:                                   ormLogger,
		DisableForeignKeyConstraintWhenMigrating: true, // 取消外键索引
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	/*
		Mdb, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				//TablePrefix:   viper.GetString("mysql_table_prefix"),
				SingularTable: true,
			},
			Logger: logger.Default.LogMode(logger.Error),
		})*/
	if err != nil {
		panic(err)
	}
	// if conf.Debug {
	// 	Mdb = Mdb.Debug()
	// }
	sqlDB, err := Mdb.DB()
	if err != nil {
		color.Red.Printf("[store_db] mysql get DB,err=%s\n", err)
		panic(err)
	}
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(time.Second * 30)
	err = sqlDB.Ping()
	if err != nil {
		color.Red.Printf("[store_db] mysql connDB err:%s", err.Error())
		panic(err)
	}
	//log.Sugar.Debug("[store_db] mysql connDB success")
}

func Migration() {
	err := Mdb.AutoMigrate(
		&mdb.Transaction{},
		&mdb.Wallet{},
		&mdb.Deposit{},
		&mdb.User{},
	)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
}
