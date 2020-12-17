package resolver

import (
	"github.com/findy-network/findy-agent-vault/agency"
	"github.com/findy-network/findy-agent-vault/db/fake"
	"github.com/findy-network/findy-agent-vault/db/store"
	"github.com/findy-network/findy-agent-vault/db/store/mock"
	"github.com/findy-network/findy-agent-vault/db/store/pg"
	"github.com/findy-network/findy-agent-vault/graph/model"
)

//go:generate go run github.com/99designs/gqlgen

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	db             store.DB
	agency         agency.Agency
	eventObservers map[string]chan *model.EventEdge
}

func InitResolver(mockDB, fakeData bool) *Resolver {
	var db store.DB
	if mockDB {
		db = mock.InitState()
	} else {
		db = pg.InitDB("file://db/migrations", "5432", false)
	}

	// TODO: configure agency
	a := agency.Mock{}
	r := &Resolver{
		db:             db,
		agency:         &agency.Mock{},
		eventObservers: map[string]chan *model.EventEdge{},
	}

	a.Init(r)

	if fakeData {
		fake.AddData(db)
	}

	return r
}

func (r *Resolver) AddConnection(id, ourDID, theirDID, theirEndpoint, theirLabel string) {

}

func (r *Resolver) AddMessage(connectionID, id, message string, sentByMe bool) {

}

func (r *Resolver) UpdateMessage(connectionID, id, delivered bool) {

}

func (r *Resolver) AddCredential(
	connectionID, id string,
	role model.CredentialRole,
	schemaID, credDefID string,
	attributes []*model.CredentialValue,
	initiatedByUs bool,
) {

}

func (r *Resolver) UpdateCredential(connectionID, id string, approvedMs, issuedMs, failedMs *int64) {

}

func (r *Resolver) AddProof(connectionID, id string, role model.ProofRole, attributes []*model.ProofAttribute, initiatedByUs bool) {

}

func (r *Resolver) UpdateProof(connectionID, id string, approvedMs, verifiedMs, failedMs *int64) {

}
