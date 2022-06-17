package findy

import (
	"context"
	"encoding/json"

	"github.com/findy-network/findy-agent-vault/agency/model"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/findy-network/findy-common-go/agency/client"
	agency "github.com/findy-network/findy-common-go/grpc/agency/v1"
	"github.com/findy-network/findy-common-go/jwt"
	"github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
	"google.golang.org/grpc"
)

type Agency struct {
	currentTimeMs func() int64
	vault         model.Listener
	archiver      model.Archiver

	agencyHost    string
	agencyPort    int
	agencyAdminID string
	tlsPath       string
	options       []grpc.DialOption

	ctx  context.Context
	conn client.Conn

	userAsyncClient func(a *model.Agent) clientConn
}

func (f *Agency) Init(
	listener model.Listener,
	agents []*model.Agent,
	archiver model.Archiver,
	config *utils.Configuration,
) {
	jwt.SetJWTSecret(config.JWTKey)
	f.currentTimeMs = utils.CurrentTimeMs

	f.agencyHost = config.AgencyHost
	f.agencyPort = config.AgencyPort
	f.agencyAdminID = config.AgencyAdminID
	f.tlsPath = config.AgencyCertPath

	f.ctx = context.Background()
	// open connection without JWT token
	// instead, token is set on each call
	f.conn = client.TryAuthOpen(
		"",
		client.BuildClientConnBase(f.tlsPath, f.agencyHost, f.agencyPort, f.options),
	)

	f.vault = listener
	f.archiver = archiver
	f.userAsyncClient = f.getUserAsyncClient

	if config.AgencyMainSubscriber {
		err := f.listenAdminHook()
		if err != nil {
			panic(err)
		}
	} else {
		glog.Warningln("DEV mode: Skipping subscribing to PSM hook.")
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

func (f *Agency) Invite(a *model.Agent) (data *model.InvitationData, err error) {
	cmd := agency.NewAgentServiceClient(f.conn)
	id := uuid.New().String()

	res, err := cmd.CreateInvitation(
		f.ctx,
		&agency.InvitationBase{Label: a.Label, ID: id},
		callOptions(a.RawJWT)...,
	)
	try.To(err)

	data = &model.InvitationData{}

	try.To(json.Unmarshal([]byte(res.GetJSON()), &data.Data))

	data.Raw = res.GetURL()
	data.ID = id

	return
}

func (f *Agency) Connect(a *model.Agent, strInvitation string) (id string, err error) {
	defer err2.Return(&err) // TODO: do not leak internal errors to client

	cmd := f.userSyncClient(a, "")

	cmd.Label = a.Label
	protocolID := try.To1(cmd.Connection(f.ctx, strInvitation))

	return protocolID.ID, err
}

func (f *Agency) SendMessage(a *model.Agent, connectionID, message string) (id string, err error) {
	defer err2.Return(&err) // TODO: do not leak internal errors to client

	cmd := f.userSyncClient(a, connectionID)

	protocolID := try.To1(cmd.BasicMessage(f.ctx, message))

	return protocolID.ID, err
}

func (f *Agency) resume(
	a *model.Agent,
	job *model.JobInfo,
	accept bool,
	protocol agency.Protocol_Type,
) (err error) {
	defer err2.Return(&err) // TODO: do not leak internal errors to client

	cmd := f.userSyncClient(a, job.ConnectionID)
	state := agency.ProtocolState_NACK
	if accept {
		state = agency.ProtocolState_ACK
	}

	try.To1(cmd.Resume(f.ctx, job.JobID, protocol, state))

	return
}

func (f *Agency) ResumeCredentialOffer(a *model.Agent, job *model.JobInfo, accept bool) (err error) {
	defer err2.Return(&err)
	try.To(f.resume(a, job, accept, agency.Protocol_ISSUE_CREDENTIAL))

	now := f.currentTimeMs()
	return f.vault.UpdateCredential(job, &model.CredentialUpdate{ApprovedMs: &now})
}

func (f *Agency) ResumeProofRequest(a *model.Agent, job *model.JobInfo, accept bool) (err error) {
	defer err2.Return(&err)
	try.To(f.resume(a, job, accept, agency.Protocol_PRESENT_PROOF))

	now := f.currentTimeMs()
	return f.vault.UpdateProof(job, &model.ProofUpdate{ApprovedMs: &now})
}
