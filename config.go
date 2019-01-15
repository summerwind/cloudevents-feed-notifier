package main

type Config struct {
	Endpoint string   `json:"endpoint"`
	Interval string   `json:"interval"`
	Feeds    []string `json:"feeds"`
}

func NewConfig() *Config {
	return &Config{
		Interval: "30m",
		Feeds:    []string{},
	}
}
