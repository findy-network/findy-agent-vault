package pg

import (
	"fmt"

	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/lainio/err2"
)

const (
	sqlAgentFields = "id, agent_id, label, raw_jwt, created, last_accessed"
	sqlAgentInsert = "INSERT INTO agent (agent_id, label, raw_jwt) VALUES ($1, $2, $3) " +
		"ON CONFLICT (agent_id) DO UPDATE SET last_accessed = (now() at time zone 'UTC') RETURNING " + sqlAgentFields
	sqlAgentSelect          = "SELECT " + sqlAgentFields + " FROM agent"
	sqlAgentSelectByID      = sqlAgentSelect + " WHERE id=$1"
	sqlAgentSelectByAgentID = sqlAgentSelect + " WHERE agent_id=$1"
)

func (pg *Database) AddAgent(a *model.Agent) (n *model.Agent, err error) {
	defer returnErr("AddAgent", &err)

	rows, err := pg.db.Query(
		sqlAgentInsert,
		a.AgentID,
		a.Label,
		a.RawJWT,
	)
	err2.Check(err)
	defer rows.Close()

	n = model.NewAgent(a)
	if rows.Next() {
		err = rows.Scan(&n.ID, &n.AgentID, &n.Label, &n.RawJWT, &n.Created, &n.LastAccessed)
		err2.Check(err)
	}

	err = rows.Err()
	err2.Check(err)

	n.TenantID = n.ID

	return
}

func (pg *Database) GetAgent(id, agentID *string) (a *model.Agent, err error) {
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

	rows, err := pg.db.Query(query, *queryID)
	err2.Check(err)
	defer rows.Close()

	a = model.NewAgent(nil)
	if rows.Next() {
		err = rows.Scan(&a.ID, &a.AgentID, &a.Label, &a.RawJWT, &a.Created, &a.LastAccessed)
		err2.Check(err)
	}

	err = rows.Err()
	err2.Check(err)

	a.TenantID = a.ID

	return
}
