package resolver

import (
	agencys "github.com/findy-network/findy-agent-vault/agency"
	agency "github.com/findy-network/findy-agent-vault/agency/model"
	"github.com/findy-network/findy-agent-vault/db/fake"
	"github.com/findy-network/findy-agent-vault/db/store"
	"github.com/findy-network/findy-agent-vault/db/store/mock"
	"github.com/findy-network/findy-agent-vault/db/store/pg"
	"github.com/findy-network/findy-agent-vault/resolver/agent"
	"github.com/findy-network/findy-agent-vault/resolver/credential"
	"github.com/findy-network/findy-agent-vault/resolver/credentialconn"
	"github.com/findy-network/findy-agent-vault/resolver/event"
	"github.com/findy-network/findy-agent-vault/resolver/eventconn"
	"github.com/findy-network/findy-agent-vault/resolver/job"
	"github.com/findy-network/findy-agent-vault/resolver/jobconn"
	"github.com/findy-network/findy-agent-vault/resolver/listen"
	"github.com/findy-network/findy-agent-vault/resolver/message"
	"github.com/findy-network/findy-agent-vault/resolver/messageconn"
	"github.com/findy-network/findy-agent-vault/resolver/mutation"
	"github.com/findy-network/findy-agent-vault/resolver/pairwise"
	"github.com/findy-network/findy-agent-vault/resolver/pairwiseconn"
	"github.com/findy-network/findy-agent-vault/resolver/playground"
	"github.com/findy-network/findy-agent-vault/resolver/proof"
	"github.com/findy-network/findy-agent-vault/resolver/proofconn"
	"github.com/findy-network/findy-agent-vault/resolver/query"
	"github.com/findy-network/findy-agent-vault/resolver/update"
	"github.com/findy-network/findy-agent-vault/utils"
)

//go:generate go run github.com/99designs/gqlgen

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type controller struct {
	agent                *agent.Resolver
	message              *message.Resolver
	credentialConnection *credentialconn.Resolver
	credential           *credential.Resolver
	eventConnection      *eventconn.Resolver
	event                *event.Resolver
	jobConnection        *jobconn.Resolver
	job                  *job.Resolver
	messageConnection    *messageconn.Resolver
	mutation             *mutation.Resolver
	pairwiseConnection   *pairwiseconn.Resolver
	pairwise             *pairwise.Resolver
	playground           *playground.Resolver
	proofConnection      *proofconn.Resolver
	proof                *proof.Resolver
	query                *query.Resolver
}

type Resolver struct {
	db       store.DB
	agency   agency.Agency
	updater  *update.Updater
	listener *listen.Listener

	resolvers *controller
}

func InitResolver(config *utils.Configuration) *Resolver {
	var db store.DB
	if config.UseMockDB {
		db = mock.InitState()
	} else {
		db = pg.InitDB(
			"file://db/migrations",
			config.DBHost,
			config.DBPassword,
			config.DBPort,
			false,
		)
	}
	if config.GenerateFakeData {
		fake.AddData(db)
	}

	r := &Resolver{db: db}

	aType := agencys.AgencyTypeMock
	if !config.UseMockAgency {
		aType = agencys.AgencyTypeFindyGRPC
	}
	r.agency = agencys.Create(aType)

	agentResolver := agent.NewResolver(db, r.agency)
	updater := update.NewUpdater(db, agentResolver)
	r.resolvers = &controller{
		agent:                agentResolver,
		message:              message.NewResolver(db, agentResolver),
		credentialConnection: credentialconn.NewResolver(db, agentResolver),
		credential:           credential.NewResolver(db, agentResolver),
		eventConnection:      eventconn.NewResolver(db, agentResolver),
		event:                event.NewResolver(db, agentResolver),
		jobConnection:        jobconn.NewResolver(db, agentResolver),
		job:                  job.NewResolver(db, agentResolver),
		messageConnection:    messageconn.NewResolver(db, agentResolver),
		mutation:             mutation.NewResolver(db, r.agency, agentResolver, updater),
		proofConnection:      proofconn.NewResolver(db, agentResolver),
		proof:                proof.NewResolver(db, agentResolver),
		pairwiseConnection:   pairwiseconn.NewResolver(db, agentResolver),
		pairwise:             pairwise.NewResolver(db, agentResolver),
		query:                query.NewResolver(db, agentResolver),
	}
	r.updater = updater

	r.listener = listen.NewListener(db, r.updater)
	r.agency.Init(r.listener, agentResolver.FetchAgents(), config)

	if config.UsePlayground {
		r.resolvers.playground = playground.NewResolver(db, agentResolver, r.listener)
	}

	return r
}

// For testing
func (r *Resolver) Store() store.DB {
	return r.db
}
