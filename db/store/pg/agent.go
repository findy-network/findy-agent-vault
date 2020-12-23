package pg

import (
	"database/sql"
	"fmt"
	"sort"

	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/paginator"
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

func (pg *Database) GetListenerAgents(info *paginator.BatchInfo) (a *model.Agents, err error) {
	defer returnErr("GetListenerAgents", &err)

	jwtNotNull := " raw_jwt IS NOT NULL "
	query, args := getBatchQuery(&queryInfo{
		Asc:        sqlAgentSelect + " WHERE " + jwtNotNull + sqlOrderByAsc("") + " $1",
		Desc:       sqlAgentSelect + " WHERE " + jwtNotNull + sqlOrderByDesc("") + " $1",
		AfterAsc:   sqlAgentSelect + " WHERE cursor > $1 AND" + jwtNotNull + sqlOrderByAsc("") + " $2",
		AfterDesc:  sqlAgentSelect + " WHERE cursor > $1 AND" + jwtNotNull + sqlOrderByDesc("") + " $2",
		BeforeAsc:  sqlAgentSelect + " WHERE cursor < $1 AND" + jwtNotNull + sqlOrderByAsc("") + " $2",
		BeforeDesc: sqlAgentSelect + " WHERE cursor < $1 AND" + jwtNotNull + sqlOrderByDesc("") + " $2",
	},
		info,
		"",
		[]interface{}{},
	)

	rows, err := pg.db.Query(query, args...)
	err2.Check(err)
	defer rows.Close()

	a = &model.Agents{
		Agents:          make([]*model.Agent, 0),
		HasNextPage:     false,
		HasPreviousPage: false,
	}
	var agent *model.Agent
	for rows.Next() {
		agent, err = readRowToAgent(rows)
		err2.Check(err)
		a.Agents = append(a.Agents, agent)
	}

	err = rows.Err()
	err2.Check(err)

	if info.Count < len(a.Agents) {
		a.Agents = a.Agents[:info.Count]
		if info.Tail {
			a.HasPreviousPage = true
		} else {
			a.HasNextPage = true
		}
	}

	if info.After > 0 {
		a.HasPreviousPage = true
	}
	if info.Before > 0 {
		a.HasNextPage = true
	}

	// Reverse order for tail first
	if info.Tail {
		sort.Slice(a.Agents, func(i, j int) bool {
			return a.Agents[i].Created.Sub(a.Agents[j].Created) < 0
		})
	}

	return a, err
}

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

func readRowToAgent(rows *sql.Rows) (a *model.Agent, err error) {
	a = model.NewAgent(nil)
	err = rows.Scan(
		&a.ID, &a.AgentID, &a.Label, &a.RawJWT, &a.Created, &a.LastAccessed,
	)
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

	if rows.Next() {
		a, err = readRowToAgent(rows)
		err2.Check(err)
	}

	err = rows.Err()
	err2.Check(err)

	a.TenantID = a.ID

	return
}
