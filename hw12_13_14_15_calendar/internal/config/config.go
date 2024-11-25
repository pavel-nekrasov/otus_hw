package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Logger LoggerConf
	HTTP   HTTPConf
}

type LoggerConf struct {
	Level  string
	Output string
}

type HTTPConf struct {
	Host string
	Port int
}

func New(filePath string) (c Config) {
	_, err := toml.DecodeFile(filePath, &c)
	if err != nil {
		log.Fatalf("Failed to load config file: %v", err)
	}
	return
}
