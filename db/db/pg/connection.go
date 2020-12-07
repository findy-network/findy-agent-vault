package pg

import (
	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/lainio/err2"
)

const (
	sqlConnectionFields = "tenant_id, our_did, their_did, their_endpoint, their_label, invited"
	sqlConnectionInsert = "INSERT INTO connection " +
		"(" + sqlConnectionFields + ") " +
		"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created, cursor"
	sqlConnectionSelect     = "SELECT connection.id, " + sqlConnectionFields + ", connection.created, approved, cursor FROM connection"
	sqlConnectionSelectByID = sqlConnectionSelect +
		" INNER JOIN agent ON tenant_id = agent.id WHERE connection.id=$1 AND agent.agent_id=$2"
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

func (p *Database) GetConnection(id, agentID string) (c *model.Connection, err error) {
	defer returnErr("GetConnection", &err)

	rows, err := p.db.Query(sqlConnectionSelectByID, id, agentID)
	err2.Check(err)
	defer rows.Close()

	c = model.NewConnection()
	if rows.Next() {
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
		err2.Check(err)
	}

	err = rows.Err()
	err2.Check(err)

	return
}
