package configs

import (
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

func GetConfig() *Config {
	if config == nil {
		config = &Config{
			DBConfig{
				Dialect:  "postgres",
				Host:     os.Getenv("DB_HOST"),
				Port:     os.Getenv("DB_PORT"),
				DBName:   os.Getenv("DB_NAME"),
				Username: os.Getenv("DB_USERNAME"),
				Password: os.Getenv("DB_PASSWORD"),
			},
			ServerConfig{
				Port: os.Getenv("SERVER_PORT"),
				Host: os.Getenv("SERVER_HOST"),
			},
		}
	}
	return config
}

func (db DBConfig) GetConnectionString() string {
	return fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", db.Host, db.Username, db.DBName, db.Password, db.Port)
}
