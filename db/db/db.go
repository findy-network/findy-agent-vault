package db

import (
	"context"

	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/paginator"
	"github.com/findy-network/findy-agent-vault/utils"
)

func GetAgent(ctx context.Context, db DB) (*model.Agent, error) {
	token := utils.ParseToken(ctx)
	return db.AddAgent(&model.Agent{AgentID: token.AgentID, Label: token.Label})
}

type DB interface {
	Close()

	AddAgent(a *model.Agent) (*model.Agent, error)
	GetAgent(id, agentID *string) (*model.Agent, error)

	AddConnection(c *model.Connection) (*model.Connection, error)
	GetConnection(id string, tenantID string) (*model.Connection, error)
	GetConnections(info *paginator.BatchInfo, tenantID string) (connections *model.Connections, err error)

	/*AddMessage(connectionID, id, message string, sentByMe bool)
	UpdateMessage(connectionID, id, delivered bool)

	AddCredential(
		connectionID, id string,
		role model.CredentialRole,
		schemaID, credDefID string,
		attributes []*model.CredentialValue,
		initiatedByUs bool,
	)
	UpdateCredential(connectionID, id string, approvedMs, issuedMs, failedMs *int64)

	AddProof(connectionID, id string, role model.ProofRole, attributes []*model.ProofAttribute, initiatedByUs bool)
	UpdateProof(connectionID, id string, approvedMs, verifiedMs, failedMs *int64)*/
}
