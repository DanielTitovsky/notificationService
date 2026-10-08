package app_transport_http_server

import "time"

type Config struct {
	Addr            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

func NewConfig() Config {
	return Config{
		Addr:            ":8086",
		ReadTimeout:     20 * time.Second,
		WriteTimeout:    40 * time.Second,
		ShutdownTimeout: 40 * time.Second,
	}
}
