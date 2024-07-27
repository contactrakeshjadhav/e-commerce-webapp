package config

import (
	"errors"
	"os"
	"strconv"
	"time"
)

const (
	DB_AUTO_MIGRATE            string = "DB_AUTO_MIGRATE"
	DB_PORT                    string = "DB_PORT"
	DB_NAME                    string = "DB_NAME"
	DB_USER                    string = "DB_USER"
	DB_PASS                    string = "DB_PASS"
	DB_HOST                    string = "DB_HOST"
	DB_MAX_OPEN_CONNECTIONS    string = "DB_MAX_OPEN_CONNECTIONS"
	DB_MAX_IDLE_CONNECTIONS    string = "DB_MAX_IDLE_CONNECTIONS"
	DB_OPEN_CONNECTION_TIMEOUT string = "DB_OPEN_CONNECTION_TIMEOUT"
	DB_IDLE_CONNECTION_TIMEOUT string = "DB_IDLE_CONNECTION_TIMEOUT"
	DB_MIGRATIONS_LOCATION     string = "DB_MIGRATIONS_LOCATION"
	DB_DISABLE_SSL             string = "DB_DISABLE_SSL"

	// default values
	defaultMaxOpenConnections    int = 20
	defaultMaxIdleConnections    int = 5
	defaultOpenConnectionTimeout int = 120
	defaultIdleConnectionTimeout int = 120
)

type DBConfig struct {
	DbHost                  string
	DbPort                  int
	DbUser                  string
	DbName                  string
	DbPass                  string
	DbMaxOpenConnections    int
	DbMaxIdleConnections    int
	DbOpenConnectionTimeout time.Duration
	DbIdleConnectionTimeout time.Duration
	AutoMigrate             bool
	MigrationsLocation      string
	DisableSsl              bool
}

// GetDatabaseVariables extract and return the postgres connection vars of the OS environment
// return an error if one of this vars doesn't exists
func GetDatabaseVariables() (DBConfig, error) {
	host, present := os.LookupEnv(DB_HOST)
	if !present {
		return DBConfig{}, errors.New("postgres database host required")
	}

	port, present := os.LookupEnv(DB_PORT)
	if !present {
		return DBConfig{}, errors.New("postgres database port required")
	}
	p, _ := strconv.Atoi(port)

	user, present := os.LookupEnv(DB_USER)
	if !present {
		return DBConfig{}, errors.New("postgres database user required")
	}

	dbName, present := os.LookupEnv(DB_NAME)
	if !present {
		return DBConfig{}, errors.New("postgres database dbName required")
	}

	dbPass, present := os.LookupEnv(DB_PASS)
	if !present {
		return DBConfig{}, errors.New("postgres database password required")
	}

	autoMigrate, present := os.LookupEnv(DB_AUTO_MIGRATE)
	if !present {
		return DBConfig{}, errors.New("postgres database auto_migrate required")
	}
	migrate, err := strconv.ParseBool(autoMigrate)
	if err != nil {
		return DBConfig{}, errors.New("could not parse postgres auto_migrate")
	}

	migrationsLocation, present := os.LookupEnv(DB_MIGRATIONS_LOCATION)
	if migrate && !present {
		return DBConfig{}, errors.New("postgres migration location required when auto_migrate")
	}

	disableSslEnv, present := os.LookupEnv(DB_DISABLE_SSL)
	disableSsl := false
	if present {
		disableSsl, _ = strconv.ParseBool(disableSslEnv)
	}

	maxOpenConn := defaultMaxOpenConnections
	openConnValue, present := os.LookupEnv(DB_MAX_OPEN_CONNECTIONS)
	if present {
		value, err := strconv.Atoi(openConnValue)
		if err != nil {
			return DBConfig{}, errors.New("could not parse postgres max_open_connections value")
		}
		maxOpenConn = value
	}

	maxIdleConn := defaultMaxIdleConnections
	idleConnValue, present := os.LookupEnv(DB_MAX_IDLE_CONNECTIONS)
	if present {
		value, err := strconv.Atoi(idleConnValue)
		if err != nil {
			return DBConfig{}, errors.New("could not parse postgres max_idle_connections value")
		}
		maxIdleConn = value
	}

	openConnTimeout := defaultOpenConnectionTimeout
	openConnTimeoutValue, present := os.LookupEnv(DB_OPEN_CONNECTION_TIMEOUT)
	if present {
		value, err := strconv.Atoi(openConnTimeoutValue)
		if err != nil {
			return DBConfig{}, errors.New("could not parse postgres open_connection_timeout value")
		}
		openConnTimeout = value
	}

	idleConnTimeout := defaultIdleConnectionTimeout
	idleConnTimeoutValue, present := os.LookupEnv(DB_IDLE_CONNECTION_TIMEOUT)
	if present {
		value, err := strconv.Atoi(idleConnTimeoutValue)
		if err != nil {
			return DBConfig{}, errors.New("could not parse postgres open_connection_timeout value")
		}
		openConnTimeout = value
	}

	return DBConfig{
		DbHost:                  host,
		DbPort:                  p,
		DbUser:                  user,
		DbName:                  dbName,
		DbPass:                  dbPass,
		DbMaxOpenConnections:    maxOpenConn,
		DbMaxIdleConnections:    maxIdleConn,
		DbOpenConnectionTimeout: time.Duration(openConnTimeout) * time.Second,
		DbIdleConnectionTimeout: time.Duration(idleConnTimeout) * time.Second,
		AutoMigrate:             migrate,
		MigrationsLocation:      migrationsLocation,
		DisableSsl:              disableSsl,
	}, nil
}
