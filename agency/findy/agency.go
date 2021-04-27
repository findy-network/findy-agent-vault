package findy

import (
	"context"
	"encoding/json"

	"github.com/findy-network/findy-agent-vault/agency/model"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/findy-network/findy-common-go/agency/client"
	agency "github.com/findy-network/findy-common-go/grpc/agency/v1"
	"github.com/findy-network/findy-common-go/jwt"
	didexchange "github.com/findy-network/findy-common-go/std/didexchange/invitation"
	"github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/lainio/err2"
	"google.golang.org/grpc"
)

type Agency struct {
	currentTimeMs func() int64
	vault         model.Listener
	archiver      model.Archiver

	agencyHost string
	agencyPort int
	tlsPath    string
	options    []grpc.DialOption

	ctx  context.Context
	conn client.Conn
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

func (f *Agency) Invite(a *model.Agent) (invitation, id string, err error) {
	cmd := agency.NewAgentServiceClient(f.conn)
	id = uuid.New().String()

	res, err := cmd.CreateInvitation(
		f.ctx,
		&agency.InvitationBase{Label: a.Label, ID: id},
		callOptions(a.RawJWT)...,
	)
	err2.Check(err)

	invitation = res.JsonStr

	return
}

func (f *Agency) Connect(a *model.Agent, strInvitation string) (id string, err error) {
	defer err2.Return(&err) // TODO: do not leak internal errors to client

	inv := didexchange.Invitation{}
	err2.Check(json.Unmarshal([]byte(strInvitation), &inv))

	cmd := f.userSyncClient(a, "")

	cmd.Label = a.Label
	protocolID, err := cmd.Connection(f.ctx, strInvitation)
	err2.Check(err)

	return protocolID.ID, err
}

func (f *Agency) SendMessage(a *model.Agent, connectionID, message string) (id string, err error) {
	defer err2.Return(&err) // TODO: do not leak internal errors to client

	cmd := f.userSyncClient(a, connectionID)

	protocolID, err := cmd.BasicMessage(f.ctx, message)
	err2.Check(err)

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

	_, err = cmd.Resume(f.ctx, job.JobID, protocol, state)
	err2.Check(err)

	return
}

func (f *Agency) ResumeCredentialOffer(a *model.Agent, job *model.JobInfo, accept bool) (err error) {
	defer err2.Return(&err)
	err2.Check(f.resume(a, job, accept, agency.Protocol_ISSUE_CREDENTIAL))

	now := f.currentTimeMs()
	return f.vault.UpdateCredential(job, &model.CredentialUpdate{ApprovedMs: &now})
}

func (f *Agency) ResumeProofRequest(a *model.Agent, job *model.JobInfo, accept bool) (err error) {
	defer err2.Return(&err)
	err2.Check(f.resume(a, job, accept, agency.Protocol_PRESENT_PROOF))

	now := f.currentTimeMs()
	return f.vault.UpdateProof(job, &model.ProofUpdate{ApprovedMs: &now})
}
