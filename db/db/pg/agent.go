package pg

import (
	"fmt"

	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/lainio/err2"
)

const (
	sqlAgentFields = "id, agent_id, label, created, last_accessed"
	sqlAgentInsert = "INSERT INTO agent (agent_id, label) VALUES ($1, $2) " +
		"ON CONFLICT (agent_id) DO UPDATE SET last_accessed = now() RETURNING " + sqlAgentFields
	sqlAgentSelect          = "SELECT " + sqlAgentFields + " FROM agent"
	sqlAgentSelectByID      = sqlAgentSelect + " WHERE id=$1"
	sqlAgentSelectByAgentID = sqlAgentSelect + " WHERE agent_id=$1"
)

func (p *Database) AddAgent(a *model.Agent) (n *model.Agent, err error) {
	defer returnErr("AddAgent", &err)

	rows, err := p.db.Query(
		sqlAgentInsert,
		a.AgentID,
		a.Label,
	)
	err2.Check(err)
	defer rows.Close()

	n = a.Copy()
	if rows.Next() {
		err = rows.Scan(&n.ID, &n.AgentID, &n.Label, &n.Created, &n.LastAccessed)
		err2.Check(err)
	}

	err = rows.Err()
	err2.Check(err)

	return
}

func (p *Database) GetAgent(id, agentID *string) (a *model.Agent, err error) {
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

	a = model.NewAgent()
	if rows.Next() {
		err = rows.Scan(&a.ID, &a.AgentID, &a.Label, &a.Created, &a.LastAccessed)
		err2.Check(err)
	}

	err = rows.Err()
	err2.Check(err)

	return
}
