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
	"github.com/lainio/err2/try"
)

const (
	user = "postgres"
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

func createNewDB(config *utils.Configuration) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s sslmode=disable",
		config.DBHost, config.DBPort, user, config.DBPassword)
	sqlDB, _ := sql.Open("postgres", psqlInfo)
	if _, err := sqlDB.Exec("CREATE database " + config.DBName); err != nil {
		glog.Warning(err)
	}
	sqlDB.Close()
}

func InitDB(config *utils.Configuration, reset, createDB bool) store.DB {
	// For testing only
	if createDB {
		createNewDB(config)
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.DBHost, config.DBPort, user, config.DBPassword, config.DBName)

	var sqlDB *sql.DB
	var err error
	if config.DBTracing {
		sqlDB = try.To1(initTraceHook(psqlInfo))
	} else {
		sqlDB = try.To1(sql.Open("postgres", psqlInfo))
	}

	driver := try.To1(postgres.WithInstance(sqlDB, &postgres.Config{}))

	m := try.To1(migrate.NewWithDatabaseInstance(
		config.DBMigrationsPath, "postgres", driver,
	))

	if reset {
		err = m.Down()
		if err == migrate.ErrNoChange {
			glog.Info("empty db, skipping db cleanup")
		} else {
			try.To(err)
			version, _, _ := m.Version()
			glog.Infof("db reset to version: %d", version)
		}
	}

	err = m.Up()
	if err == migrate.ErrNoChange {
		glog.Info("no new migrations, skipping db modifications")
	} else {
		try.To(err)
		version, _, _ := m.Version()
		glog.Infof("migrations ok: %d", version)
	}

	try.To(sqlDB.Ping())

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
	defer err2.Handle(&err)

	query := "SELECT count(id) FROM " + tableName
	args := make([]interface{}, 0)
	args = append(args, tenantID)
	if connectionID == nil {
		query += batchWhere
	} else {
		query += batchWhereConnection
		args = append(args, *connectionID)
	}

	try.To(pg.doRowQuery(
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

func (pg *Database) doRowQuery(scan func(*sql.Rows) error, query string, args ...interface{}) (err error) {
	defer err2.Handle(&err)

	rows := try.To1(pg.db.Query(query, args...))
	defer rows.Close()

	if rows.Next() {
		err = scan(rows)
	} else {
		err = store.NewError(store.ErrCodeNotFound, "no rows returned")
	}
	try.To(err)
	try.To(rows.Err())

	return nil
}

func (pg *Database) doRowsQuery(scan func(*sql.Rows) error, query string, args ...interface{}) (err error) {
	defer err2.Handle(&err)

	rows := try.To1(pg.db.Query(query, args...))
	defer rows.Close()

	scanCount := 0
	for rows.Next() {
		err = scan(rows)
		scanCount++
	}

	if scanCount == 0 {
		err = store.NewError(store.ErrCodeNotFound, "no rows returned")
	}
	try.To(err)
	try.To(rows.Err())

	return nil
}
