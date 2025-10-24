package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type DBConfig struct {
	Host     string `json:"db_host"`
	Port     int    `json:"db_port"`
	Name     string `json:"db_name"`
	Username string `json:"db_username"`
	Password string `json:"db_password"`
}

var (
	DBConfigure DBConfig
	DB          *sql.DB
)

func InitDBConfig(configData map[string]interface{}) {
	DBConfigure = DBConfig{
		Host:     getConfigString(getJSONTag(DBConfig{}, "Host"), configData, "127.0.0.1"),
		Port:     getConfigInt(getJSONTag(DBConfig{}, "Port"), configData, 3306),
		Name:     getConfigString(getJSONTag(DBConfig{}, "Name"), configData, ""),
		Username: getConfigString(getJSONTag(DBConfig{}, "Username"), configData, "root"),
		Password: getConfigString(getJSONTag(DBConfig{}, "Password"), configData, ""),
	}
	initDBInstance()
}

func initDBInstance() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		DBConfigure.Username,
		DBConfigure.Password,
		DBConfigure.Host,
		DBConfigure.Port,
		DBConfigure.Name,
	)
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Database connection established successfully:", dsn)
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
