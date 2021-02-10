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

var (
	sqlJwtNotSet            = " raw_jwt IS NOT NULL AND raw_jwt != ''"
	sqlAgentJwtNotNullAsc   = sqlJwtNotSet + sqlOrderByCursorAsc
	sqlAgentjJwtNotNullDesc = sqlJwtNotSet + sqlOrderByCursorDesc
	agentQueryInfo          = &queryInfo{
		Asc:        sqlAgentSelect + " WHERE " + sqlAgentJwtNotNullAsc + " $1",
		Desc:       sqlAgentSelect + " WHERE " + sqlAgentjJwtNotNullDesc + " $1",
		AfterAsc:   sqlAgentSelect + " WHERE cursor > $1 AND" + sqlAgentJwtNotNullAsc + " $2",
		AfterDesc:  sqlAgentSelect + " WHERE cursor > $1 AND" + sqlAgentjJwtNotNullDesc + " $2",
		BeforeAsc:  sqlAgentSelect + " WHERE cursor < $1 AND" + sqlAgentJwtNotNullAsc + " $2",
		BeforeDesc: sqlAgentSelect + " WHERE cursor < $1 AND" + sqlAgentjJwtNotNullDesc + " $2",
	}
)

func (pg *Database) GetListenerAgents(info *paginator.BatchInfo) (a *model.Agents, err error) {
	defer returnErr("GetListenerAgents", &err)

	query, args := getBatchQuery(agentQueryInfo,
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
		agent, err = rowToAgent(rows)
		err2.Check(err)
		a.Agents = append(a.Agents, agent)
	}

	err2.Check(rows.Err())

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

	n = model.NewAgent(a)

	err2.Check(pg.doQuery(
		readRowToAgent(n),
		sqlAgentInsert,
		a.AgentID,
		a.Label,
		a.RawJWT,
	))

	n.TenantID = n.ID

	return
}

func rowToAgent(rows *sql.Rows) (a *model.Agent, err error) {
	a = model.NewAgent(nil)
	return a, readRowToAgent(a)(rows)
}

func readRowToAgent(a *model.Agent) func(*sql.Rows) error {
	return func(rows *sql.Rows) error {
		return rows.Scan(
			&a.ID, &a.AgentID, &a.Label, &a.RawJWT, &a.Created, &a.LastAccessed,
		)
	}
}

func (pg *Database) GetAgent(id, agentID *string) (a *model.Agent, err error) {
	defer err2.Annotate("GetAgent", &err)

	if id == nil && agentID == nil {
		panic(fmt.Errorf("either id or agent id is required"))
	}
	query := sqlAgentSelectByID
	queryID := id
	if id == nil {
		query = sqlAgentSelectByAgentID
		queryID = agentID
	}
	a = model.NewAgent(nil)

	err2.Check(pg.doQuery(readRowToAgent(a), query, *queryID))

	a.TenantID = a.ID

	return
}
