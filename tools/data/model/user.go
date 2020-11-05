package model

import "github.com/findy-network/findy-agent-vault/graph/model"

type InternalUser struct {
	ID   string `faker:"uuid_hyphenated"`
	Name string `faker:"first_name"`
}

func (u *InternalUser) ToNode() *model.User {
	return &model.User{
		ID:   u.ID,
		Name: u.Name,
	}
}
