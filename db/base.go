package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	C "github.com/atharvbhadange/go-api-template/config"
)

var PostgresConn *sql.DB

func GetPostgresURL() string {
	dbHost := C.Conf.PostgresHost
	dbPort := C.Conf.PostgresPort
	dbUser := C.Conf.PostgresUser
	dbPass := C.Conf.PostgresPassword
	dbName := C.Conf.PostgresDB

	if C.Conf.PostgresSSLMode == "disable" {
		return fmt.Sprintf("host=%s port=%s user=%s "+
			"password=%s dbname=%s sslmode=disable",
			dbHost, dbPort, dbUser, dbPass, dbName) 
	} else {
		return fmt.Sprintf("host=%s port=%s user=%s "+
			"password=%s dbname=%s sslmode=%s sslrootcert=%s",
			dbHost, dbPort, dbUser, dbPass, dbName, C.Conf.PostgresSSLMode, C.Conf.PostgresRootCertLoc)
	}
}

func Init() error {
	var err error
	PostgresConn, err = sql.Open("postgres", GetPostgresURL())
	if err != nil {
		return fmt.Errorf("error opening database connection: %w", err)
	}

	err = PostgresConn.Ping()
	if err != nil {
		return fmt.Errorf("error pinging database: %w", err)
	}

	PostgresConn.SetMaxOpenConns(C.Conf.PostgresMaxOpenConns)
	PostgresConn.SetMaxIdleConns(C.Conf.PostgresMaxIdleConns)

	return nil
}

func PGTransaction(ctx context.Context) (*sql.Tx, error) {
	tx, err := PostgresConn.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func Close() {
	PostgresConn.Close()
}
