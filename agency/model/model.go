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

type InvitationData struct {
	Raw string
	ID  string
}

type Agency interface {
	Init(l Listener, agents []*Agent, archiver Archiver, config *utils.Configuration)
	AddAgent(agent *Agent) error

	Invite(a *Agent) (*InvitationData, error)
	Connect(a *Agent, invitation string) (string, error)
	SendMessage(a *Agent, connectionID, message string) (string, error)

	ResumeCredentialOffer(a *Agent, job *JobInfo, accept bool) error
	ResumeProofRequest(a *Agent, job *JobInfo, accept bool) error
}
