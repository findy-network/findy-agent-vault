package agency

import "github.com/findy-network/findy-agent-vault/graph/model"

type Listener interface {
	AddConnection(id, ourDID, theirDID, theirEndpoint, theirLabel string)
	AddMessage(connectionID, id, message string, sentByMe bool)
	UpdateMessage(connectionID, id, delivered bool)
	AddCredential(connectionID, id string, role model.CredentialRole, schemaID, credDefID string, attributes []*model.CredentialValue, initiatedByUs bool)
	AddProof(connectionID, id string, role model.ProofRole, attributes []*model.ProofAttribute, initiatedByUs bool)
}

type Agency interface {
	Init(l Listener)
	Invite() (string, string, error)
	Connect(invitation string) (string, error)
	SendMessage(connectionID, message string) (string, error)
}
