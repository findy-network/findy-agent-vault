package model

import "time"

type base struct {
	ID      string `faker:"uuid_hyphenated"`
	Created time.Time
}

type Agent struct {
	*base
	AgentID string
	Label   string `faker:"first_name"`
}

func NewAgent() *Agent { return &Agent{base: &base{}} }

type Connection struct {
	*base
	TenantID      string
	OurDid        string
	TheirDid      string
	TheirEndpoint string `faker:"url"`
	TheirLabel    string `faker:"organisationLabel"`
	Invited       bool
	Approved      *time.Time
	Cursor        uint64
}

func NewConnection() *Connection { return &Connection{base: &base{}} }

func (c *Connection) Copy() (n *Connection) {
	n = NewConnection()
	n.TenantID = c.TenantID
	n.OurDid = c.OurDid
	n.TheirDid = c.TheirDid
	n.TheirEndpoint = c.TheirEndpoint
	n.TheirLabel = c.TheirLabel
	n.Invited = c.Invited
	return n
}
