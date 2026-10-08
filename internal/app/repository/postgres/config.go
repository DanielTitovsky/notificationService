package app_repository_postgres

import "time"

type Config struct {
	User            string        `env: User`
	Password        string        `env: Password`
	Host            string        `env: Host`
	Port            string        `env: Port`
	DataBase        string        `env: DataBase`
	MaxConns        int           `env: MaxConns`
	MinConns        int           `env: MinConns`
	MaxConnLifeTime time.Duration `env: MaxConnLifeTime`
	MaxConnIdleTime time.Duration `env: MaxConnIdleTime`
}

func NewConfig() Config {
	return Config{
		User:            "articles_user",
		Password:        "articles_password",
		Host:            "localhost",
		Port:            "5433",
		DataBase:        "articles_db",
		MaxConns:        25,
		MinConns:        5,
		MaxConnLifeTime: 5 * time.Minute,
		MaxConnIdleTime: 2 * time.Minute,
	}
}
