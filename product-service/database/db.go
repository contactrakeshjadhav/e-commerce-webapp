package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/contactrakeshjadhav/e-commerce-webapp/product-service/commons/database"
	"github.com/contactrakeshjadhav/e-commerce-webapp/product-service/commons/pkg/config"
	"github.com/golang-migrate/migrate/v4"
	databasePostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func init() {
	godotenv.Load()
}

type DB struct {
	Conn   *sql.DB
	Config config.DBConfig
}

func InitDB(config config.DBConfig) (*DB, error) {
	sslMode := database.GetSslMode(config)
	dsn := fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v sslmode=%v",
		config.DbHost,
		config.DbPort,
		config.DbUser,
		config.DbName,
		config.DbPass,
		sslMode,
	)

	// we create the DB is does not exists
	if err := database.SetupDatabase(config); err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, errors.New("failed to connect Database")
	}

	return &DB{
		Conn:   db,
		Config: config,
	}, nil
}

func (db *DB) Automigrate() error {
	if db.Config.AutoMigrate {
		return migrateTables(db.Conn, db.Config.MigrationsLocation)
	}
	return nil
}

func InitDBTest(port int, migrationsPath string, dbName string, migrate bool) *sql.DB {
	dsn := fmt.Sprintf("host=localhost port=%d user=postgres dbname=postgres password=%s sslmode=disable", port, getTestPassword())
	defaultDB, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect Test Database: " + err.Error())
	}
	defer defaultDB.Close()

	stmt := fmt.Sprintf("CREATE DATABASE %s", dbName)
	_, err = defaultDB.Exec(stmt)
	if err != nil {
		log.Fatal(err)
	}

	dsn = fmt.Sprintf("host=localhost port=%d user=postgres dbname=%v password=%s sslmode=disable", port, dbName, getTestPassword())
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	if migrate {
		if err = migrateTables(db, fmt.Sprintf("file://%s", migrationsPath)); err != nil {
			log.Fatalf("failed to automigrate Test Database: " + err.Error())
		}
	}

	return db
}

func getTestPassword() string {
	return "postgres"
}

func migrateTables(db *sql.DB, migrationsPath string) error {
	driver, err := databasePostgres.WithInstance(db, &databasePostgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err == migrate.ErrNoChange {
		// We are up to date
		return nil
	}

	return err
}
