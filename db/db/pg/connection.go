package pg

import (
	"database/sql"
	"sort"

	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/paginator"
	"github.com/lainio/err2"
)

func sqlWhereTenantAsc(orderBy string) string {
	return " WHERE tenant_id=$1 " + sqlOrderByAsc(orderBy)
}

func sqlWhereTenantDesc(orderBy string) string {
	return " WHERE tenant_id=$1 " + sqlOrderByDesc(orderBy)
}

func sqlWhereTenantAscAfter(orderBy string) string {
	return " WHERE tenant_id=$1 AND cursor > $2" + sqlOrderByAsc(orderBy)
}

func sqlWhereTenantDescAfter(orderBy string) string {
	return " WHERE tenant_id=$1 AND cursor > $2" + sqlOrderByDesc(orderBy)
}

func sqlWhereTenantAscBefore(orderBy string) string {
	return " WHERE tenant_id=$1 AND cursor < $2" + sqlOrderByAsc(orderBy)
}

func sqlWhereTenantDescBefore(orderBy string) string {
	return " WHERE tenant_id=$1 AND cursor < $2" + sqlOrderByDesc(orderBy)
}

const (
	sqlConnectionFields = "tenant_id, our_did, their_did, their_endpoint, their_label, invited"
	sqlConnectionInsert = "INSERT INTO connection " + "(" + sqlConnectionFields + ") " +
		"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created, cursor"
	sqlConnectionSelect = "SELECT id, " + sqlConnectionFields + ", created, approved, cursor FROM connection"
)

func (p *Database) AddConnection(c *model.Connection) (n *model.Connection, err error) {
	defer returnErr("AddConnection", &err)

	rows, err := p.db.Query(
		sqlConnectionInsert,
		c.TenantID,
		c.OurDid,
		c.TheirDid,
		c.TheirEndpoint,
		c.TheirLabel,
		c.Invited,
	)
	err2.Check(err)
	defer rows.Close()

	n = c.Copy()
	if rows.Next() {
		err = rows.Scan(&n.ID, &n.Created, &n.Cursor)
		err2.Check(err)
	}

	err = rows.Err()
	err2.Check(err)

	return
}

func readRowToConnection(rows *sql.Rows) (c *model.Connection, err error) {
	c = model.NewConnection()
	err = rows.Scan(
		&c.ID,
		&c.TenantID,
		&c.OurDid,
		&c.TheirDid,
		&c.TheirEndpoint,
		&c.TheirLabel,
		&c.Invited,
		&c.Created,
		&c.Approved,
		&c.Cursor,
	)
	return
}

func (p *Database) GetConnection(id, tenantID string) (c *model.Connection, err error) {
	defer returnErr("GetConnection", &err)

	const sqlConnectionSelectByID = sqlConnectionSelect + " WHERE id=$1 AND tenant_id=$2"

	rows, err := p.db.Query(sqlConnectionSelectByID, id, tenantID)
	err2.Check(err)
	defer rows.Close()

	if rows.Next() {
		c, err = readRowToConnection(rows)
		err2.Check(err)
	}

	err = rows.Err()
	err2.Check(err)

	return
}

func (p *Database) GetConnections(info *paginator.BatchInfo, tenantID string) (c *model.Connections, err error) {
	defer returnErr("GetConnections", &err)

	query, args := getBatchQuery(&queryInfo{
		Asc:        sqlConnectionSelect + sqlWhereTenantAsc("") + " $2",
		Desc:       sqlConnectionSelect + sqlWhereTenantDesc("") + " $2",
		AfterAsc:   sqlConnectionSelect + sqlWhereTenantAscAfter("") + " $3",
		AfterDesc:  sqlConnectionSelect + sqlWhereTenantDescAfter("") + " $3",
		BeforeAsc:  sqlConnectionSelect + sqlWhereTenantAscBefore("") + " $3",
		BeforeDesc: sqlConnectionSelect + sqlWhereTenantDescBefore("") + " $3",
	},
		info,
		[]interface{}{tenantID},
	)

	rows, err := p.db.Query(query, args...)
	err2.Check(err)
	defer rows.Close()

	c = &model.Connections{
		Connections:     make([]*model.Connection, 0),
		HasNextPage:     false,
		HasPreviousPage: false,
	}
	var connection *model.Connection
	for rows.Next() {
		connection, err = readRowToConnection(rows)
		err2.Check(err)
		c.Connections = append(c.Connections, connection)
	}

	err = rows.Err()
	err2.Check(err)

	if info.Count < len(c.Connections) {
		c.Connections = c.Connections[:info.Count]
		if info.Tail {
			c.HasPreviousPage = true
		} else {
			c.HasNextPage = true
		}
	}

	if info.After > 0 {
		c.HasPreviousPage = true
	}
	if info.Before > 0 {
		c.HasNextPage = true
	}

	// Reverse order for tail first
	if info.Tail {
		sort.Slice(c.Connections, func(i, j int) bool {
			return c.Connections[i].Created.Sub(c.Connections[j].Created) < 0
		})
	}

	return c, err
}

func (p *Database) GetConnectionCount(tenantID string) (count int, err error) {
	defer returnErr("GetConnectionCount", &err)

	const sqlConnectionSelectCount = "SELECT count(id) FROM connection WHERE tenant_id=$1"

	rows, err := p.db.Query(sqlConnectionSelectCount, tenantID)
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
