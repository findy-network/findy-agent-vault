package playground

import (
	"context"
	"time"

	"github.com/findy-network/findy-agent-vault/agency/model"
	"github.com/findy-network/findy-agent-vault/db/fake"
	"github.com/findy-network/findy-agent-vault/db/store"
	"github.com/findy-network/findy-agent-vault/paginator"
	"github.com/findy-network/findy-agent-vault/resolver/agent"
	"github.com/findy-network/findy-agent-vault/resolver/listen"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/google/uuid"
	"github.com/lainio/err2"
)

type Resolver struct {
	db store.DB
	*agent.Resolver
	*listen.Listener
}

func NewResolver(
	db store.DB,
	agentResolver *agent.Resolver,
	listener *listen.Listener,
) *Resolver {
	return &Resolver{db, agentResolver, listener}
}

func (r *Resolver) AddRandomEvent(ctx context.Context) (ok bool, err error) {
	utils.LogLow().Info("mutationResolver:addRandomEvent")

	_, err = r.GetAgent(ctx)
	err2.Check(err)

	// TODO
	return
}

func (r *Resolver) AddRandomMessage(ctx context.Context) (ok bool, err error) {
	utils.LogLow().Info("mutationResolver:addRandomMessage")

	tenant, err := r.GetAgent(ctx)
	err2.Check(err)

	res, err := r.db.GetConnections(
		&paginator.BatchInfo{Count: 1},
		tenant.ID,
	)
	err2.Check(err)

	if len(res.Connections) > 0 {
		connectionID := res.Connections[0].ID
		message := fake.Message(tenant.ID, connectionID)
		job := &model.JobInfo{
			TenantID:     tenant.ID,
			JobID:        uuid.New().String(),
			ConnectionID: connectionID,
		}

		r.AddMessage(job, &model.Message{Message: message.Message, SentByMe: message.SentByMe})
		ok = true
	}

	return
}

func (r *Resolver) AddRandomCredential(ctx context.Context) (ok bool, err error) {
	utils.LogLow().Info("mutationResolver:addRandomCredential")

	tenant, err := r.GetAgent(ctx)
	err2.Check(err)

	res, err := r.db.GetConnections(
		&paginator.BatchInfo{Count: 1},
		tenant.ID,
	)
	err2.Check(err)

	if len(res.Connections) > 0 {
		connectionID := res.Connections[0].ID
		credential := fake.Credential(tenant.ID, connectionID)
		job := &model.JobInfo{
			TenantID:     tenant.ID,
			JobID:        uuid.New().String(),
			ConnectionID: connectionID,
		}

		r.AddCredential(
			job,
			&model.Credential{
				Role:          credential.Role,
				SchemaID:      credential.SchemaID,
				CredDefID:     credential.CredDefID,
				Attributes:    credential.Attributes,
				InitiatedByUs: credential.InitiatedByUs,
			},
		)
		time.AfterFunc(time.Second, func() {
			now := utils.CurrentTimeMs()
			r.UpdateCredential(job, &model.CredentialUpdate{ApprovedMs: &now, IssuedMs: &now})
		})
		ok = true
	}

	return ok, err
}

func (r *Resolver) AddRandomProof(ctx context.Context) (ok bool, err error) {
	utils.LogLow().Info("mutationResolver:addRandomProof")

	tenant, err := r.GetAgent(ctx)
	err2.Check(err)

	res, err := r.db.GetConnections(
		&paginator.BatchInfo{Count: 1},
		tenant.ID,
	)
	err2.Check(err)

	if len(res.Connections) > 0 {
		connectionID := res.Connections[0].ID
		proof := fake.Proof(tenant.ID, connectionID)
		job := &model.JobInfo{
			TenantID:     tenant.ID,
			JobID:        uuid.New().String(),
			ConnectionID: connectionID,
		}

		r.AddProof(
			job,
			&model.Proof{
				Role:          proof.Role,
				Attributes:    proof.Attributes,
				InitiatedByUs: proof.InitiatedByUs,
			},
		)
		time.AfterFunc(time.Second, func() {
			now := utils.CurrentTimeMs()
			r.UpdateProof(job, &model.ProofUpdate{
				ApprovedMs: &now,
				VerifiedMs: &now,
			})
		})
		ok = true
	}
	return ok, err
}
