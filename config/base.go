package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/atharvbhadange/go-api-template/constants"
)

type Config struct {
	Port        string
	Environment string
	ServiceName string
	Version     string

	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresDB       string
	PostgresPassword string

	PostgresSSLMode      string
	PostgresRootCertLoc  string
	PostgresMaxOpenConns int
	PostgresMaxIdleConns int
	PostgresMaxIdleTime  time.Duration
}

type confVars struct {
	missing   []string //name of the mandatory environment variable that are missing
	malformed []string //errors describing malformed environment varibale values
}

var Conf *Config

func New() (*Config, error) {
	vars := &confVars{}

	port := vars.mandatoryInt("PORT")
	environment := vars.mandatory("ENVIRONMENT")
	serviceName := vars.optional("SERVICE_NAME", "go-service")
	version := vars.optional("VERSION", "1.0.0")

	postgresHost := vars.mandatory("POSTGRES_HOST")
	postgresPort := vars.mandatory("POSTGRES_PORT")
	postgresUser := vars.mandatory("POSTGRES_USER")
	postgresDB := vars.mandatory("POSTGRES_DB")
	postgresPassword := vars.mandatory("POSTGRES_PASSWORD")

	postgresSSLMode := vars.optional("POSTGRES_SSL_MODE", "disable")
	postgresRootCertLoc := vars.optional("POSTGRES_ROOT_CERT_LOC", "")

	postgresMaxOpenConns := vars.optionalInt("POSTGRES_MAX_OPEN_CONNS", constants.POSTGRES_MAX_OPEN_CONNS)
	postgresMaxIdleConns := vars.optionalInt("POSTGRES_MAX_IDLE_CONNS", constants.POSTGRES_MAX_IDLE_CONNS)
	postgresMaxIdleTime := vars.optionalDuration("POSTGRES_MAX_IDLE_TIME", 5*time.Minute)

	if err := vars.Error(); err != nil {
		return nil, fmt.Errorf("error loading configuration: %w", err)
	}

	config := &Config{
		Port:        fmt.Sprintf(":%d", port),
		Environment: environment,
		ServiceName: serviceName,
		Version:     version,

		PostgresHost:     postgresHost,
		PostgresPort:     postgresPort,
		PostgresUser:     postgresUser,
		PostgresDB:       postgresDB,
		PostgresPassword: postgresPassword,

		PostgresSSLMode:      postgresSSLMode,
		PostgresRootCertLoc:  postgresRootCertLoc,
		PostgresMaxOpenConns: postgresMaxOpenConns,
		PostgresMaxIdleConns: postgresMaxIdleConns,
		PostgresMaxIdleTime:  postgresMaxIdleTime,
	}

	Conf = config

	return config, nil
}

func (vars *confVars) optional(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func (vars *confVars) optionalInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	valueInt, err := strconv.Atoi(value)

	if err != nil {
		vars.malformed = append(vars.malformed, key)
		return fallback
	}

	return valueInt
}

func (vars *confVars) optionalBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	valueBool, err := strconv.ParseBool(value)

	if err != nil {
		vars.malformed = append(vars.malformed, key)
		return fallback
	}

	return valueBool
}

func (vars *confVars) optionalDuration(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	valueDuration, err := time.ParseDuration(value)

	if err != nil {
		vars.malformed = append(vars.malformed, key)
		return fallback
	}

	return valueDuration
}

func (vars *confVars) mandatory(key string) string {
	value := os.Getenv(key)
	if value == "" {
		vars.missing = append(vars.missing, key)
	}
	return value
}

func (vars *confVars) mandatoryInt(key string) int {
	value := os.Getenv(key)
	if value == "" {
		vars.missing = append(vars.missing, key)
		return 0
	}

	valueInt, err := strconv.Atoi(value)

	if err != nil {
		vars.malformed = append(vars.malformed, key)
		return 0
	}

	return valueInt
}

func (vars *confVars) mandatoryDuration(key string) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		vars.missing = append(vars.missing, key)
		return 0
	}

	valueDuration, err := time.ParseDuration(value)

	if err != nil {
		vars.malformed = append(vars.malformed, key)
		return 0
	}

	return valueDuration
}

func (vars *confVars) mandatoryBool(key string) bool {
	value := os.Getenv(key)
	if value == "" {
		vars.missing = append(vars.missing, key)
		return false
	}

	valueBool, err := strconv.ParseBool(value)

	if err != nil {
		vars.malformed = append(vars.malformed, key)
		return false
	}

	return valueBool
}

func (vars confVars) Error() error {
	if len(vars.missing) > 0 {
		return fmt.Errorf("missing mandatory configurations: %s", strings.Join(vars.missing, ", "))
	}

	if len(vars.malformed) > 0 {
		return fmt.Errorf("malformed configurations: %s", strings.Join(vars.malformed, "; "))
	}
	return nil
}
