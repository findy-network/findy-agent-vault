// +build findy_grpc

package findy

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/findy-network/findy-agent-api/grpc/agency"
	"github.com/findy-network/findy-agent-vault/agency/model"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/findy-network/findy-grpc/agency/client"
	"github.com/findy-network/findy-grpc/agency/client/async"
	"github.com/golang/glog"
	"github.com/google/uuid"
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

func userCmdConn(a *model.Agent) client.Conn {
	config := client.BuildClientConnBase("", agencyHost, agencyPort, nil)
	return client.TryAuthOpen(a.RawJWT, config)
}

func userCmdPw(a *model.Agent, connectionID string) *async.Pairwise {
	return async.NewPairwise(userCmdConn(a), connectionID)
}

func (f *Agency) Init(listener model.Listener, agents []*model.Agent) {
	f.ctx = context.Background()
	f.vault = listener
	// TODO: release protocol when saved
	err := f.listenAdminHook()
	if err != nil {
		panic(err)
	}
	for _, a := range agents {
		err := f.listenAgent(a)
		if err != nil {
			glog.Error(err)
		}
	}
}

func (f *Agency) AddAgent(agent *model.Agent) error {
	return f.listenAgent(agent)
}

func (f *Agency) Invite(a *model.Agent) (invitation, id string, err error) {
	conn := userCmdConn(a)
	defer conn.Close()

	cmd := agency.NewAgentClient(conn)
	id = uuid.New().String()

	res, err := cmd.CreateInvitation(f.ctx, &agency.InvitationBase{Label: a.Label, Id: id})
	err2.Check(err)

	invitation = res.JsonStr

	return
}

func (f *Agency) Connect(a *model.Agent, strInvitation string) (id string, err error) {
	defer err2.Return(&err) // TODO: do not leak internal errors to client

	inv := model.Invitation{}
	err2.Check(json.Unmarshal([]byte(strInvitation), &inv))

	connect := userCmdPw(a, "")
	defer connect.Close()

	connect.Label = a.Label
	protocolID, err := connect.Connection(context.Background(), strInvitation)
	fmt.Println(err)
	err2.Check(err)

	return protocolID.Id, err
}

func (f *Agency) SendMessage(a *model.Agent, connectionID, message string) (id string, err error) {
	defer err2.Return(&err) // TODO: do not leak internal errors to client

	pairwise := userCmdPw(a, connectionID)
	defer pairwise.Close()

	protocolID, err := pairwise.BasicMessage(context.Background(), message)
	err2.Check(err)

	return protocolID.Id, err
}

func (f *Agency) resume(
	conn client.Conn,
	id string,
	protocol agency.Protocol_Type,
	state agency.ProtocolState_State,
) (*agency.ProtocolID, error) {
	ctx := context.Background()
	didComm := agency.NewDIDCommClient(conn)
	return didComm.Resume(ctx, &agency.ProtocolState{
		ProtocolId: &agency.ProtocolID{
			TypeId: protocol,
			Role:   agency.Protocol_RESUME,
			Id:     id,
		},
		State: state,
	})
}

func (f *Agency) ResumeCredentialOffer(a *model.Agent, job *model.JobInfo, accept bool) (err error) {
	defer err2.Return(&err) // TODO: do not leak internal errors to client

	conn := userCmdConn(a)
	defer conn.Close()

	state := agency.ProtocolState_NACK
	if accept {
		state = agency.ProtocolState_ACK
	}

	_, err = f.resume(conn, job.JobID, agency.Protocol_ISSUE, state)
	err2.Check(err)

	now := utils.CurrentTimeMs()
	f.vault.UpdateCredential(
		job,
		&now,
		nil,
		nil,
	)

	return
}

func (f *Agency) ResumeProofRequest(a *model.Agent, job *model.JobInfo, accept bool) (err error) {
	defer err2.Return(&err) // TODO: do not leak internal errors to client

	conn := userCmdConn(a)
	defer conn.Close()

	state := agency.ProtocolState_NACK
	if accept {
		state = agency.ProtocolState_ACK
	}

	_, err = f.resume(conn, job.JobID, agency.Protocol_PROOF, state)
	err2.Check(err)

	now := utils.CurrentTimeMs()
	f.vault.UpdateProof(
		job,
		&now,
		nil,
		nil,
	)
	return
}
