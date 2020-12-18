package agency

import (
	"os"
	"testing"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

type agencyListener struct{}

func (l *agencyListener) AddConnection(job *JobInfo, ourDID, theirDID, theirEndpoint, theirLabel string) {
}

func (l *agencyListener) AddMessage(job *JobInfo, message string, sentByMe bool) {}
func (l *agencyListener) UpdateMessage(job *JobInfo, delivered bool)             {}

func (l *agencyListener) AddCredential(
	job *JobInfo,
	role model.CredentialRole,
	schemaID, credDefID string,
	attributes []*model.CredentialValue,
	initiatedByUs bool,
) {
}
func (l *agencyListener) UpdateCredential(job *JobInfo, approvedMs, issuedMs, failedMs *int64) {}

func (l *agencyListener) AddProof(job *JobInfo, role model.ProofRole, attributes []*model.ProofAttribute, initiatedByUs bool) {
}
func (l *agencyListener) UpdateProof(job *JobInfo, approvedMs, verifiedMs, failedMs *int64) {}

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
