package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func Init(path ...string) error {
	return godotenv.Load(path...)
}

type ConfigDB struct {
	host     string
	port     string
	user     string
	password string
	dbName   string
	sslMode  string
}

func LoadConfigDB() (*ConfigDB, error) {
	host, err := getEnv("POSTGRES_HOST")
	if err != nil {
		return nil, err
	}

	port, err := getEnv("POSTGRES_PORT")
	if err != nil {
		return nil, err
	}

	user, err := getEnv("POSTGRES_USER")
	if err != nil {
		return nil, err
	}

	password, err := getEnv("POSTGRES_PASSWORD")
	if err != nil {
		return nil, err
	}

	dbName, err := getEnv("POSTGRES_DB")
	if err != nil {
		return nil, err
	}

	sslMode, err := getEnv("SSLMode")
	if err != nil {
		return nil, err
	}

	return &ConfigDB{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		dbName:   dbName,
		sslMode:  sslMode,
	}, nil
}

func getEnv(key string) (string, error) {
	if value, ok := os.LookupEnv(key); ok {
		return value, nil
	}

	return "", fmt.Errorf("There is no value for key %s", key)
}

func (cfg ConfigDB) GetUser() string {
	return cfg.user
}

func (cfg ConfigDB) GetPassword() string {
	return cfg.password
}

func (cfg ConfigDB) GetHost() string {
	return cfg.host
}

func (cfg ConfigDB) GetPort() string {
	return cfg.port
}

func (cfg ConfigDB) GetDBName() string {
	return cfg.dbName
}

func (cfg ConfigDB) GetSSLMode() string {
	return cfg.sslMode
}
