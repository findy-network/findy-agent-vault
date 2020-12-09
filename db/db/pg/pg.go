package pg

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/findy-network/findy-agent-vault/db/db"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // blank for migrate driver

	"github.com/golang/glog"
	"github.com/lainio/err2"
	"github.com/lib/pq"
)

const (
	host   = "localhost"
	user   = "postgres"
	dbName = "vault"
)

type PostgresErrorCode string

const (
	PostgresErrorUniqueViolation PostgresErrorCode = "unique_violation"
)

type PostgresError struct {
	operation string
	code      PostgresErrorCode
	error     *pq.Error
}

func returnErr(op string, err *error) {
	if r := recover(); r != nil {
		e, ok := r.(error)
		if !ok {
			panic(r)
		}
		*err = e
	}

	if pgErr, ok := (*err).(*pq.Error); ok {
		*err = &PostgresError{operation: op, code: PostgresErrorCode(pgErr.Code.Name()), error: pgErr}
	}
}

func (e *PostgresError) Error() string {
	return e.error.Error()
}

type Database struct {
	db *sql.DB
}

func InitDB(migratePath, port string, reset bool) db.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, os.Getenv("POSTGRES_PASSWORD"), dbName)
	sqlDB, err := sql.Open("postgres", psqlInfo)
	err2.Check(err)

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	err2.Check(err)

	m, err := migrate.NewWithDatabaseInstance(
		migratePath, "postgres", driver,
	)
	err2.Check(err)

	if reset {
		err = m.Down()
		if err == migrate.ErrNoChange {
			glog.Info("empty db, skipping db cleanup")
		} else {
			err2.Check(err)
			version, _, _ := m.Version()
			glog.Infof("db reset to version: %d", version)
		}
	}

	err = m.Up()
	if err == migrate.ErrNoChange {
		glog.Info("no new migrations, skipping db modifications")
	} else {
		err2.Check(err)
		version, _, _ := m.Version()
		glog.Infof("migrations ok: %d", version)
	}

	err = sqlDB.Ping()
	err2.Check(err)

	glog.Infof("successfully connected to postgres %s:%s\n", host, port)
	return &Database{db: sqlDB}
}

func (p *Database) Close() {
	p.db.Close()
}
