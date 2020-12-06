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
	sqlConnectionSelect     = "SELECT id, " + sqlConnectionFields + ", created, approved, cursor FROM connection"
	sqlConnectionSelectByID = sqlConnectionSelect + " WHERE id=$1"
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

func (p *Database) GetConnection(id string) (c *model.Connection, err error) {
	defer returnErr("GetConnection", &err)

	rows, err := p.db.Query(sqlConnectionSelectByID, id)
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
