package model

import "time"

type Agent struct {
	ID      string `faker:"uuid_hyphenated"`
	AgentID string
	Label   string `faker:"first_name"`
	Created time.Time
}
