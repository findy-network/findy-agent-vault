package agent

import (
	"context"

	agency "github.com/findy-network/findy-agent-vault/agency/model"
	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/db/store"
	"github.com/findy-network/findy-agent-vault/paginator"
	"github.com/lainio/err2"
)

type Resolver struct {
	db     store.DB
	agency agency.Agency
}

func NewResolver(db store.DB, agencyInstance agency.Agency) *Resolver {
	return &Resolver{db: db, agency: agencyInstance}
}

func (r *Resolver) GetAgent(ctx context.Context) (agent *model.Agent, err error) {
	err2.Return(&err)

	agent, err = store.GetAgent(ctx, r.db)
	err2.Check(err)

	// make sure we are listening events for this agent
	if agent.IsNewOnboard() {
		err2.Check(r.agency.AddAgent(r.AgencyAuth(agent)))
	}
	return
}

func (r *Resolver) AgencyAuth(agent *model.Agent) *agency.Agent {
	return &agency.Agent{
		Label:    agent.Label,
		RawJWT:   agent.RawJWT,
		TenantID: agent.ID,
		AgentID:  agent.AgentID,
	}
}

func (r *Resolver) FetchAgents() []*agency.Agent {
	nextPage := true
	after := uint64(0)
	allAgents := make([]*model.Agent, 0)
	for nextPage {
		agents, err := r.db.GetListenerAgents(&paginator.BatchInfo{Count: 100, After: after})
		if err != nil && store.ErrorCode(err) != store.ErrCodeNotFound {
			panic(err)
		}
		count := len(agents.Agents)
		if count > 0 {
			allAgents = append(allAgents, agents.Agents...)
			nextPage = agents.HasNextPage
			after = agents.Agents[count-1].Cursor
		} else {
			nextPage = false
		}
	}

	listenerAgents := make([]*agency.Agent, len(allAgents))
	for index := range allAgents {
		listenerAgents[index] = r.AgencyAuth(allAgents[index])
	}
	return listenerAgents
}
