package utils

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver      string        `mapstructure:"DB_DRIVER"`
	DBSource      string        `mapstructure:"DB_SOURCE"`
	ServerAddress string        `mapstructure:"SERVER_ADDRESS"`
	Version       string        `mapstructure:"VERSION"`
	Secret        string        `mapstructure:"SECRET"` // This might be the TokenKey
	TokenKey      string        `mapstructure:"TOKEN_KEY"`
	TokenDuration time.Duration `mapstructure:"TOKEN_DURATION"`
}

func LoadConfig(path string, configName string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(configName)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Printf("Warning: Could not load config file '%s' from path '%s'. Error: %v. Using default values.", configName, path, err)

		// Populate with default values
		config = Config{
			DBDriver:      "postgres",
			DBSource:      "postgresql://root:secret@localhost:5432/image_processing_db?sslmode=disable",
			ServerAddress: "0.0.0.0:3000",
			Version:       "1.0.0", // Default version
			Secret:        "defaultsecretkey12345678901234", // Default secret
			TokenKey:      "defaultsecretkey12345678901234", // Default token key
			TokenDuration: time.Hour * 1,                 // Default token duration 1 hour
		}
		// Check if essential defaults are present
		if config.DBSource == "" || config.TokenKey == "" {
			return Config{}, fmt.Errorf("critical default values (DBSource, TokenKey) are not set in fallback logic")
		}
		return config, nil // Return nil error because we handled it with defaults
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Printf("Error unmarshalling config: %v. Check config file structure and values.", err)
		return Config{}, err
	}

	// Ensure essential values loaded from file are present
	if config.DBSource == "" || (config.TokenKey == "" && config.Secret == "") { // Check both TokenKey and Secret
		log.Printf("Warning: DBSOURCE or TOKEN_KEY/SECRET is empty after loading config file '%s'. Check your config file content.", configName)
		// Fallback to defaults for critical missing fields if file loaded but fields are empty
		// This part can be adjusted based on how strictly missing fields from file should be handled
		if config.DBSource == "" {
			log.Println("DBSOURCE missing from file, applying default.")
			config.DBSource = "postgresql://root:secret@localhost:5432/image_processing_db?sslmode=disable"
		}
		if config.TokenKey == "" && config.Secret == "" {
			log.Println("TOKEN_KEY/SECRET missing from file, applying default.")
			config.TokenKey = "defaultsecretkey12345678901234"
			config.Secret = "defaultsecretkey12345678901234"
		}
        if config.TokenDuration == 0 {
            log.Println("TOKEN_DURATION missing or zero, applying default (1 hour).")
            config.TokenDuration = time.Hour * 1
        }
		// Return error if still missing after trying to apply defaults for missing fields
		if config.DBSource == "" || (config.TokenKey == "" && config.Secret == "") {
			return Config{}, fmt.Errorf("critical values (DBSource, TokenKey/SECRET) are missing from loaded config file '%s' and defaults could not be applied", configName)
		}
	}
	return config, nil
}
