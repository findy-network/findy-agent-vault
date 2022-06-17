package pg

import (
	"database/sql"
	"sort"

	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/db/store"
	"github.com/findy-network/findy-agent-vault/paginator"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
)

var (
	connectionFields        = []string{"id", "tenant_id", "our_did", "their_did", "their_endpoint", "their_label", "invited", "archived"}
	connectionExtraFields   = []string{"created", "approved", "cursor"}
	sqlConnectionBaseFields = sqlFields("", connectionFields)
	sqlConnectionInsert     = "INSERT INTO connection " + "(" + sqlConnectionBaseFields + ") " +
		"VALUES (" + sqlArguments(connectionFields) + ") RETURNING " + sqlInsertFields
	sqlConnectionSelect = "SELECT " + sqlConnectionBaseFields + ", " + sqlFields("connection", connectionExtraFields) + " FROM connection"

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
	defer err2.Annotate("getConnectionForObject", &err)

	sqlConnectionJoinSelect := "SELECT " + sqlFields("connection", connectionFields) +
		", connection.created, connection.approved, connection.cursor FROM connection"
	sqlConnectionSelectByObjectID := sqlConnectionJoinSelect +
		" INNER JOIN " + objectName + " ON " + objectName +
		"." + columnName + "=connection.id WHERE " + objectName + ".id = $1 AND connection.tenant_id = $2"

	c = &model.Connection{}
	try.To(pg.doRowQuery(
		readRowToConnection(c),
		sqlConnectionSelectByObjectID,
		objectID,
		tenantID,
	))

	return
}

func (pg *Database) AddConnection(c *model.Connection) (newConnection *model.Connection, err error) {
	defer err2.Annotate("AddConnection", &err)

	newConnection = &model.Connection{}
	*newConnection = *c
	try.To(pg.doRowQuery(
		func(rows *sql.Rows) error {
			return rows.Scan(&newConnection.ID, &newConnection.Created, &newConnection.Cursor)
		},
		sqlConnectionInsert,
		c.ID,
		c.TenantID,
		c.OurDid,
		c.TheirDid,
		c.TheirEndpoint,
		c.TheirLabel,
		c.Invited,
		c.Archived,
	))

	return
}

func rowToConnection(rows *sql.Rows) (c *model.Connection, err error) {
	c = &model.Connection{}
	return c, readRowToConnection(c)(rows)
}

func readRowToConnection(c *model.Connection) func(*sql.Rows) error {
	return func(rows *sql.Rows) error {
		err := rows.Scan(
			&c.ID,
			&c.TenantID,
			&c.OurDid,
			&c.TheirDid,
			&c.TheirEndpoint,
			&c.TheirLabel,
			&c.Invited,
			&c.Archived,
			&c.Created,
			&c.Approved,
			&c.Cursor,
		)
		return err
	}
}

func (pg *Database) GetConnection(id, tenantID string) (c *model.Connection, err error) {
	defer err2.Annotate("GetConnection", &err)

	sqlConnectionSelectByID := sqlConnectionSelect + " WHERE id=$1 AND tenant_id=$2"

	c = &model.Connection{}
	try.To(pg.doRowQuery(
		readRowToConnection(c),
		sqlConnectionSelectByID,
		id,
		tenantID,
	))

	return
}

func (pg *Database) GetConnections(info *paginator.BatchInfo, tenantID string) (c *model.Connections, err error) {
	defer err2.Annotate("GetConnections", &err)

	query, args := getBatchQuery(connectionQueryInfo, info, tenantID, []interface{}{})

	c = &model.Connections{
		Connections:     make([]*model.Connection, 0),
		HasNextPage:     false,
		HasPreviousPage: false,
	}
	var connection *model.Connection
	if err = pg.doRowsQuery(func(rows *sql.Rows) (err error) {
		defer err2.Return(&err)
		connection = try.To1(rowToConnection(rows))
		c.Connections = append(c.Connections, connection)
		return
	}, query, args...); err != nil && store.ErrorCode(err) == store.ErrCodeNotFound {
		// no connections is not an error
		err = nil
	}
	try.To(err)

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
	defer err2.Annotate("GetCredentialCount", &err)
	count, err = pg.getCount(
		"connection",
		" WHERE tenant_id=$1 ",
		"",
		tenantID,
		nil,
	)
	try.To(err)
	return
}

func (pg *Database) ArchiveConnection(id, tenantID string) (err error) {
	defer err2.Annotate("ArchiveConnection", &err)

	var (
		sqlConnectionArchive = "UPDATE connection SET archived=$1 WHERE id = $2 and tenant_id = $3 RETURNING " +
			sqlConnectionBaseFields + "," + sqlFields("", connectionExtraFields)
	)

	now := utils.CurrentTime()
	newConnection := &model.Connection{}
	try.To(pg.doRowQuery(
		readRowToConnection(newConnection),
		sqlConnectionArchive,
		now,
		id,
		tenantID,
	))
	return
}
