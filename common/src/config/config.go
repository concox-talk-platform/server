/**
 * 公共配置类，用于加载数据库配置等
 * Author: tesion
 * Date: 20th March 2019
 */
package config

import (
	"fmt"
	"github.com/go-ini/ini"
)

type DBConfig struct {
	Host		string
	Port		int64
	User 		string
	Password	string
	DB 			string
	MaxConn		int64
	Driver 		string
}

func NewDBConfig() *DBConfig {
	return new(DBConfig)
}

func (cfg *DBConfig) LoadConfig(section, configPath string) error {
	if cfg == nil {
		return fmt.Errorf("config obj is null")
	}

	config, err := ini.Load(configPath)
	if err != nil {
		return err
	}

	cfg.Host = config.Section(section).Key("host").String()
	cfg.Port = config.Section(section).Key("port").MustInt64(9999)
	cfg.User = config.Section(section).Key("user").String()
	cfg.Password = config.Section(section).Key("password").String()
	cfg.DB = config.Section(section).Key("db").String()
	cfg.MaxConn = config.Section(section).Key("max_conn").MustInt64(1)
	cfg.Driver = config.Section(section).Key("driver").String()

	return nil
}

func (cfg *DBConfig) GetDSN() string {
	str := ""
	if cfg != nil {
		str = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB)
	}
	return str
}
