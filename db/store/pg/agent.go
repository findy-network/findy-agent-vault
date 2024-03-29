//nolint:goconst
package pg

import (
	"database/sql"
	"fmt"
	"sort"

	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/paginator"
	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
)

const (
	sqlAgentFields          = "id, agent_id, label, raw_jwt, created, last_accessed, cursor"
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
	defer err2.Handle(&err, "GetListenerAgents")

	query, args := getBatchQuery(agentQueryInfo,
		info,
		"",
		[]interface{}{},
	)

	a = &model.Agents{
		Agents:          make([]*model.Agent, 0),
		HasNextPage:     false,
		HasPreviousPage: false,
	}
	var agent *model.Agent
	try.To(pg.doRowsQuery(func(rows *sql.Rows) (err error) {
		defer err2.Handle(&err)
		agent = try.To1(rowToAgent(rows))
		a.Agents = append(a.Agents, agent)
		return
	}, query, args...))

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

func (pg *Database) AddAgent(a *model.Agent) (newAgent *model.Agent, err error) {
	defer err2.Handle(&err, "AddAgent")

	const sqlAgentInsert = "INSERT INTO agent (agent_id, label, raw_jwt) VALUES ($1, $2, $3) " +
		"ON CONFLICT (agent_id) DO UPDATE SET " +
		"last_accessed = (now() at time zone 'UTC'), raw_jwt = $4 " +
		"RETURNING " + sqlAgentFields

	newAgent = &model.Agent{}
	*newAgent = *a

	try.To(pg.doRowQuery(
		readRowToAgent(newAgent),
		sqlAgentInsert,
		a.AgentID,
		a.Label,
		a.RawJWT,
		a.RawJWT,
	))

	newAgent.TenantID = newAgent.ID

	return
}

func rowToAgent(rows *sql.Rows) (a *model.Agent, err error) {
	a = &model.Agent{}
	return a, readRowToAgent(a)(rows)
}

func readRowToAgent(a *model.Agent) func(*sql.Rows) error {
	return func(rows *sql.Rows) error {
		return rows.Scan(
			&a.ID, &a.AgentID, &a.Label, &a.RawJWT, &a.Created, &a.LastAccessed, &a.Cursor,
		)
	}
}

func (pg *Database) GetAgent(id, agentID *string) (a *model.Agent, err error) {
	defer err2.Handle(&err, "GetAgent")

	if id == nil && agentID == nil {
		panic(fmt.Errorf("either id or agent id is required"))
	}
	query := sqlAgentSelectByID
	queryID := id
	if id == nil {
		query = sqlAgentSelectByAgentID
		queryID = agentID
	}
	a = &model.Agent{}

	try.To(pg.doRowQuery(readRowToAgent(a), query, *queryID))

	a.TenantID = a.ID

	return
}
