package model

import (
	"github.com/findy-network/findy-agent-vault/utils"
)

type JobInfo struct {
	TenantID     string
	JobID        string
	ConnectionID string
}

type Agent struct {
	Label    string
	RawJWT   string
	TenantID string
	AgentID  string
}

type Agency interface {
	Init(l Listener, agents []*Agent, config *utils.Configuration)
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
