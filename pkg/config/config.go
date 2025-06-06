package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

// AppConfig holds the application configuration.
// We'll start with Kafka settings.
type AppConfig struct {
	Kafka KafkaConfig
}

// KafkaConfig holds Kafka related configuration.
type KafkaConfig struct {
	Brokers []string
	Topic   string
}

// Global application config instance
var Cfg AppConfig

// LoadConfig loads configuration from file and environment variables.
func LoadConfig(configPaths ...string) error {
	v := viper.New()
	v.SetConfigName("config") // name of config file (without extension)
	v.SetConfigType("yaml")   // or json, toml, etc.

	// Add default config paths
	if len(configPaths) == 0 {
		v.AddConfigPath("./conf") // path to look for the config file in
		v.AddConfigPath(".")      // optionally look for config in the working directory
	} else {
		for _, path := range configPaths {
			v.AddConfigPath(path)
		}
	}

	v.AutomaticEnv() // Read in environment variables that match
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // Replace . with _ in env var names (e.g. kafka.brokers -> KAFKA_BROKERS)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found; using default values or environment variables if set.")
			// You could set default values here if needed:
			// v.SetDefault("kafka.brokers", []string{"localhost:9092"})
			// v.SetDefault("kafka.topic", "case-events")
		} else {
			log.Printf("Error reading config file: %s", err)
			return err
		}
	}

	err := v.Unmarshal(&Cfg)
	if err != nil {
		log.Printf("Unable to decode into struct: %v", err)
		return err
	}

	// If Kafka brokers or topic are not set via config file, try to load from ENV as a fallback
	// This provides flexibility if the config file is missing but ENVs are present.
	if len(Cfg.Kafka.Brokers) == 0 {
		brokersStr := viper.GetString("KAFKA_BROKERS") // e.g., "broker1:9092,broker2:9092"
		if brokersStr != "" {
			Cfg.Kafka.Brokers = strings.Split(brokersStr, ",")
		} else {
			// Set a hardcoded default if nothing is found
			Cfg.Kafka.Brokers = []string{"localhost:9092"}
			log.Println("Kafka brokers not found in config or ENV, using default:", Cfg.Kafka.Brokers)
		}
	}
	if Cfg.Kafka.Topic == "" {
		Cfg.Kafka.Topic = viper.GetString("KAFKA_TOPIC")
		if Cfg.Kafka.Topic == "" {
			// Set a hardcoded default if nothing is found
			Cfg.Kafka.Topic = "case-events"
			log.Println("Kafka topic not found in config or ENV, using default:", Cfg.Kafka.Topic)
		}
	}

	log.Printf("Configuration loaded: %+v", Cfg)
	if len(Cfg.Kafka.Brokers) > 0 {
		log.Printf("Kafka Brokers: %v", Cfg.Kafka.Brokers)
	}
	log.Printf("Kafka Topic: %s", Cfg.Kafka.Topic)

	return nil
}
