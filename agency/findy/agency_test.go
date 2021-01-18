package findy

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/findy-network/findy-agent-api/grpc/agency"
	"github.com/findy-network/findy-agent-vault/agency/model"
	graph "github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/findy-network/findy-grpc/jwt"
	"github.com/findy-network/findy-grpc/rpc"
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

type mockDIDCommServer struct {
	agency.UnimplementedDIDCommServer
}

func (*mockDIDCommServer) Run(*agency.Protocol, agency.DIDComm_RunServer) error {
	return status.Errorf(codes.Unimplemented, "method Run not implemented")
}
func (*mockDIDCommServer) Start(context.Context, *agency.Protocol) (*agency.ProtocolID, error) {
	return &agency.ProtocolID{Id: testID}, nil
}
func (*mockDIDCommServer) Status(context.Context, *agency.ProtocolID) (*agency.ProtocolStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Status not implemented")
}
func (*mockDIDCommServer) Resume(context.Context, *agency.ProtocolState) (*agency.ProtocolID, error) {
	return &agency.ProtocolID{Id: testID}, nil
}
func (*mockDIDCommServer) Release(context.Context, *agency.ProtocolID) (*agency.ProtocolID, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Release not implemented")
}

func dialer(t *testing.T) func(context.Context, string) (net.Conn, error) {
	const bufSize = 1024 * 1024

	listener := bufconn.Listen(bufSize)
	// TODO:
	pki := rpc.LoadPKI("../../../findy-grpc/cert")
	glog.V(1).Infof("starting gRPC server with\ncrt:\t%s\nkey:\t%s\nclient:\t%s",
		pki.Server.CertFile, pki.Server.KeyFile, pki.Client.CertFile)

	s, lis, err := rpc.PrepareServe(&rpc.ServerCfg{
		Port:    50051,
		PKI:     pki,
		TestLis: listener,
		Register: func(s *grpc.Server) error {
			agency.RegisterDIDCommServer(s, &mockDIDCommServer{})
			glog.V(10).Infoln("GRPC registration all done")
			return nil
		},
	})
	if err != nil {
		t.Errorf("Unable to register mock server %v", err)
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

func (m *mockListener) AddConnection(job *model.JobInfo, ourDID, theirDID, theirEndpoint, theirLabel string) {
	panic("Not implemented")
}
func (m *mockListener) AddMessage(job *model.JobInfo, message string, sentByMe bool) {
	panic("Not implemented")
}
func (m *mockListener) UpdateMessage(job *model.JobInfo, delivered bool) { panic("Not implemented") }

func (m *mockListener) AddCredential(
	job *model.JobInfo,
	role graph.CredentialRole,
	schemaID, credDefID string,
	attributes []*graph.CredentialValue,
	initiatedByUs bool,
) {
	panic("Not implemented")
}
func (m *mockListener) UpdateCredential(job *model.JobInfo, approvedMs, issuedMs, failedMs *int64) {
	m.credTS = *approvedMs
}

func (m *mockListener) AddProof(job *model.JobInfo, role graph.ProofRole, attributes []*graph.ProofAttribute, initiatedByUs bool) {
	panic("Not implemented")
}
func (m *mockListener) UpdateProof(job *model.JobInfo, approvedMs, verifiedMs, failedMs *int64) {
	m.proofTS = *approvedMs
}

func testData(t *testing.T) (*Agency, *model.Agent) {
	return &Agency{
			vault:   &mockListener{},
			options: []grpc.DialOption{grpc.WithContextDialer(dialer(t))},
		},
		&model.Agent{RawJWT: jwt.BuildJWT("test-user")}
}

func TestConnect(t *testing.T) {
	findy, agent := testData(t)
	id, err := findy.Connect(agent, testInvitation)
	if err != nil {
		t.Errorf("Encountered error on connect %v", err)
	}
	if id != testID {
		t.Errorf("Mismatch with id expecting %v, got %v", testID, id)
	}
}

func TestSendMessage(t *testing.T) {
	findy, agent := testData(t)
	id, err := findy.SendMessage(agent, "id", "message")
	if err != nil {
		t.Errorf("Encountered error on connect %v", err)
	}
	if id != testID {
		t.Errorf("Mismatch with id expecting %v, got %v", testID, id)
	}
}

func TestResumeCredentialOffer(t *testing.T) {
	findy, agent := testData(t)
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
	findy, agent := testData(t)
	err := findy.ResumeProofRequest(agent, &model.JobInfo{}, true)
	if err != nil {
		t.Errorf("Encountered error on resume proof request %v", err)
	}
	vault := findy.vault.(*mockListener)
	if vault.proofTS != utils.CurrentTimeMs() {
		t.Errorf("Expected valid proof timestamp %v", vault.proofTS)
	}
}
