package db

import "github.com/findy-network/findy-agent-vault/db/model"

type Db interface {
	Close()

	AddAgent(agentId, label string) error
	GetAgent(id, agentID *string) (*model.Agent, error)

	/*AddConnection(id, ourDID, theirDID, theirEndpoint, theirLabel string)

	AddMessage(connectionID, id, message string, sentByMe bool)
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
