package pg

import (
	"database/sql"
	"fmt"

	"github.com/findy-network/findy-agent-vault/db/store"
	"github.com/findy-network/findy-agent-vault/paginator"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // blank for migrate driver

	"github.com/golang/glog"
	"github.com/lainio/err2"
	"github.com/lib/pq"
)

const (
	user   = "postgres"
	dbName = "vault"
)

const (
	sqlGreaterThan = " > "
	sqlLessThan    = " < "
	sqlAsc         = "ASC"
	sqlDesc        = "DESC"
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

type queryInfo struct {
	Asc        string
	Desc       string
	AfterAsc   string
	AfterDesc  string
	BeforeAsc  string
	BeforeDesc string
}

func getBatchQuery(
	queries *queryInfo,
	batch *paginator.BatchInfo,
	tenantID string,
	initialArgs []interface{},
) (query string, args []interface{}) {
	args = make([]interface{}, 0)
	if tenantID != "" {
		args = append(args, tenantID)
	}

	if batch.Tail {
		query = queries.Desc
		if batch.After > 0 {
			query = queries.AfterDesc
		} else if batch.Before > 0 {
			query = queries.BeforeDesc
		}
	} else {
		query = queries.Asc
		if batch.After > 0 {
			query = queries.AfterAsc
		} else if batch.Before > 0 {
			query = queries.BeforeAsc
		}
	}
	if batch.After > 0 {
		args = append(args, batch.After)
	} else if batch.Before > 0 {
		args = append(args, batch.Before)
	}
	args = append(args, initialArgs...)

	args = append(args, batch.Count+1)
	return query, args
}

func sqlOrderByAsc(orderBy string) string {
	if orderBy != "" {
		orderBy = ", " + orderBy
	}
	return fmt.Sprintf(" ORDER BY cursor ASC %s LIMIT", orderBy)
}

func sqlOrderByDesc(orderBy string) string {
	if orderBy != "" {
		orderBy = ", " + orderBy
	}
	return fmt.Sprintf(" ORDER BY cursor DESC %s LIMIT", orderBy)
}

type Database struct {
	db *sql.DB
}

func InitDB(migratePath, host, password string, port int, reset bool) store.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)
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

	glog.Infof("successfully connected to postgres %s:%d\n", host, port)
	return &Database{db: sqlDB}
}

func (pg *Database) Close() {
	pg.db.Close()
}

func (pg *Database) getCount(
	tableName, batchWhere, batchWhereConnection, tenantID string,
	connectionID *string,
) (count int, err error) {
	defer err2.Return(&err)

	query := "SELECT count(id) FROM " + tableName
	args := make([]interface{}, 0)
	args = append(args, tenantID)
	if connectionID == nil {
		query += batchWhere
	} else {
		query += batchWhereConnection
		args = append(args, *connectionID)
	}

	rows, err := pg.db.Query(query, args...)
	err2.Check(err)
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&count)
		err2.Check(err)
	}

	err = rows.Err()
	err2.Check(err)

	return
}
