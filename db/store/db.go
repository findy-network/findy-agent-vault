package store

import (
	"context"

	"github.com/findy-network/findy-agent-vault/db/model"
	graph "github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/paginator"
	"github.com/findy-network/findy-grpc/jwt"
)

func GetAgent(ctx context.Context, db DB) (*model.Agent, error) {
	token, err := jwt.TokenFromContext(ctx, "user")
	if err != nil {
		return nil, err
	}
	a := model.NewAgent(nil)
	a.AgentID = token.AgentID
	a.Label = token.Label
	a.RawJWT = &token.Raw
	return db.AddAgent(a)
}

type DB interface {
	GetListenerAgents(info *paginator.BatchInfo) (*model.Agents, error)
	Close()

	AddAgent(a *model.Agent) (*model.Agent, error)
	GetAgent(id, agentID *string) (*model.Agent, error)

	AddConnection(c *model.Connection) (*model.Connection, error)
	GetConnection(id, tenantID string) (*model.Connection, error)
	GetConnections(info *paginator.BatchInfo, tenantID string) (connections *model.Connections, err error)
	GetConnectionCount(tenantID string) (int, error)

	AddCredential(c *model.Credential) (*model.Credential, error)
	UpdateCredential(c *model.Credential) (*model.Credential, error)
	GetCredential(id, tenantID string) (*model.Credential, error)
	GetCredentials(info *paginator.BatchInfo, tenantID string, connectionID *string) (connections *model.Credentials, err error)
	GetCredentialCount(tenantID string, connectionID *string) (int, error)
	GetConnectionForCredential(id, tenantID string) (*model.Connection, error)

	AddProof(p *model.Proof) (*model.Proof, error)
	UpdateProof(p *model.Proof) (*model.Proof, error)
	GetProof(id, tenantID string) (*model.Proof, error)
	GetProofs(info *paginator.BatchInfo, tenantID string, connectionID *string) (connections *model.Proofs, err error)
	GetProofCount(tenantID string, connectionID *string) (int, error)
	GetConnectionForProof(id, tenantID string) (*model.Connection, error)

	AddMessage(m *model.Message) (*model.Message, error)
	UpdateMessage(m *model.Message) (*model.Message, error)
	GetMessage(id, tenantID string) (*model.Message, error)
	GetMessages(info *paginator.BatchInfo, tenantID string, connectionID *string) (connections *model.Messages, err error)
	GetMessageCount(tenantID string, connectionID *string) (int, error)
	GetConnectionForMessage(id, tenantID string) (*model.Connection, error)

	AddEvent(e *model.Event) (*model.Event, error)
	MarkEventRead(id, tenantID string) (*model.Event, error)
	GetEvent(id, tenantID string) (*model.Event, error)
	GetEvents(info *paginator.BatchInfo, tenantID string, connectionID *string) (connections *model.Events, err error)
	GetEventCount(tenantID string, connectionID *string) (int, error)
	GetConnectionForEvent(id, tenantID string) (*model.Connection, error)
	GetJobForEvent(id, tenantID string) (*model.Job, error)
	GetJobOutput(id, tenantID string, protocolType graph.ProtocolType) (*model.JobOutput, error)

	AddJob(j *model.Job) (*model.Job, error)
	UpdateJob(j *model.Job) (*model.Job, error)
	GetJob(id, tenantID string) (*model.Job, error)
	GetJobs(info *paginator.BatchInfo, tenantID string, connectionID *string, completed *bool) (connections *model.Jobs, err error)
	GetJobCount(tenantID string, connectionID *string, completed *bool) (int, error)
	GetConnectionForJob(id, tenantID string) (*model.Connection, error)
}
