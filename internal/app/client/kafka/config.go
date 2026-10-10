package app_client_kafka

type Config struct {
	addr           string
	consumerGroup  string
	consumerTopics string
}

func newConfig(addr string, consumerGroup string, consumerTopics string) Config {
	return Config{
		addr:           addr,
		consumerGroup:  consumerGroup,
		consumerTopics: consumerTopics,
	}
}
