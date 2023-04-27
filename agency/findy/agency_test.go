package findy

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/findy-network/findy-agent-vault/agency/model"
	dbModel "github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/utils"
	agency "github.com/findy-network/findy-common-go/grpc/agency/v1"
	ops "github.com/findy-network/findy-common-go/grpc/ops/v1"
	"github.com/findy-network/findy-common-go/jwt"
	"github.com/findy-network/findy-common-go/rpc"
	"github.com/golang/glog"
	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

const (
	testInvitation = `{"serviceEndpoint":` +
		`"http://url",` +
		`"recipientKeys":["Hmk4756ry7fqBCKPf634SRvaM3xss1QBhoFC1uAbwkVL"],"@id":"d679e4c6-b8db-4c39-99ca-783034b51bd4"` +
		`,"label":"findy-issuer","@type":"did:sov:BzCbsNYhMrjHiqZDTUASHg;spec/connections/1.0/invitation"}`
	testID = "d679e4c6-b8db-4c39-99ca-783034b51bd4"
)

var (
	testInvitationURL = fmt.Sprintf("didcomm://aries_connection_invitation?c_i=%s", base64.StdEncoding.EncodeToString([]byte(testInvitation)))
)

type mockServer struct {
	agency.UnimplementedProtocolServiceServer
	ops.UnimplementedAgencyServiceServer
	agency.UnimplementedAgentServiceServer
	hookID    string
	clientIDs []string
}

func (*mockServer) Run(*agency.Protocol, agency.ProtocolService_RunServer) error {
	return status.Errorf(codes.Unimplemented, "method Run not implemented")
}
func (*mockServer) Start(context.Context, *agency.Protocol) (*agency.ProtocolID, error) {
	return &agency.ProtocolID{ID: testID}, nil
}
func (*mockServer) Status(context.Context, *agency.ProtocolID) (*agency.ProtocolStatus, error) {
	return &agency.ProtocolStatus{}, nil
}
func (*mockServer) Resume(context.Context, *agency.ProtocolState) (*agency.ProtocolID, error) {
	return &agency.ProtocolID{ID: testID}, nil
}
func (*mockServer) Release(context.Context, *agency.ProtocolID) (*agency.ProtocolID, error) {
	return &agency.ProtocolID{ID: testID}, nil
}

func (m *mockServer) PSMHook(dataHook *ops.DataHook, _ ops.AgencyService_PSMHookServer) error {
	m.hookID = dataHook.ID
	return nil
}
func (*mockServer) Onboard(context.Context, *ops.Onboarding) (*ops.OnboardResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Onboard not implemented")
}

func (m *mockServer) Listen(id *agency.ClientID, _ agency.AgentService_ListenServer) error {
	m.clientIDs = append(m.clientIDs, id.ID)
	return nil
}
func (*mockServer) Give(context.Context, *agency.Answer) (*agency.ClientID, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Give not implemented")
}
func (*mockServer) CreateInvitation(context.Context, *agency.InvitationBase) (*agency.Invitation, error) {
	return &agency.Invitation{JSON: testInvitation, URL: testInvitationURL}, nil
}

func dialer(insecure bool) func(context.Context, string) (net.Conn, error) {
	const bufSize = 1024 * 1024

	listener := bufconn.Listen(bufSize)
	var pki *rpc.PKI
	if !insecure {
		// TODO:
		pki = rpc.LoadPKI("../../scripts/test-cert")
		glog.V(1).Infof("starting gRPC server with\ncrt:\t%s\nkey:\t%s\nclient:\t%s",
			pki.Server.CertFile, pki.Server.KeyFile, pki.Client.CertFile)
	}

	s, lis, err := rpc.PrepareServe(&rpc.ServerCfg{
		Port:    50051,
		PKI:     pki,
		TestLis: listener,
		Register: func(s *grpc.Server) error {
			agency.RegisterProtocolServiceServer(s, mockAgencyServer)
			ops.RegisterAgencyServiceServer(s, mockAgencyServer)
			agency.RegisterAgentServiceServer(s, mockAgencyServer)
			glog.V(10).Infoln("GRPC registration all done")
			return nil
		},
	})
	if err != nil {
		panic(fmt.Errorf("unable to register mock server %w", err))
	}

	go func() {
		defer err2.Catch(func(err error) {
			log.Fatal(err)
		})
		try.To(s.Serve(lis))
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

type mockListener struct {
	credTS  int64
	proofTS int64
}

func (m *mockListener) AddConnection(_ *model.JobInfo, _ *model.Connection) error {
	panic("Not implemented")
}
func (m *mockListener) AddMessage(_ *model.JobInfo, _ *model.Message) error {
	panic("Not implemented")
}
func (m *mockListener) UpdateMessage(_ *model.JobInfo, _ *model.MessageUpdate) error {
	panic("Not implemented")
}

func (m *mockListener) AddCredential(_ *model.JobInfo, _ *model.Credential) (*dbModel.Job, error) {
	panic("Not implemented")
}
func (m *mockListener) UpdateCredential(_ *model.JobInfo, _ *model.Credential, update *model.CredentialUpdate) error {
	m.credTS = *update.ApprovedMs
	return nil
}

func (m *mockListener) AddProof(_ *model.JobInfo, _ *model.Proof) (*dbModel.Job, error) {
	panic("Not implemented")
}

func (m *mockListener) UpdateProof(_ *model.JobInfo, _ *model.Proof, update *model.ProofUpdate) error {
	m.proofTS = *update.ApprovedMs
	return nil
}

func (m *mockListener) FailJob(_ *model.JobInfo) error {
	panic("Not implemented")
}

var (
	tlsPath             = "../../scripts/test-cert"
	dialOptions         = []grpc.DialOption{grpc.WithContextDialer(dialer(false))}
	insecureDialOptions = []grpc.DialOption{grpc.WithContextDialer(dialer(true))}
	findy               = &Agency{
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
		&utils.Configuration{JWTKey: "mySuperSecretKeyLol", AgencyCertPath: tlsPath, AgencyAdminID: "admin-id", AgencyMainSubscriber: true},
	)
	// Wait for a while that calls complete
	time.Sleep(time.Millisecond * 100)
	if mockAgencyServer.hookID == "" {
		t.Errorf("psm hook registration failed")
	}
	mockAgencyServer.hookID = ""
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

func TestInitInsecure(t *testing.T) {
	const testClientID = "test"
	testAgency := &Agency{options: insecureDialOptions}
	testAgency.Init(
		&mockListener{},
		[]*model.Agent{{AgentID: testClientID, TenantID: testClientID}},
		&mockArchiver{},
		&utils.Configuration{
			JWTKey:               "mySuperSecretKeyLol",
			AgencyCertPath:       "",
			AgencyAdminID:        "admin-id",
			AgencyMainSubscriber: true,
			AgencyInsecure:       true,
		},
	)
	// Wait for a while that calls complete
	time.Sleep(time.Millisecond * 100)
	if mockAgencyServer.hookID == "" {
		t.Errorf("psm hook registration failed")
	}
	mockAgencyServer.hookID = ""
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
	data, err := findy.Invite(agent)
	if err != nil {
		t.Errorf("Encountered error on invite %v", err)
	}
	if data.ID == "" {
		t.Errorf("Received empty job id ")
	}
	if data.Raw != testInvitationURL {
		t.Errorf("Mismatch with invitation expecting %v, got %v", testInvitationURL, data.Raw)
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
	if vault, ok := findy.vault.(*mockListener); ok {
		if vault.credTS != utils.CurrentTimeMs() {
			t.Errorf("Expected valid credential timestamp %v", vault.credTS)
		}
	}
}

func TestResumeProofRequest(t *testing.T) {
	err := findy.ResumeProofRequest(agent, &model.JobInfo{}, true)
	if err != nil {
		t.Errorf("Encountered error on resume proof request %v", err)
	}
	if vault, ok := findy.vault.(*mockListener); ok {
		if vault.proofTS != utils.CurrentTimeMs() {
			t.Errorf("Expected valid proof timestamp %v", vault.proofTS)
		}
	}
}
