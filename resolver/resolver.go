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
	db db.DB
}

func InitResolver() *Resolver {
	store := pg.InitDB("file://db/migrations", "5432", false)

	r := &Resolver{db: store}
	fake.AddData(store)

	return r
}
