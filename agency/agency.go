package agency

import (
	"github.com/findy-network/findy-agent-vault/graph/model"
)

type JobInfo struct {
	TenantID     string
	JobID        string
	ConnectionID string
}

type Listener interface {
	AddConnection(job *JobInfo, ourDID, theirDID, theirEndpoint, theirLabel string)

	AddMessage(job *JobInfo, message string, sentByMe bool)
	UpdateMessage(job *JobInfo, delivered bool)

	AddCredential(
		job *JobInfo,
		role model.CredentialRole,
		schemaID, credDefID string,
		attributes []*model.CredentialValue,
		initiatedByUs bool,
	)
	UpdateCredential(job *JobInfo, approvedMs, issuedMs, failedMs *int64)

	AddProof(job *JobInfo, role model.ProofRole, attributes []*model.ProofAttribute, initiatedByUs bool)
	UpdateProof(job *JobInfo, approvedMs, verifiedMs, failedMs *int64)
}

type Agent struct {
	RawJWT   string
	TenantID string
	AgentID  string
}

type Agency interface {
	Init(l Listener)

	Invite(a *Agent) (string, string, error)
	Connect(a *Agent, invitation string) (string, error)
	SendMessage(a *Agent, connectionID, message string) (string, error)

	ResumeCredentialOffer(a *Agent, id string, accept bool) error
	ResumeProofRequest(a *Agent, id string, accept bool) error
}
