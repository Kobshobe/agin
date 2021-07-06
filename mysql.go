package agin

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type Mysql struct {
	Path            string
	Config          string `mapstructure:"config" json:"config" yaml:"config"`
	Dbname          string `mapstructure:"db-name" json:"db-name" yaml:"db-name"`
	Username        string `mapstructure:"username" json:"username" yaml:"username"`
	Password        string `mapstructure:"password" json:"password" yaml:"password"`
	MaxIdleConns    int    `mapstructure:"max-idle-conns" json:"maxIdleConns" yaml:"max-idle-conns"`
	MaxOpenConns    int    `mapstructure:"max-open-conns" json:"maxOpenConns" yaml:"max-open-conns"`
	LogMode         bool   `mapstructure:"log-mode" json:"logMode" yaml:"log-mode"`
	LogZap          string `mapstructure:"log-zap" json:"logZap" yaml:"log-zap"`
	LocalPath       string `mapstructure:"localPath" json:"localPath" yaml:"localPath"`
	CloudPath       string `mapstructure:"cloudPath" json:"cloudPath" yaml:"cloudPath"`
	DockerPath      string `mapstructure:"dockerPath" json:"dockerPath" yaml:"dockerPath"`
	InnerDockerPath string `mapstructure:"innerDockerPath" json:"innerDockerPath" yaml:"innerDockerPath"`
	LocalPwd        string `mapstructure:"localPwd" json:"localPwd" yaml:"localPwd"`
	DockerPwd       string `mapstructure:"dockerPwd" json:"dockerPwd" yaml:"dockerPwd"`
}

// 获取dsn配置
func (config *Mysql) DSN(mode string) string {
	if mode == "local" || mode == "test" {
		config.Path = config.LocalPath
		config.Password = config.LocalPwd
	} else if mode == "docker" {
		config.Path = config.DockerPath
		config.Password = config.DockerPwd
	} else if mode == "cloud" {
		config.Path = config.CloudPath
		config.Password = config.DockerPwd
	} else if mode == "dockerInner" {
		config.Path = config.InnerDockerPath
		config.Password = config.DockerPwd
	}
	return config.Username + ":" + config.Password + "@tcp(" + config.Path + ")/" + config.Dbname + "?" + config.Config
}

func (config *Mysql) InitDB(mode string) *gorm.DB {

	fmt.Println(config.DSN(mode), config.Username)

	mysqlConfig := mysql.Config{
		DSN: config.DSN(mode), // DSN data source name
		//DefaultStringSize:         191,   	// string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: true, // 根据版本自动配置
	}
	fmt.Println("db begin connect ...")
	if db, err := gorm.Open(mysql.New(mysqlConfig), gormConfig(true)); err != nil {
		log.Fatal("db connect fail: ", err)
		return nil
	} else {
		fmt.Println("db connect in ...")
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(config.MaxIdleConns)
		sqlDB.SetMaxOpenConns(config.MaxOpenConns)
		fmt.Println("db connect ok ...")
		return db
	}
}

//@author: SliverHorn
//@function: gormConfig
//@description: 根据配置决定是否开启日志
//@param: mod bool
//@return: *gorm.Config

func gormConfig(mod bool) *gorm.Config {
	var config = &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true}
	//switch global.CONFIG.Mysql.LogZap {
	//case "silent", "Silent":
	//	config.Logger = internal.Default.LogMode(logger.Silent)
	//case "error", "Error":
	//	config.Logger = internal.Default.LogMode(logger.Error)
	//case "warn", "Warn":
	//	config.Logger = internal.Default.LogMode(logger.Warn)
	//case "info", "Info":
	//	config.Logger = internal.Default.LogMode(logger.Info)
	//case "zap", "Zap":
	//	config.Logger = internal.Default.LogMode(logger.Info)
	//default:
	//	if mod {
	//		config.Logger = internal.Default.LogMode(logger.Info)
	//		break
	//	}
	//	config.Logger = internal.Default.LogMode(logger.Silent)
	//}
	return config
}