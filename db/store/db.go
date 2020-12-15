package store

import (
	"context"

	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/paginator"
	"github.com/findy-network/findy-agent-vault/utils"
)

func GetAgent(ctx context.Context, db DB) (*model.Agent, error) {
	token, err := utils.ParseToken(ctx)
	if err != nil {
		return nil, err
	}
	a := model.NewAgent(nil)
	a.AgentID = token.AgentID
	a.Label = token.Label
	return db.AddAgent(a)
}

type DB interface {
	Close()

	AddAgent(a *model.Agent) (*model.Agent, error)
	GetAgent(id, agentID *string) (*model.Agent, error)

	AddConnection(c *model.Connection) (*model.Connection, error)
	GetConnection(id string, tenantID string) (*model.Connection, error)
	GetConnections(info *paginator.BatchInfo, tenantID string) (connections *model.Connections, err error)
	GetConnectionCount(tenantID string) (int, error)

	AddCredential(c *model.Credential) (*model.Credential, error)
	UpdateCredential(c *model.Credential) (*model.Credential, error)
	GetCredential(id string, tenantID string) (*model.Credential, error)
	GetCredentials(info *paginator.BatchInfo, tenantID string, connectionID *string) (connections *model.Credentials, err error)
	GetCredentialCount(tenantID string, connectionID *string) (int, error)

	AddEvent(e *model.Event) (*model.Event, error)
	MarkEventRead(id, tenantID string) (*model.Event, error)
	GetEvent(id, tenantID string) (*model.Event, error)
	GetEvents(info *paginator.BatchInfo, tenantID string, connectionID *string) (connections *model.Events, err error)
	GetEventCount(tenantID string, connectionID *string) (int, error)

	AddMessage(c *model.Message) (*model.Message, error)
	UpdateMessage(c *model.Message) (*model.Message, error)
	GetMessage(id string, tenantID string) (*model.Message, error)
	GetMessages(info *paginator.BatchInfo, tenantID string, connectionID *string) (connections *model.Messages, err error)
	GetMessageCount(tenantID string, connectionID *string) (int, error)
}
