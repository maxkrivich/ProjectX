package configs

import (
	"log"
	"os"
	"fmt"
)

type DBConfig struct {
	Dialect  string
	Host     string
	Port     string
	DBName   string
	Username string
	Password string
}

type ServerConfig struct {
	Port string
	Host string
}

type Config struct {
	DBConfig
	ServerConfig
}

var config *Config

func getEnvVal(key string, empty bool) string {
	res := os.Getenv(key)
	if len(res) == 0 && !empty {
		log.Fatal(fmt.Sprintf("Expected environment variable %s not set", key))
	}
	return res
}

func newDBConfig() DBConfig {
	return DBConfig{
		Dialect:  "postgres",
		Host:     getEnvVal("DB_HOST", false),
		Port:     getEnvVal("DB_PORT", false),
		DBName:   getEnvVal("DB_NAME", false),
		Username: getEnvVal("DB_USERNAME", false),
		Password: getEnvVal("DB_PASSWORD", true),
	}
}

func newServerConfig() ServerConfig {
	return ServerConfig{
		Port: getEnvVal("SERVER_PORT", false),
		Host: getEnvVal("SERVER_HOST", false),
	}
}

func GetConfig() *Config {
	if config == nil {
		config = &Config{newDBConfig(), newServerConfig()}
	}
	return config
}

func (db DBConfig) GetConnectionString() string {
	return fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", db.Host, db.Username, db.DBName, db.Password, db.Port)
}
