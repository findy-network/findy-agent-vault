package pg

import (
	"database/sql"
	"sort"

	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/paginator"
	"github.com/lainio/err2"
)

const (
	sqlConnectionFields = "tenant_id, our_did, their_did, their_endpoint, their_label, invited"
	sqlConnectionInsert = "INSERT INTO connection " +
		"(" + sqlConnectionFields + ") " +
		"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created, cursor"
	sqlConnectionSelect = "SELECT connection.id, " +
		sqlConnectionFields +
		", connection.created, approved, cursor FROM connection"
	sqlConnectionSelectByID = sqlConnectionSelect +
		" WHERE connection.id=$1 AND tenant_id=$2"
	sqlConnectionOrderByAsc  = " ORDER BY cursor ASC LIMIT"
	sqlConnectionOrderByDesc = " ORDER BY cursor DESC LIMIT"
	sqlConnectionSelectBatch = sqlConnectionSelect +
		" WHERE tenant_id=$1 " + sqlConnectionOrderByAsc + " $2"
	sqlConnectionSelectBatchTail = sqlConnectionSelect +
		" WHERE tenant_id=$1" + sqlConnectionOrderByDesc + " $2"
	sqlConnectionSelectBatchAfter = sqlConnectionSelect +
		" WHERE tenant_id=$1 AND connection.cursor > $2" + sqlConnectionOrderByAsc + " $3"
	sqlConnectionSelectBatchAfterTail = sqlConnectionSelect +
		" WHERE tenant_id=$1 AND connection.cursor > $2" + sqlConnectionOrderByDesc + " $3"
	sqlConnectionSelectBatchBefore = sqlConnectionSelect +
		" WHERE tenant_id=$1 AND connection.cursor < $2" + sqlConnectionOrderByAsc + " $3"
	sqlConnectionSelectBatchBeforeTail = sqlConnectionSelect +
		" WHERE tenant_id=$1 AND connection.cursor < $2" + sqlConnectionOrderByDesc + " $3"
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

	query := ""
	args := make([]interface{}, 0)
	args = append(args, tenantID)

	if info.Tail {
		query = sqlConnectionSelectBatchTail
		if info.After > 0 {
			query = sqlConnectionSelectBatchAfterTail
		} else if info.Before > 0 {
			query = sqlConnectionSelectBatchBeforeTail
		}
	} else {
		query = sqlConnectionSelectBatch
		if info.After > 0 {
			query = sqlConnectionSelectBatchAfter
		} else if info.Before > 0 {
			query = sqlConnectionSelectBatchBefore
		}
	}
	if info.After > 0 {
		args = append(args, info.After)
	} else if info.Before > 0 {
		args = append(args, info.Before)
	}

	args = append(args, info.Count+1)

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

	return
}
