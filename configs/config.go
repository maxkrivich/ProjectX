package configs

import (
	"fmt"
	"log"
	"os"
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

type MinioConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
}

type Config struct {
	DBConfig
	ServerConfig
	MinioConfig
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

func newMinioConfig() MinioConfig {
	return MinioConfig{
		Endpoint:        getEnvVal("MINIO_ENDPOINT", false),
		AccessKeyID:     getEnvVal("MINIO_ACCESS_KEY", false),
		SecretAccessKey: getEnvVal("MINIO_SECRET_KEY", false),
		UseSSL:          false,
	}
}

func GetConfig() *Config {
	if config == nil {
		config = &Config{newDBConfig(), newServerConfig(), newMinioConfig()}
	}
	return config
}

func (db DBConfig) GetConnectionString() string {
	return fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", db.Host, db.Username, db.DBName, db.Password, db.Port)
}
