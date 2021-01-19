package findy

import (
	"context"
	"encoding/json"

	"github.com/findy-network/findy-agent-api/grpc/agency"
	"github.com/findy-network/findy-agent-vault/agency/model"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/findy-network/findy-grpc/agency/client"
	"github.com/findy-network/findy-grpc/agency/client/async"
	"github.com/findy-network/findy-grpc/jwt"
	"github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/lainio/err2"
	"google.golang.org/grpc"
)

const (
	agencyHost = "localhost"
	agencyPort = 50051
)

type Agency struct {
	vault   model.Listener
	ctx     context.Context
	tlsPath string
	options []grpc.DialOption
}

func (f *Agency) userCmdConn(a *model.Agent) client.Conn {
	config := client.BuildClientConnBase(f.tlsPath, agencyHost, agencyPort, f.options)
	return client.TryAuthOpen(a.RawJWT, config)
}

func (f *Agency) userCmdPw(a *model.Agent, connectionID string) *async.Pairwise {
	return async.NewPairwise(f.userCmdConn(a), connectionID)
}

func (f *Agency) Init(listener model.Listener, agents []*model.Agent, config *utils.Configuration) {
	jwt.SetJWTSecret(config.JWTKey)

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
	conn := f.userCmdConn(a)
	defer conn.Close()

	cmd := agency.NewAgentClient(conn)
	id = uuid.New().String()

	res, err := cmd.CreateInvitation(context.Background(), &agency.InvitationBase{Label: a.Label, Id: id})
	err2.Check(err)

	invitation = res.JsonStr

	return
}

func (f *Agency) Connect(a *model.Agent, strInvitation string) (id string, err error) {
	defer err2.Return(&err) // TODO: do not leak internal errors to client

	inv := model.Invitation{}
	err2.Check(json.Unmarshal([]byte(strInvitation), &inv))

	connect := f.userCmdPw(a, "")
	defer connect.Close()

	connect.Label = a.Label
	protocolID, err := connect.Connection(context.Background(), strInvitation)
	err2.Check(err)

	return protocolID.Id, err
}

func (f *Agency) SendMessage(a *model.Agent, connectionID, message string) (id string, err error) {
	defer err2.Return(&err) // TODO: do not leak internal errors to client

	pairwise := f.userCmdPw(a, connectionID)
	defer pairwise.Close()

	protocolID, err := pairwise.BasicMessage(context.Background(), message)
	err2.Check(err)

	return protocolID.Id, err
}

func (f *Agency) resume(
	a *model.Agent,
	job *model.JobInfo,
	accept bool,
	protocol agency.Protocol_Type,
	update func(*model.JobInfo, *int64, *int64, *int64),
) (err error) {
	defer err2.Return(&err) // TODO: do not leak internal errors to client

	conn := f.userCmdConn(a)
	defer conn.Close()

	state := agency.ProtocolState_NACK
	if accept {
		state = agency.ProtocolState_ACK
	}

	ctx := context.Background()
	didComm := agency.NewDIDCommClient(conn)
	_, err = didComm.Resume(ctx, &agency.ProtocolState{
		ProtocolId: &agency.ProtocolID{
			TypeId: protocol,
			Role:   agency.Protocol_RESUME,
			Id:     job.JobID,
		},
		State: state,
	})
	err2.Check(err)

	now := utils.CurrentTimeMs()
	update(job, &now, nil, nil)

	return err
}

func (f *Agency) ResumeCredentialOffer(a *model.Agent, job *model.JobInfo, accept bool) (err error) {
	return f.resume(a, job, accept, agency.Protocol_ISSUE, f.vault.UpdateCredential)
}

func (f *Agency) ResumeProofRequest(a *model.Agent, job *model.JobInfo, accept bool) (err error) {
	return f.resume(a, job, accept, agency.Protocol_PROOF, f.vault.UpdateProof)
}
