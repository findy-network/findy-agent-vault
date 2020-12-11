package resolver

import (
	"github.com/findy-network/findy-agent-vault/db/fake"
	"github.com/findy-network/findy-agent-vault/db/store"
	"github.com/findy-network/findy-agent-vault/db/store/pg"
)

//go:generate go run github.com/99designs/gqlgen

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	db store.DB
}

func InitResolver() *Resolver {
	db := pg.InitDB("file://db/migrations", "5432", false)

	r := &Resolver{db: db}
	fake.AddData(db)

	return r
}
