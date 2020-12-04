package resolver

import (
	"github.com/findy-network/findy-agent-vault/db/db/pg"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

//go:generate go run github.com/99designs/gqlgen

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{}

func InitResolver() *Resolver {
	//listener := &agencyListener{}
	//agency.Instance.Init(listener)

	pg.InitDb("file://db/migrations", "5432", false)

	return &Resolver{}
}
