package dao

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"strings"
	"time"
)

var Db *gorm.DB
var DatabaseConfig = new(Database) //设置全局的引用型指针变量
type Database struct {
	userName string
	password string
	host     string
	port     int
	dbName   string
	MaxConn  int
	MaxOpen  int
	args     string
}

func GetConfig() *Database {
	viper.SetConfigFile("conf/settings.yaml")
	content, err := os.ReadFile("conf/settings.yaml")
	if err != nil {
		fmt.Println("os获取配置文件失败！")
	}
	err = viper.ReadConfig(strings.NewReader(os.ExpandEnv(string(content))))
	if err != nil {
		fmt.Println("viper获取配置文件失败！")
	}
	cfgDatabase := viper.Sub("datasource")
	DatabaseConfig = InitDatabase(cfgDatabase)
	return DatabaseConfig
}

// 对上面的配置
func InitDatabase(cfg *viper.Viper) *Database {
	db := &Database{
		userName: cfg.GetString("userName"),
		port:     cfg.GetInt("port"),
		password: cfg.GetString("password"),
		host:     cfg.GetString("host"),
		dbName:   cfg.GetString("dbName"),
		args:     cfg.GetString("args"),
		MaxConn:  cfg.GetInt("maxConn"),
		MaxOpen:  cfg.GetInt("maxOpen"),
	}
	return db
}

/*
	*type Tabler interface {
		TableName() string
	}
*/

// 定义数据库连接对象
func Init() {
	DatabaseConfig = GetConfig()
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: 2 * time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info,     // Log level
			Colorful:      true,            // 彩色打印
		},
	)
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		DatabaseConfig.userName, DatabaseConfig.password, DatabaseConfig.host, DatabaseConfig.port, DatabaseConfig.dbName, DatabaseConfig.args)
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                                   newLogger,
		SkipDefaultTransaction:                   true,
		PrepareStmt:                              true,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Panicln("err:", err.Error())
	}
	sqlDb, err := Db.DB()
	//defer sqlDb.Close()
	fmt.Println("maxConn: ", DatabaseConfig.MaxConn)
	fmt.Println("maxOpen: ", DatabaseConfig.MaxOpen)
	sqlDb.SetMaxIdleConns(DatabaseConfig.MaxConn) //设置最大连接数
	sqlDb.SetMaxOpenConns(DatabaseConfig.MaxOpen) //设置最大的空闲连接数
}
