package database

import "fmt"

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
	DbType   string
}

func (conf *DatabaseConfig) DSN() string {
	switch conf.DbType {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			conf.User, conf.Password, conf.Host, conf.Port, conf.DbName)
	case "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
			conf.Host, conf.Port, conf.User, conf.DbName, conf.Password)
	default:
		panic("unsupported database type")
	}
}
