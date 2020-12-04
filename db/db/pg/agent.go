package pg

import (
	"fmt"

	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/lainio/err2"
)

const (
	sqlAgentInsert          = "INSERT INTO agent (agent_id, label) VALUES ($1, $2)"
	sqlAgentSelect          = "SELECT id, agent_id, label, created FROM agent"
	sqlAgentSelectByID      = sqlAgentSelect + " WHERE id=$1"
	sqlAgentSelectByAgentID = sqlAgentSelect + " WHERE agent_id=$1"
)

func (p *Database) AddAgent(a *model.Agent) (err error) {
	defer returnErr("AddAgent", &err)

	_, err = p.db.Exec(sqlAgentInsert, a.AgentID, a.Label)
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

	fmt.Println("Haaa1")

	rows, err := p.db.Query(query, *queryID)
	err2.Check(err)
	defer rows.Close()

	fmt.Println("Haaa2", query, *queryID)

	a = &model.Agent{}
	if rows.Next() {
		fmt.Println("Haaa4")
		err = rows.Scan(&a.ID, &a.AgentID, &a.Label, &a.Created)
		err2.Check(err)
		fmt.Println("Haaa5")
	}

	fmt.Println("Haaa3")

	err = rows.Err()
	err2.Check(err)

	return
}
