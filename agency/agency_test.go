package agency

import (
	"os"
	"testing"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

type agencyListener struct{}

func (l *agencyListener) AddConnection(id, ourDID, theirDID, theirEndpoint, theirLabel string) {

}

func (l *agencyListener) AddMessage(connectionID, id, message string, sentByMe bool) {

}

func (l *agencyListener) UpdateMessage(connectionID, id, delivered bool) {

}

func (l *agencyListener) AddCredential(
	connectionID, id string,
	role model.CredentialRole,
	schemaID, credDefID string,
	attributes []*model.CredentialValue,
	initiatedByUs bool,
) {

}

func (l *agencyListener) UpdateCredential(connectionID, id string, approvedMs, issuedMs, failedMs *int64) {

}

func (l *agencyListener) AddProof(connectionID, id string, role model.ProofRole, attributes []*model.ProofAttribute, initiatedByUs bool) {

}

func (l *agencyListener) UpdateProof(connectionID, id string, approvedMs, verifiedMs, failedMs *int64) {

}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	Instance.Init(&agencyListener{})
}

func teardown() {
}
