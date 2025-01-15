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
	Logger   LoggerConf
	Storage  StorageConf
	Queue    QueueProducerConf
	Schedule ScheduleConf
}

type SenderConfig struct {
	Logger LoggerConf
	Queue  QueueConsumerConf
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

type QueueServerConf struct {
	Host     string
	Port     int
	User     string
	Password string
}

type QueueConsumerConf struct {
	QueueServerConf
	Exchange   string
	Queue      string
	RoutingKey string
}

type QueueProducerConf struct {
	QueueServerConf
	Exchange   string
	RoutingKey string
}

type ScheduleConf struct {
	RetentionPeriod string
	Interval        string
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
