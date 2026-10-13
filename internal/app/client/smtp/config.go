package app_client_smtp

type Config struct {
	Host     string
	Port     int
	AuthPlan string
	UserName string
	Password string
}

func NewConfig() Config {
	return Config{
		Host:     host,
		Port:     port,
		AuthPlan: authPlan,
		UserName: userName,
		Password: password,
	}
}
