package pg

import (
	"database/sql"
	"fmt"
	"sort"

	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/paginator"
	"github.com/lainio/err2"
)

var (
	connectionFields        = []string{"id", "tenant_id", "our_did", "their_did", "their_endpoint", "their_label", "invited"}
	sqlConnectionBaseFields = sqlFields("", connectionFields)
	sqlConnectionInsert     = "INSERT INTO connection " + "(" + sqlConnectionBaseFields + ") " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created, cursor"
	sqlConnectionSelect = "SELECT " + sqlConnectionBaseFields +
		", connection.created, connection.approved, connection.cursor FROM connection"

	connectionQueryInfo = &queryInfo{
		Asc:        sqlConnectionSelect + " WHERE tenant_id=$1 " + sqlOrderByCursorAsc + " $2",
		Desc:       sqlConnectionSelect + " WHERE tenant_id=$1 " + sqlOrderByCursorDesc + " $2",
		AfterAsc:   sqlConnectionSelect + " WHERE tenant_id=$1 AND cursor > $2" + sqlOrderByCursorAsc + " $3",
		AfterDesc:  sqlConnectionSelect + " WHERE tenant_id=$1 AND cursor > $2" + sqlOrderByCursorDesc + " $3",
		BeforeAsc:  sqlConnectionSelect + " WHERE tenant_id=$1 AND cursor < $2" + sqlOrderByCursorAsc + " $3",
		BeforeDesc: sqlConnectionSelect + " WHERE tenant_id=$1 AND cursor < $2" + sqlOrderByCursorDesc + " $3",
	}
)

func (pg *Database) getConnectionForObject(objectName, columnName, objectID, tenantID string) (c *model.Connection, err error) {
	defer returnErr("getConnectionForObject", &err)

	sqlConnectionJoinSelect := "SELECT " + sqlFields("connection", connectionFields) +
		", connection.created, connection.approved, connection.cursor FROM connection"
	sqlConnectionSelectByObjectID := sqlConnectionJoinSelect +
		" INNER JOIN " + objectName + " ON " + objectName +
		"." + columnName + "=connection.id WHERE " + objectName + ".id = $1 AND connection.tenant_id = $2"

	rows, err := pg.db.Query(sqlConnectionSelectByObjectID, objectID, tenantID)
	err2.Check(err)
	defer rows.Close()

	if rows.Next() {
		c, err = readRowToConnection(rows)
	} else {
		err = fmt.Errorf("not found connection for %s id %s", objectName, objectID)
	}
	err2.Check(err)

	return
}

func (pg *Database) AddConnection(c *model.Connection) (n *model.Connection, err error) {
	defer returnErr("AddConnection", &err)

	rows, err := pg.db.Query(
		sqlConnectionInsert,
		c.ID,
		c.TenantID,
		c.OurDid,
		c.TheirDid,
		c.TheirEndpoint,
		c.TheirLabel,
		c.Invited,
	)
	err2.Check(err)
	defer rows.Close()

	n = model.NewConnection(c.ID, c.TenantID, c)
	if rows.Next() {
		err = rows.Scan(&n.ID, &n.Created, &n.Cursor)
	} else {
		err = fmt.Errorf("no rows returned from insert connection query")
	}
	err2.Check(err)

	return
}

func readRowToConnection(rows *sql.Rows) (c *model.Connection, err error) {
	c = model.EmptyConnection()
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

func (pg *Database) GetConnection(id, tenantID string) (c *model.Connection, err error) {
	defer returnErr("GetConnection", &err)

	sqlConnectionSelectByID := sqlConnectionSelect + " WHERE id=$1 AND tenant_id=$2"

	rows, err := pg.db.Query(sqlConnectionSelectByID, id, tenantID)
	err2.Check(err)
	defer rows.Close()

	if rows.Next() {
		c, err = readRowToConnection(rows)
	} else {
		err = fmt.Errorf("no rows returned from select connection query (%s)", id)
	}
	err2.Check(err)

	return
}

func (pg *Database) GetConnections(info *paginator.BatchInfo, tenantID string) (c *model.Connections, err error) {
	defer returnErr("GetConnections", &err)

	query, args := getBatchQuery(connectionQueryInfo, info, tenantID, []interface{}{})

	rows, err := pg.db.Query(query, args...)
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

func (pg *Database) GetConnectionCount(tenantID string) (count int, err error) {
	defer returnErr("GetCredentialCount", &err)
	count, err = pg.getCount(
		"connection",
		" WHERE tenant_id=$1 ",
		"",
		tenantID,
		nil,
	)
	err2.Check(err)
	return
}

func (pg *Database) ArchiveConnection(c *model.Connection) (*model.Connection, error) {
	return nil, nil
}
