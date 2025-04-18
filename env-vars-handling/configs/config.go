package configs

import (
	"fmt"
	"sync"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

// Config is a struct that will receive configuration options via environment variables.
type Config struct {
	App struct {
		Name     string `envconfig:"NAME"`
		URL      string `envconfig:"URL"`
		Host     string `envconfig:"HOST"`
		BasePath string `envconfig:"BASE_PATH"`
		Revision string `envconfig:"REVISION"`

		CORS struct {
			Enable           bool     `envconfig:"ENABLE"`
			AllowCredentials bool     `envconfig:"ALLOW_CREDENTIALS"`
			AllowedHeaders   []string `envconfig:"ALLOWED_HEADERS"`
			AllowedMethods   []string `envconfig:"ALLOWED_METHODS"`
			AllowedOrigins   []string `envconfig:"ALLOWED_ORIGINS"`
			MaxAgeSeconds    int      `envconfig:"MAX_AGE_SECONDS"`
		} `envconfig:"CORS"`
	} `envconfig:"APP"`

	DB struct {
		MySQL struct {
			Write struct {
				Host     string `envconfig:"HOST"`
				Port     string `envconfig:"PORT"`
				Name     string `envconfig:"NAME"`
				Username string `envconfig:"USER"`
				Password string `envconfig:"PASSWORD"`
				Timezone string `envconfig:"TIMEZONE"`
			} `envconfig:"WRITE"`
			Read struct {
				Host     string `envconfig:"HOST"`
				Port     string `envconfig:"PORT"`
				Name     string `envconfig:"NAME"`
				Username string `envconfig:"USER"`
				Password string `envconfig:"PASSWORD"`
				Timezone string `envconfig:"TIMEZONE"`
			} `envconfig:"READ"`
		} `envconfig:"MYSQL"`
	} `envconfig:"DB"`

	Server struct {
		Env                   string `envconfig:"ENV"`
		LogLevel              string `envconfig:"LOG_LEVEL"`
		Port                  string `envconfig:"PORT"`
		Host                  string `envconfig:"HOST"`
		ShutdownCleanupPeriod int    `envconfig:"SHUTDOWN_CLEANUP_PERIOD_SECONDS"`
		ShutdownGracePeriod   int    `envconfig:"SHUTDOWN_GRACE_PERIOD_SECONDS"`
	} `envconfig:"SERVER"`
}

var (
	conf        Config
	once        sync.Once
	initialized bool
)

// Init initializes the configuration system
func Init() error {
	var err error
	once.Do(func() {
		// Load .env file if provided
		err = godotenv.Load(".env")
		if err != nil {
			log.Warn().Err(err).Msg("Could not load .env file, continuing with existing environment variables")
		} else {
			log.Info().Msg("Successfully loaded variables from .env file into environment")
		}

		// Process environment variables into the config struct
		err = envconfig.Process("", &conf)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to process environment variables")
		}

		initialized = true
		log.Info().Msg("Service configuration initialized successfully")
	})

	return err
}

// Get returns the configuration
func Get() *Config {
	// Ensure configuration is initialized
	if !initialized {
		if err := Init(); err != nil {
			log.Fatal().Err(err).Msg("Failed to initialize configuration")
		}
	}
	return &conf
}

// Debug prints out the current configuration
func (c *Config) Debug() {
	fmt.Println("=== Configuration Debug ===")
	fmt.Printf("App Name: %s\n", c.App.Name)
	fmt.Printf("App URL: %s\n", c.App.URL)
	fmt.Printf("App Host: %s\n", c.App.Host)
	fmt.Printf("App Base Path: %s\n", c.App.BasePath)
	fmt.Printf("App Revision: %s\n", c.App.Revision)

	fmt.Println("\nCORS Configuration:")
	fmt.Printf("  Enabled: %v\n", c.App.CORS.Enable)
	fmt.Printf("  Allow Credentials: %v\n", c.App.CORS.AllowCredentials)
	fmt.Printf("  Allowed Headers: %v\n", c.App.CORS.AllowedHeaders)
	fmt.Printf("  Allowed Methods: %v\n", c.App.CORS.AllowedMethods)
	fmt.Printf("  Allowed Origins: %v\n", c.App.CORS.AllowedOrigins)
	fmt.Printf("  Max Age Seconds: %d\n", c.App.CORS.MaxAgeSeconds)

	fmt.Println("\nDatabase Write Configuration:")
	fmt.Printf("  Host: %s\n", c.DB.MySQL.Write.Host)
	fmt.Printf("  Port: %s\n", c.DB.MySQL.Write.Port)
	fmt.Printf("  Name: %s\n", c.DB.MySQL.Write.Name)
	fmt.Printf("  Username: %s\n", c.DB.MySQL.Write.Username)
	fmt.Printf("  Timezone: %s\n", c.DB.MySQL.Write.Timezone)

	fmt.Println("\nServer Configuration:")
	fmt.Printf("  Environment: %s\n", c.Server.Env)
	fmt.Printf("  Log Level: %s\n", c.Server.LogLevel)
	fmt.Printf("  Port: %s\n", c.Server.Port)
	fmt.Printf("  Host: %s\n", c.Server.Host)
}
