package findy

import (
	"context"
	"encoding/json"

	"github.com/findy-network/findy-agent-api/grpc/agency"
	"github.com/findy-network/findy-agent-vault/agency/model"
	"github.com/findy-network/findy-grpc/agency/client"
	"github.com/findy-network/findy-grpc/agency/client/async"
	"github.com/golang/glog"
	"github.com/lainio/err2"
)

const (
	agencyHost = "localhost"
	agencyPort = 50051
)

type Agency struct {
	vault model.Listener
	ctx   context.Context
}

func userCmdClient(a *model.Agent) *async.Pairwise {
	config := client.BuildClientConnBase("", agencyHost, agencyPort, nil)
	return async.NewPairwise(client.TryAuthOpen(a.RawJWT, config), "")
}

func userListenClient(a *model.Agent) client.Conn {
	config := client.BuildClientConnBase("", agencyHost, agencyPort, nil)
	return client.TryOpen(a.AgentID, config)
}

func (f *Agency) listenAgent(a *model.Agent) (err error) {
	defer err2.Return(&err)
	// TODO: cancellation

	conn := userListenClient(a)

	ch, err := conn.Listen(f.ctx, &agency.ClientID{Id: a.TenantID})
	err2.Check(err)

	go func() {
		for {
			select {
			case status, ok := <-ch:
				if !ok {
					glog.V(2).Infoln("closed from server")
					break
				}
				glog.V(5).Infoln("listen status:",
					status.Notification.TypeId,
					status.Notification.Role,
					status.Notification.ProtocolId)
			}
		}
	}()

	return
}

func (f *Agency) Init(listener model.Listener, agents []*model.Agent) {
	f.ctx = context.Background()
	f.vault = listener
	// TODO: create JWT on demand
	// TODO: get all agents from db and start listening
	// TODO: start listening when onboarding
	f.listenAdminHook()
	for _, a := range agents {
		f.listenAgent(a)
	}
}

func (f *Agency) Invite(a *model.Agent) (invitation, id string, err error) {
	return
}

func (f *Agency) Connect(a *model.Agent, strInvitation string) (id string, err error) {
	defer err2.Return(&err) // TODO: do not leak internal errors to client

	inv := model.Invitation{}
	err2.Check(json.Unmarshal([]byte(strInvitation), &inv))

	client := userCmdClient(a)
	defer client.Close()

	client.Label = a.Label
	protocolID, err := client.Connection(context.Background(), strInvitation)
	err2.Check(err)

	return protocolID.Id, err
}

func (f *Agency) SendMessage(a *model.Agent, connectionID, message string) (id string, err error) {
	return
}

func (f *Agency) ResumeCredentialOffer(a *model.Agent, id string, accept bool) (err error) {
	return
}

func (f *Agency) ResumeProofRequest(a *model.Agent, id string, accept bool) (err error) {
	return
}
