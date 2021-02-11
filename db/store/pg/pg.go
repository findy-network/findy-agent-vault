package pg

import (
	"database/sql"
	"fmt"

	"github.com/findy-network/findy-agent-vault/db/store"
	"github.com/findy-network/findy-agent-vault/paginator"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // blank for migrate driver

	"github.com/golang/glog"
	"github.com/lainio/err2"
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

	sqlOrderByCursorAsc  = " ORDER BY cursor ASC LIMIT"
	sqlOrderByCursorDesc = " ORDER BY cursor DESC LIMIT"
)

var (
	sqlInsertFields = sqlFields("", []string{"id", "created", "cursor"})
)

type pgError struct {
	error
	code store.ErrCode
}

func (e *pgError) Code() store.ErrCode {
	return e.code
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

type Database struct {
	db *sql.DB
}

func InitDB(migratePath string, config *utils.Configuration, reset bool) store.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.DBHost, config.DBPort, user, config.DBPassword, dbName)

	var sqlDB *sql.DB
	var err error
	if config.DBTracing {
		sqlDB, err = initTraceHook(psqlInfo)
	} else {
		sqlDB, err = sql.Open("postgres", psqlInfo)
	}
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

	glog.Infof("successfully connected to postgres %s:%d\n", config.DBHost, config.DBPort)
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

	err2.Check(pg.doQuery(
		func(rows *sql.Rows) error {
			return rows.Scan(&count)
		},
		query,
		args...,
	))

	return
}

func sqlFields(tableName string, fields []string) string {
	if tableName != "" {
		tableName += "."
	}
	q := ""
	for i, field := range fields {
		if i != 0 {
			q += ","
		}
		q += tableName + field
	}
	return q
}

func sqlArguments(fields []string) string {
	q := ""
	for i := range fields {
		if i != 0 {
			q += ","
		}
		q += fmt.Sprintf("$%d", i+1)
	}
	return q
}

func (pg *Database) doQuery(scan func(*sql.Rows) error, query string, args ...interface{}) (err error) {
	defer err2.Return(&err)

	rows, err := pg.db.Query(query, args...)
	err2.Check(err)
	defer rows.Close()

	if rows.Next() {
		err = scan(rows)
	} else {
		err = pgError{
			fmt.Errorf("no rows returned"), store.NotExists,
		}
	}
	err2.Check(err)
	err2.Check(rows.Err())

	return nil
}
