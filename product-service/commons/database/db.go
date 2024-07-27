package database

import (
	"database/sql"
	"fmt"

	"github.com/contactrakeshjadhav/e-commerce-webapp/product-service/commons/pkg/config"
)

type DB struct {
	Conn   *sql.DB
	Config config.DBConfig
}

func InitDB(config config.DBConfig) (*DB, error) {
	sslMode := GetSslMode(config)
	dsn := fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v sslmode=%v",
		config.DbHost,
		config.DbPort,
		config.DbUser,
		config.DbName,
		config.DbPass,
		sslMode,
	)

	// we create the DB is does not exists
	if err := SetupDatabase(config); err != nil {
		return nil, err
	}

	//here is where the actual connection happens
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	return &DB{
		Conn:   db,
		Config: config,
	}, nil
}

func SetupDatabase(config config.DBConfig) error {
	// we create a secondary SQL connection
	// it will connect to the default DB
	sslMode := GetSslMode(config)
	var rows int64
	dsn := fmt.Sprintf("host=%s port=%v user=%s password=%s sslmode=%s",
		config.DbHost,
		config.DbPort,
		config.DbUser,
		config.DbPass,
		sslMode,
	)
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	defer conn.Close()

	// check if db exists
	result := conn.QueryRow("SELECT COUNT(*) FROM pg_database WHERE datname=$1;", config.DbName)
	if result.Err() != nil {
		return result.Err()
	}

	err = result.Scan(&rows)
	if err != nil {
		return err
	}

	// the DB exists
	if rows > 0 {
		return nil
	}

	//postgres bind parameters cannot be used  here because it can only bind to a literal.
	//I have wrapped the variable in double quotes to prevent injection
	stmt := fmt.Sprintf("CREATE DATABASE \"%s\";", config.DbName)
	_, err = conn.Exec(stmt)
	return err
}

func GetSslMode(config config.DBConfig) string {
	sslMode := "require"
	if config.DisableSsl {
		sslMode = "disable"
	}
	return sslMode
}
