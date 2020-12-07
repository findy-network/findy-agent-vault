package resolver

import (
	"github.com/findy-network/findy-agent-vault/db/db"
	"github.com/findy-network/findy-agent-vault/db/db/pg"
	"github.com/findy-network/findy-agent-vault/db/fake"
)

//go:generate go run github.com/99designs/gqlgen

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	db db.Db
}

func InitResolver() *Resolver {
	//listener := &agencyListener{}
	//agency.Instance.Init(listener)

	db := pg.InitDB("file://db/migrations", "5432", false)

	r := &Resolver{db}
	fake.AddData(db)

	return r
}
