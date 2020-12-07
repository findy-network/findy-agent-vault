package pg

import (
	"database/sql"

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
		", connection.created, approved, cursor FROM connection" +
		" INNER JOIN agent ON tenant_id = agent.id"
	sqlConnectionSelectByID = sqlConnectionSelect +
		" WHERE connection.id=$1 AND agent.agent_id=$2"
	sqlConnectionOrderByAsc  = " ORDER BY cursor ASC LIMIT"
	sqlConnectionOrderByDesc = " ORDER BY cursor ASC LIMIT"
	sqlConnectionSelectBatch = sqlConnectionSelect +
		" WHERE agent.agent_id=$1 " + sqlConnectionOrderByAsc + " $2"
	sqlConnectionSelectBatchTail = sqlConnectionSelect +
		" WHERE agent.agent_id=$1" + sqlConnectionOrderByDesc + " $2"
	sqlConnectionSelectBatchAfter = sqlConnectionSelect +
		" WHERE agent.agent_id=$1 AND connection.cursor > $2" + sqlConnectionOrderByAsc + " $3"
	sqlConnectionSelectBatchAfterTail = sqlConnectionSelect +
		" WHERE agent.agent_id=$1 AND connection.cursor > $2" + sqlConnectionOrderByDesc + " $3"
	sqlConnectionSelectBatchBefore = sqlConnectionSelect +
		" WHERE agent.agent_id=$1 AND connection.cursor < $2" + sqlConnectionOrderByAsc + " $3"
	sqlConnectionSelectBatchBeforeTail = sqlConnectionSelect +
		" WHERE agent.agent_id=$1 AND connection.cursor < $2" + sqlConnectionOrderByDesc + " $3"
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

func (p *Database) GetConnection(id, agentID string) (c *model.Connection, err error) {
	defer returnErr("GetConnection", &err)

	rows, err := p.db.Query(sqlConnectionSelectByID, id, agentID)
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

/*

if first, last missing, return error

Start from the greedy query: SELECT * FROM table ORDER BY created
If the after argument is provided, add id > parsed_cursor to the WHERE clause
If the before argument is provided, add id < parsed_cursor to the WHERE clause
If the first argument is provided, add ORDER BY id DESC LIMIT first+1 to the query
If the last argument is provided, add ORDER BY id ASC LIMIT last+1 to the query
If the last argument is provided, I reverse the order of the results
If the first argument is provided then I set hasPreviousPage: false (see spec for a description of this behavior).
If no less than first+1 results are returned, I set hasNextPage: true, otherwise I set it to false.
If the last argument is provided then I set hasNextPage: false (see spec for a description of this behavior).
If no less last+1 results are returned, I set hasPreviousPage: true, otherwise I set it to false.
*/

func (p *Database) GetConnections(info *paginator.BatchInfo, agentID string) (c *model.Connections, err error) {
	defer returnErr("GetConnections", &err)

	query := ""
	args := make([]interface{}, 0)
	args = append(args, agentID)

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

	//fmt.Println(query, args)
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

	return
}
