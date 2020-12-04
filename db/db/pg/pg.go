package pg

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/findy-network/findy-agent-vault/db/db"
	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/golang/glog"
	"github.com/lainio/err2"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

const (
	sqlAgentInsert          = "INSERT INTO agent (agent_id, label) VALUES ($1, $2)"
	sqlAgentSelect          = "SELECT id, agent_id, label, created FROM agent"
	sqlAgentSelectByID      = sqlAgentSelect + " WHERE id=$1"
	sqlAgentSelectByAgentID = sqlAgentSelect + " WHERE agent_id=$1"
	/*insertToConnection = `INSERT INTO connection (
		tenant_id, our_did, their_did, their_endpoint, their_label, invited
	) VALUES ($1, $2, $3, $4, $5, $6)`*/
)

const (
	host   = "localhost"
	user   = "postgres"
	dbName = "vault"
)

type PgErrorCode string

const (
	PgErrorUniqueViolation PgErrorCode = "unique_violation"
)

type PgError struct {
	operation string
	code      PgErrorCode
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
		*err = &PgError{operation: op, code: PgErrorCode(pgErr.Code.Name()), error: pgErr}
	}
}

func (e *PgError) Error() string {
	return e.error.Error()
}

type PgDb struct {
	db *sql.DB
}

func InitDb(migratePath string, port string, reset bool) db.Db {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, os.Getenv("POSTGRES_PASSWORD"), dbName)
	sqlDB, err := sql.Open("postgres", psqlInfo)
	err2.Check(err)

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
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
	return &PgDb{db: sqlDB}
}

func (p *PgDb) Close() {
	p.db.Close()
}

func (p *PgDb) AddAgent(agentId, label string) (err error) {
	defer returnErr("AddAgent", &err)

	_, err = p.db.Exec(sqlAgentInsert, agentId, label)
	err2.Check(err)

	return
}

func (p *PgDb) GetAgent(id, agentID *string) (a *model.Agent, err error) {
	defer returnErr("GetAgent", &err)

	if id == nil && agentID == nil {
		panic(fmt.Errorf("either id or agent id is required"))
	}
	query := sqlAgentSelectByID
	queryID := id
	if id == nil {
		query = sqlAgentSelectByAgentID
		queryID = agentID
	}

	rows, err := p.db.Query(query, *queryID)
	err2.Check(err)
	defer rows.Close()

	a = &model.Agent{}
	if rows.Next() {
		err = rows.Scan(&a.ID, &a.AgentID, &a.Label, &a.Created)
		err2.Check(err)
	}

	err = rows.Err()
	err2.Check(err)

	return
}
