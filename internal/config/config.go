package config

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	EnvDevelopment = "development"
	EnvProduction  = "production"
	EnvConfigPath  = "configs/development.json"
)

func mustGetEnv(env string) string {
	value, ok := os.LookupEnv(env)
	if !ok {
		panic(fmt.Errorf("environment variable %s must be set", env))
	}

	return value
}

func mustGetInt(s string) int {
	value, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Errorf("%s could not be converted to int: %v", s, err))
	}
	return value
}

type Config struct {
	Server   *Server
	Logger   *Logger
	Database *Database
	Redis    *Redis
	Mqtt     *Mqtt
	Jwt      *Jwt
}

type Server struct {
	Port        int    `json:"port,omitempty"`
	Environment string `json:"environment,omitempty"`
}

func NewServerConfigFromEnv() *Server {
	return &Server{
		Port:        mustGetInt(mustGetEnv("PORT")),
		Environment: mustGetEnv("ENVIRONMENT"),
	}
}

type Logger struct {
	Level     string `json:"level,omitempty"`
	AddSource bool   `json:"addSource,omitempty"`
}

func NewLoggerConfig() *Logger {
	level, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		level = "info"
	}

	addSourceStr, ok := os.LookupEnv("LOG_ADD_SOURCE")
	if !ok {
		addSourceStr = "true"
	}

	addSource, err := strconv.ParseBool(addSourceStr)
	if err != nil {
		panic(fmt.Errorf("error parsing LOG_ADD_SOURCE: %v", err))
	}

	return &Logger{
		Level:     level,
		AddSource: addSource,
	}
}

type Database struct {
	Uri string `json:"uri"`
}

func NewDatabaseConfigFromEnv() *Database {
	return &Database{
		Uri: mustGetEnv("DATABASE_URI"),
	}
}

type Redis struct {
	ConnectionString string        `json:"connectionString,omitempty"`
	Expire           time.Duration `json:"expire,omitempty"`
}

func NewRedisConfigFromEnv() *Redis {
	expire := mustGetInt(mustGetEnv("REDIS_EXPIRE_MINUTES"))
	return &Redis{
		ConnectionString: mustGetEnv("REDIS_CONNECTION_STRING"),
		Expire:           time.Duration(expire) * time.Minute,
	}
}

type Mqtt struct {
	Host     string `json:"host,omitempty"`
	Username string `json:"username,omitempty"`
	Port     string `json:"port,omitempty"`
}

func NewMqttConfigFromEnv() *Mqtt {
	return &Mqtt{
		Host:     mustGetEnv("MQTT_HOST"),
		Username: mustGetEnv("MQTT_USERNAME"),
		Port:     mustGetEnv("MQTT_PORT"),
	}
}

type Jwt struct {
	Secret   string        `json:"secret,omitempty"`
	Audience string        `json:"audience,omitempty"`
	Issuer   string        `json:"issuer,omitempty"`
	Expire   time.Duration `json:"expire,omitempty"`
}

func NewJwtConfigFromEnv() *Jwt {
	expire := mustGetInt(mustGetEnv("JWT_EXPIRE_MINUTES"))
	return &Jwt{
		Secret:   mustGetEnv("JWT_SECRET"),
		Audience: mustGetEnv("JWT_AUDIENCE"),
		Issuer:   mustGetEnv("JWT_ISSUER"),
		Expire:   time.Duration(expire) * time.Minute,
	}
}

func loadConfigFromFile(filePath string) (*Config, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			slog.Error(fmt.Errorf("error closing file: %v at file path: %s", err, filePath).Error(), slog.String("filePath", filePath))
		}
	}(f)

	var config Config
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to load config from %s: %w", filePath, err)
	}

	if config.Redis != nil {
		config.Redis.Expire = config.Redis.Expire * time.Minute
	}

	return &config, nil
}

func loadConfigFromEnv() *Config {
	config := &Config{
		Server:   NewServerConfigFromEnv(),
		Logger:   NewLoggerConfig(),
		Database: NewDatabaseConfigFromEnv(),
		Redis:    NewRedisConfigFromEnv(),
		Jwt:      NewJwtConfigFromEnv(),
		Mqtt:     NewMqttConfigFromEnv(),
	}

	return config
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load("configs/.env"); err != nil {
		slog.Warn("Error loading .env file, falling back to environment variables")
	}

	env, ok := os.LookupEnv("ENVIRONMENT")
	if !ok {
		env = EnvDevelopment
	}

	slog.Info("Loading configuration", slog.String("env", env))
	if strings.ToLower(env) == EnvProduction {
		return loadConfigFromEnv(), nil
	}

	filePath, ok := os.LookupEnv("CONFIG_PATH")
	if !ok {
		filePath = EnvConfigPath
	}

	slog.Info("Loading development config from file", slog.String("file", filePath))
	return loadConfigFromFile(filePath)
}
