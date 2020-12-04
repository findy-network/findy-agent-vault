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
