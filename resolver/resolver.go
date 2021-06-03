package resolver

import (
	agency "github.com/findy-network/findy-agent-vault/agency/model"
	"github.com/findy-network/findy-agent-vault/db/fake"
	"github.com/findy-network/findy-agent-vault/db/store"
	"github.com/findy-network/findy-agent-vault/db/store/pg"
	"github.com/findy-network/findy-agent-vault/resolver/archive"
	"github.com/findy-network/findy-agent-vault/resolver/listen"
	"github.com/findy-network/findy-agent-vault/resolver/mutation"
	"github.com/findy-network/findy-agent-vault/resolver/query"
	"github.com/findy-network/findy-agent-vault/resolver/query/agent"
	"github.com/findy-network/findy-agent-vault/resolver/query/credential"
	"github.com/findy-network/findy-agent-vault/resolver/query/credentialconn"
	"github.com/findy-network/findy-agent-vault/resolver/query/event"
	"github.com/findy-network/findy-agent-vault/resolver/query/eventconn"
	"github.com/findy-network/findy-agent-vault/resolver/query/job"
	"github.com/findy-network/findy-agent-vault/resolver/query/jobconn"
	"github.com/findy-network/findy-agent-vault/resolver/query/message"
	"github.com/findy-network/findy-agent-vault/resolver/query/messageconn"
	"github.com/findy-network/findy-agent-vault/resolver/query/pairwise"
	"github.com/findy-network/findy-agent-vault/resolver/query/pairwiseconn"
	"github.com/findy-network/findy-agent-vault/resolver/query/proof"
	"github.com/findy-network/findy-agent-vault/resolver/query/proofconn"
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
	proofConnection      *proofconn.Resolver
	proof                *proof.Resolver
	query                *query.Resolver
}

type Resolver struct {
	db       store.DB
	agency   agency.Agency
	updater  *update.Updater
	listener *listen.Listener
	archiver *archive.Archiver

	resolvers *controller
}

func InitResolverWithDB(config *utils.Configuration, coreAgency agency.Agency, db store.DB) *Resolver {
	r := &Resolver{db: db}

	r.agency = coreAgency

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
	r.archiver = archive.NewArchiver(db)
	r.agency.Init(r.listener, agentResolver.FetchAgents(), r.archiver, config)

	return r
}

func InitResolver(config *utils.Configuration, coreAgency agency.Agency) *Resolver {
	db := pg.InitDB(config, false, false)
	if config.GenerateFakeData {
		fake.AddData(db)
	}
	return InitResolverWithDB(config, coreAgency, db)
}

// For testing
func (r *Resolver) Store() store.DB {
	return r.db
}
