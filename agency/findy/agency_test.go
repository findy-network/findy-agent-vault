package findy

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/findy-network/findy-agent-api/grpc/agency"
	"github.com/findy-network/findy-agent-api/grpc/ops"
	"github.com/findy-network/findy-agent-vault/agency/model"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/findy-network/findy-common-go/jwt"
	"github.com/findy-network/findy-common-go/rpc"
	"github.com/golang/glog"
	"github.com/lainio/err2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

const (
	testInvitation = `{"serviceEndpoint":` +
		`"http://findy-agent.op-ai.fi/a2a/Xmjk7cFr8TT2j5kWLWyhDB/Xmjk7cFr8TT2j5kWLWyhDB/GqmnSTxevze48yio5m2fUE",` +
		`"recipientKeys":["Hmk4756ry7fqBCKPf634SRvaM3xss1QBhoFC1uAbwkVL"],"@id":"d679e4c6-b8db-4c39-99ca-783034b51bd4"` +
		`,"label":"findy-issuer","@type":"did:sov:BzCbsNYhMrjHiqZDTUASHg;spec/connections/1.0/invitation"}`
	testID = "d679e4c6-b8db-4c39-99ca-783034b51bd4"
)

type mockServer struct {
	agency.UnimplementedDIDCommServer
	ops.UnimplementedAgencyServer
	agency.UnimplementedAgentServer
	hookID    string
	clientIDs []string
}

func (*mockServer) Run(*agency.Protocol, agency.DIDComm_RunServer) error {
	return status.Errorf(codes.Unimplemented, "method Run not implemented")
}
func (*mockServer) Start(context.Context, *agency.Protocol) (*agency.ProtocolID, error) {
	return &agency.ProtocolID{Id: testID}, nil
}
func (*mockServer) Status(context.Context, *agency.ProtocolID) (*agency.ProtocolStatus, error) {
	return &agency.ProtocolStatus{}, nil
}
func (*mockServer) Resume(context.Context, *agency.ProtocolState) (*agency.ProtocolID, error) {
	return &agency.ProtocolID{Id: testID}, nil
}
func (*mockServer) Release(context.Context, *agency.ProtocolID) (*agency.ProtocolID, error) {
	return &agency.ProtocolID{Id: testID}, nil
}

func (m *mockServer) PSMHook(dataHook *ops.DataHook, server ops.Agency_PSMHookServer) error {
	m.hookID = dataHook.Id
	return nil
}
func (*mockServer) Onboard(context.Context, *ops.Onboarding) (*ops.OnboardResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Onboard not implemented")
}

func (m *mockServer) Listen(id *agency.ClientID, server agency.Agent_ListenServer) error {
	m.clientIDs = append(m.clientIDs, id.Id)
	return nil
}
func (*mockServer) Give(context.Context, *agency.Answer) (*agency.ClientID, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Give not implemented")
}
func (*mockServer) CreateInvitation(context.Context, *agency.InvitationBase) (*agency.Invitation, error) {
	return &agency.Invitation{JsonStr: testInvitation}, nil
}

func dialer() func(context.Context, string) (net.Conn, error) {
	const bufSize = 1024 * 1024

	listener := bufconn.Listen(bufSize)
	// TODO:
	pki := rpc.LoadPKI("../../.github/workflows/cert")
	glog.V(1).Infof("starting gRPC server with\ncrt:\t%s\nkey:\t%s\nclient:\t%s",
		pki.Server.CertFile, pki.Server.KeyFile, pki.Client.CertFile)

	s, lis, err := rpc.PrepareServe(&rpc.ServerCfg{
		Port:    50051,
		PKI:     pki,
		TestLis: listener,
		Register: func(s *grpc.Server) error {
			agency.RegisterDIDCommServer(s, mockAgencyServer)
			ops.RegisterAgencyServer(s, mockAgencyServer)
			agency.RegisterAgentServer(s, mockAgencyServer)
			glog.V(10).Infoln("GRPC registration all done")
			return nil
		},
	})
	if err != nil {
		panic(fmt.Errorf("unable to register mock server %v", err))
	}

	go func() {
		defer err2.Catch(func(err error) {
			log.Fatal(err)
		})
		err2.Check(s.Serve(lis))
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

type mockListener struct {
	credTS  int64
	proofTS int64
}

func (m *mockListener) AddConnection(job *model.JobInfo, connection *model.Connection) error {
	panic("Not implemented")
}
func (m *mockListener) AddMessage(job *model.JobInfo, message *model.Message) error {
	panic("Not implemented")
}
func (m *mockListener) UpdateMessage(job *model.JobInfo, update *model.MessageUpdate) error {
	panic("Not implemented")
}

func (m *mockListener) AddCredential(job *model.JobInfo, credential *model.Credential) error {
	panic("Not implemented")
}
func (m *mockListener) UpdateCredential(job *model.JobInfo, update *model.CredentialUpdate) error {
	m.credTS = *update.ApprovedMs
	return nil
}

func (m *mockListener) AddProof(job *model.JobInfo, proof *model.Proof) error {
	panic("Not implemented")
}

func (m *mockListener) UpdateProof(job *model.JobInfo, update *model.ProofUpdate) error {
	m.proofTS = *update.ApprovedMs
	return nil
}

func (m *mockListener) FailJob(job *model.JobInfo) error {
	panic("Not implemented")
}

var (
	tlsPath     = "../../.github/workflows/cert"
	dialOptions = []grpc.DialOption{grpc.WithContextDialer(dialer())}
	findy       = &Agency{
		vault:   &mockListener{},
		tlsPath: tlsPath,
		options: dialOptions,
	}
	agent *model.Agent

	mockAgencyServer = &mockServer{clientIDs: make([]string, 0)}
)

func setup() {
	utils.SetLogDefaults()
	findy.Init(
		&mockListener{},
		[]*model.Agent{},
		&mockArchiver{},
		&utils.Configuration{JWTKey: "mySuperSecretKeyLol", AgencyCertPath: tlsPath},
	)
	agent = &model.Agent{RawJWT: jwt.BuildJWT("test-user")}
}

func teardown() {
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestInit(t *testing.T) {
	const testClientID = "test"
	testAgency := &Agency{options: dialOptions}
	testAgency.Init(
		&mockListener{},
		[]*model.Agent{{AgentID: testClientID, TenantID: testClientID}},
		&mockArchiver{},
		&utils.Configuration{JWTKey: "mySuperSecretKeyLol", AgencyCertPath: tlsPath, AgencyMainSubscriber: true},
	)
	// Wait for a while that calls complete
	time.Sleep(time.Millisecond * 100)
	if mockAgencyServer.hookID == "" {
		t.Errorf("psm hook registration failed")
	}
	found := false
	for _, clientID := range mockAgencyServer.clientIDs {
		if clientID == testClientID {
			found = true
		}
	}
	if !found {
		t.Errorf("client listener registration failed")
	}
}

func TestInvite(t *testing.T) {
	invitation, id, err := findy.Invite(agent)
	if err != nil {
		t.Errorf("Encountered error on invite %v", err)
	}
	if id == "" {
		t.Errorf("Received empty job id ")
	}
	if invitation != testInvitation {
		t.Errorf("Mismatch with invitation expecting %v, got %v", testInvitation, invitation)
	}
}

func TestConnect(t *testing.T) {
	id, err := findy.Connect(agent, testInvitation)
	if err != nil {
		t.Errorf("Encountered error on connect %v", err)
	}
	if id != testID {
		t.Errorf("Mismatch with id expecting %v, got %v", testID, id)
	}
}

func TestSendMessage(t *testing.T) {
	id, err := findy.SendMessage(agent, "id", "message")
	if err != nil {
		t.Errorf("Encountered error on connect %v", err)
	}
	if id != testID {
		t.Errorf("Mismatch with id expecting %v, got %v", testID, id)
	}
}

func TestResumeCredentialOffer(t *testing.T) {
	err := findy.ResumeCredentialOffer(agent, &model.JobInfo{}, true)
	if err != nil {
		t.Errorf("Encountered error on resume credential offer %v", err)
	}
	vault := findy.vault.(*mockListener)
	if vault.credTS != utils.CurrentTimeMs() {
		t.Errorf("Expected valid credential timestamp %v", vault.credTS)
	}
}

func TestResumeProofRequest(t *testing.T) {
	err := findy.ResumeProofRequest(agent, &model.JobInfo{}, true)
	if err != nil {
		t.Errorf("Encountered error on resume proof request %v", err)
	}
	vault := findy.vault.(*mockListener)
	if vault.proofTS != utils.CurrentTimeMs() {
		t.Errorf("Expected valid proof timestamp %v", vault.proofTS)
	}
}
