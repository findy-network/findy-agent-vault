package findy

import (
	"reflect"
	"testing"

	agency "github.com/findy-network/findy-common-go/grpc/agency/v1"
	"github.com/findy-network/findy-agent-vault/agency/model"
)

var (
	testArchiveInfo = &model.ArchiveInfo{
		AgentID:      "agent-id",
		ConnectionID: "connection-id",
	}
)

type mockArchive struct {
	info       *model.ArchiveInfo
	connection *model.Connection
	message    *model.Message
	credential *model.Credential
	proof      *model.Proof
}

type mockArchiver struct {
	connection *mockArchive
	message    *mockArchive
	credential *mockArchive
	proof      *mockArchive
}

func (m *mockArchiver) ArchiveConnection(info *model.ArchiveInfo, connection *model.Connection) {
	m.connection = &mockArchive{info: info, connection: connection}
}

func (m *mockArchiver) ArchiveMessage(info *model.ArchiveInfo, message *model.Message) {
	m.message = &mockArchive{info: info, message: message}
}

func (m *mockArchiver) ArchiveCredential(info *model.ArchiveInfo, credential *model.Credential) {
	m.credential = &mockArchive{info: info, credential: credential}
}

func (m *mockArchiver) ArchiveProof(info *model.ArchiveInfo, proof *model.Proof) {
	m.proof = &mockArchive{info: info, proof: proof}
}

func (m *mockArchiver) connectionArchive() *mockArchive { return m.connection }
func (m *mockArchiver) messageArchive() *mockArchive    { return m.message }
func (m *mockArchiver) credentialArchive() *mockArchive { return m.credential }
func (m *mockArchiver) proofArchive() *mockArchive      { return m.proof }

func TestArchive(t *testing.T) {
	archiver := &mockArchiver{}
	testFindy := &Agency{archiver: archiver}

	tests := []struct {
		name   string
		status *agency.ProtocolStatus
		exp    *mockArchive
		got    func() *mockArchive
	}{
		{
			"connection",
			testConnectionStatus(""),
			&mockArchive{info: testArchiveInfo, connection: testConnection},
			archiver.connectionArchive,
		},
		{
			"message",
			testMessageStatus("", agency.ProtocolState_OK),
			&mockArchive{info: testArchiveInfo, message: testMessage},
			archiver.messageArchive,
		},
		{
			"credential",
			testCredentialStatus("", agency.ProtocolState_OK),
			&mockArchive{info: testArchiveInfo, credential: testCredential},
			archiver.credentialArchive,
		},
		{
			"proof",
			testProofStatus("", agency.ProtocolState_OK),
			&mockArchive{info: testArchiveInfo, proof: testProof},
			archiver.proofArchive,
		},
	}

	for _, testCase := range tests {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			testFindy.archive(testArchiveInfo, tc.status)
			if !reflect.DeepEqual(tc.exp, tc.got()) {
				t.Errorf("Mismatch in archive %s, expected: %v  got: %v", tc.name, tc.exp, tc.got())
			}
		})
	}
}
