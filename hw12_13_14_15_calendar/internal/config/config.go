package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

const (
	StorageModePostgres = "postgres"
	StorageModeMemory   = "memory"
)

type CalendarConfig struct {
	Logger   LoggerConf
	Endpoint EndpointConf
	Storage  StorageConf
}

type SchedulerConfig struct {
	Logger  LoggerConf
	Storage StorageConf
	Queue   QueueConfig
}

type SenderConfig struct {
	Logger LoggerConf
	Queue  QueueConfig
}

type LoggerConf struct {
	Level  string
	Output string
}

type EndpointConf struct {
	Host     string
	HTTPPort int
	GRPCPort int
}

type StorageConf struct {
	Mode     string
	Host     string
	Port     int
	DBName   string
	User     string
	Password string
}

type QueueConfig struct {
	Host     string
	Port     int
	Exchnage string
	Queue    string
	User     string
	Password string
}

func NewCalendarConfig(filePath string) (c CalendarConfig) {
	_, err := toml.DecodeFile(filePath, &c)
	if err != nil {
		log.Fatalf("Failed to load config file: %v", err)
	}
	return c
}

func NewSchedulerConfig(filePath string) (c SchedulerConfig) {
	_, err := toml.DecodeFile(filePath, &c)
	if err != nil {
		log.Fatalf("Failed to load config file: %v", err)
	}
	return c
}

func NewSenderConfig(filePath string) (c SenderConfig) {
	_, err := toml.DecodeFile(filePath, &c)
	if err != nil {
		log.Fatalf("Failed to load config file: %v", err)
	}
	return c
}
