package app_repository_postgres

import "time"

type Config struct {
	User            string
	Password        string
	Host            string
	Port            string
	DataBase        string
	MaxConns        int
	MinConns        int
	MaxConnLifeTime time.Duration
	MaxConnIdleTime time.Duration
}

func NewConfig() Config {
	return Config{
		User:            "postgres",
		Password:        "root",
		Host:            "localhost",
		Port:            "5434",
		DataBase:        "postgres",
		MaxConns:        25,
		MinConns:        5,
		MaxConnLifeTime: 5 * time.Minute,
		MaxConnIdleTime: 2 * time.Minute,
	}
}
