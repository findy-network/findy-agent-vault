package model

import "github.com/findy-network/findy-agent-vault/graph/model"

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
	Label    string
	RawJWT   string
	TenantID string
	AgentID  string
}

type Agency interface {
	Init(l Listener, agents []*Agent)
	AddAgent(agent *Agent) error

	Invite(a *Agent) (string, string, error)
	Connect(a *Agent, invitation string) (string, error)
	SendMessage(a *Agent, connectionID, message string) (string, error)

	ResumeCredentialOffer(a *Agent, job *JobInfo, accept bool) error
	ResumeProofRequest(a *Agent, job *JobInfo, accept bool) error
}

// TODO: use invitation struct defined by agency
type Invitation struct {
	ServiceEndpoint string   `json:"serviceEndpoint,omitempty" faker:"url"`
	RecipientKeys   []string `json:"recipientKeys,omitempty" faker:"-"`
	ID              string   `json:"@id,omitempty" faker:"uuid_hyphenated"`
	Label           string   `json:"label,omitempty" faker:"first_name"`
	Type            string   `json:"@type,omitempty" faker:"-"`
}
