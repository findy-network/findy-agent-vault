package listen

import (
	"os"
	"testing"
	"time"

	agency "github.com/findy-network/findy-agent-vault/agency/model"
	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/db/store"
	graph "github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/resolver/query/agent"
	"github.com/findy-network/findy-agent-vault/resolver/update"
	gomock "github.com/golang/mock/gomock"

	"github.com/findy-network/findy-agent-vault/utils"
)

func setup() {
	utils.SetLogDefaults()
}

func teardown() {
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func createListener(db store.DB) *Listener {
	agentResolver := agent.NewResolver(db, nil)
	updater := update.NewUpdater(db, agentResolver)
	return &Listener{db, updater}
}

func TestAddConnection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockDB(ctrl)

	var (
		job        = &agency.JobInfo{JobID: "job-id", TenantID: "tenant-id", ConnectionID: "connection-id"}
		connection = &agency.Connection{
			OurDID:        "ourDID",
			TheirDID:      "theirDID",
			TheirEndpoint: "theirEndpoint",
			TheirLabel:    "theirLabel",
		}
		resultJob        = &model.Job{Base: model.Base{ID: job.JobID, TenantID: job.TenantID}}
		now              = utils.CurrentTime()
		resultConnection = &model.Connection{
			Base: model.Base{
				ID:       job.ConnectionID,
				TenantID: job.TenantID,
			},
			OurDid:        connection.OurDID,
			TheirDid:      connection.TheirDID,
			TheirEndpoint: connection.TheirEndpoint,
			TheirLabel:    connection.TheirLabel,
			Approved:      now,
			Invited:       false,
		}
		event = &model.Event{
			Base:         model.Base{TenantID: job.TenantID},
			Read:         false,
			Description:  "Established connection to theirLabel",
			ConnectionID: &job.ConnectionID,
			JobID:        &job.JobID,
		}
	)

	m.
		EXPECT().
		GetJob(gomock.Eq(job.JobID), gomock.Eq(job.TenantID)).
		Return(resultJob, nil)
	m.
		EXPECT().
		AddConnection(gomock.Any()). // TODO: custom matcher
		Return(resultConnection, nil)
	m.
		EXPECT().
		UpdateJob(resultJob).
		Return(resultJob, nil)
	m.
		EXPECT().
		AddEvent(event).
		Return(event, nil)

	l := createListener(m)

	_ = l.AddConnection(job, connection)
}

func TestAddMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockDB(ctrl)
	var (
		job     = &agency.JobInfo{JobID: "job-id", TenantID: "tenant-id", ConnectionID: "connection-id"}
		message = &agency.Message{
			Message:  "message",
			SentByMe: false,
		}
		resultMessage = &model.Message{
			Base:         model.Base{TenantID: job.TenantID},
			ConnectionID: job.ConnectionID,
			Message:      message.Message,
			SentByMe:     message.SentByMe,
		}
		resultJob = &model.Job{
			Base:              model.Base{ID: job.JobID, TenantID: job.TenantID},
			ConnectionID:      &job.ConnectionID,
			ProtocolType:      graph.ProtocolTypeBasicMessage,
			ProtocolMessageID: &resultMessage.ID,
			InitiatedByUs:     message.SentByMe,
			Status:            graph.JobStatusComplete,
			Result:            graph.JobResultSuccess,
		}
		event = &model.Event{
			Base:         model.Base{TenantID: job.TenantID},
			Read:         false,
			Description:  resultMessage.Description(),
			ConnectionID: &job.ConnectionID,
			JobID:        &job.JobID,
		}
	)

	m.
		EXPECT().
		AddMessage(resultMessage).
		Return(resultMessage, nil)
	m.
		EXPECT().
		AddJob(resultJob).
		Return(resultJob, nil)
	m.
		EXPECT().
		AddEvent(event).
		Return(event, nil)

	l := createListener(m)

	_ = l.AddMessage(job, message)
}

func TestAddCredential(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockDB(ctrl)
	var (
		job        = &agency.JobInfo{JobID: "job-id", TenantID: "tenant-id", ConnectionID: "connection-id"}
		credential = &agency.Credential{
			Role:      graph.CredentialRoleHolder,
			SchemaID:  "schema-id",
			CredDefID: "cred-def-id",
			Attributes: []*graph.CredentialValue{{
				Name:  "attribute-name",
				Value: "attribute-value",
			}},
			InitiatedByUs: false,
		}
		resultCredential = &model.Credential{
			Base:          model.Base{TenantID: job.TenantID},
			ConnectionID:  job.ConnectionID,
			Role:          credential.Role,
			SchemaID:      credential.SchemaID,
			CredDefID:     credential.CredDefID,
			Attributes:    credential.Attributes,
			InitiatedByUs: credential.InitiatedByUs,
		}
		resultJob = &model.Job{
			Base:                 model.Base{ID: job.JobID, TenantID: job.TenantID},
			ConnectionID:         &job.ConnectionID,
			ProtocolType:         graph.ProtocolTypeCredential,
			ProtocolCredentialID: &resultCredential.ID,
			InitiatedByUs:        credential.InitiatedByUs,
			Status:               graph.JobStatusPending,
			Result:               graph.JobResultNone,
		}
		event = &model.Event{
			Base:         model.Base{TenantID: job.TenantID},
			Read:         false,
			Description:  resultCredential.Description(),
			ConnectionID: &job.ConnectionID,
			JobID:        &job.JobID,
		}
	)

	m.
		EXPECT().
		AddCredential(resultCredential).
		Return(resultCredential, nil)
	m.
		EXPECT().
		AddJob(resultJob).
		Return(resultJob, nil)
	m.
		EXPECT().
		AddEvent(event).
		Return(event, nil)

	l := createListener(m)

	_ = l.AddCredential(job, credential)
}

func TestUpdateCredential(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockDB(ctrl)
	var (
		now              = utils.CurrentTimeMs()
		credentialID     = "credential-id"
		job              = &agency.JobInfo{JobID: "job-id", TenantID: "tenant-id", ConnectionID: "connection-id"}
		credentialUpdate = &agency.CredentialUpdate{
			ApprovedMs: &now,
		}
		resultCredential = &model.Credential{
			Base:     model.Base{TenantID: job.TenantID},
			Role:     graph.CredentialRoleHolder,
			Approved: utils.TSToTimeIfNotSet(nil, credentialUpdate.ApprovedMs),
			Issued:   utils.TSToTimeIfNotSet(nil, &now),
		}
		resultJob = &model.Job{
			Base:                 model.Base{ID: job.JobID, TenantID: job.TenantID},
			ConnectionID:         &job.ConnectionID,
			ProtocolCredentialID: &credentialID,
		}
		event = &model.Event{
			Base:         model.Base{TenantID: job.TenantID},
			Read:         false,
			Description:  resultCredential.Description(),
			ConnectionID: &job.ConnectionID,
			JobID:        &job.JobID,
		}
	)

	m.
		EXPECT().
		GetJob(gomock.Eq(job.JobID), gomock.Eq(job.TenantID)).
		Return(resultJob, nil)
	m.
		EXPECT().
		GetCredential(gomock.Eq(credentialID), gomock.Eq(job.TenantID)).
		Return(resultCredential, nil)
	m.
		EXPECT().
		UpdateCredential(resultCredential).
		Return(resultCredential, nil)
	m.
		EXPECT().
		UpdateJob(resultJob).
		Return(resultJob, nil)
	m.
		EXPECT().
		AddEvent(event).
		Return(event, nil)

	m.
		EXPECT().
		GetOpenProofJobs(job.TenantID, []*graph.ProofAttribute{})

	l := createListener(m)

	_ = l.UpdateCredential(job, credentialUpdate)
}

func TestAddProof(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	utils.CurrentStaticTime = utils.CurrentTime()

	m := NewMockDB(ctrl)
	var (
		now   = utils.CurrentTime()
		job   = &agency.JobInfo{JobID: "job-id", TenantID: "tenant-id", ConnectionID: "connection-id"}
		proof = &agency.Proof{
			Role: graph.ProofRoleProver,
			Attributes: []*graph.ProofAttribute{{
				Name:      "attribute-name",
				CredDefID: "cred-def-id",
			}},
			InitiatedByUs: false,
		}
		resultProof = &model.Proof{
			Base:          model.Base{TenantID: job.TenantID},
			ConnectionID:  job.ConnectionID,
			Role:          proof.Role,
			Attributes:    proof.Attributes,
			Result:        false,
			InitiatedByUs: proof.InitiatedByUs,
			Provable:      &now,
		}
		resultJob = &model.Job{
			Base:            model.Base{ID: job.JobID, TenantID: job.TenantID},
			ConnectionID:    &job.ConnectionID,
			ProtocolType:    graph.ProtocolTypeProof,
			ProtocolProofID: &resultProof.ID,
			InitiatedByUs:   proof.InitiatedByUs,
			Status:          graph.JobStatusBlocked,
			Result:          graph.JobResultNone,
		}
		event = &model.Event{
			Base:         model.Base{TenantID: job.TenantID},
			Read:         false,
			Description:  resultProof.Description(),
			ConnectionID: &job.ConnectionID,
			JobID:        &job.JobID,
		}
	)

	m.
		EXPECT().
		AddProof(resultProof).
		Return(resultProof, nil)
	m.
		EXPECT().
		SearchCredentials(job.TenantID, proof.Attributes)
	m.
		EXPECT().
		AddJob(resultJob).
		Return(resultJob, nil)
	m.
		EXPECT().
		AddEvent(event).
		Return(event, nil)

	l := createListener(m)

	_ = l.AddProof(job, proof)

	utils.CurrentStaticTime = time.Time{}
}

func TestUpdateProof(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockDB(ctrl)
	var (
		now         = utils.CurrentTimeMs()
		proofID     = "proof-id"
		job         = &agency.JobInfo{JobID: "job-id", TenantID: "tenant-id", ConnectionID: "connection-id"}
		proofUpdate = &agency.ProofUpdate{
			ApprovedMs: &now,
		}
		resultProof = &model.Proof{
			Base:     model.Base{TenantID: job.TenantID},
			Role:     graph.ProofRoleProver,
			Approved: utils.TSToTimePtrIfNotSet(nil, proofUpdate.ApprovedMs),
		}

		resultJob = &model.Job{
			Base:            model.Base{ID: job.JobID, TenantID: job.TenantID},
			ConnectionID:    &job.ConnectionID,
			ProtocolProofID: &proofID,
		}
		event = &model.Event{
			Base:         model.Base{TenantID: job.TenantID},
			Read:         false,
			Description:  resultProof.Description(),
			ConnectionID: &job.ConnectionID,
			JobID:        &job.JobID,
		}
	)

	m.
		EXPECT().
		GetJob(gomock.Eq(job.JobID), gomock.Eq(job.TenantID)).
		Return(resultJob, nil)
	m.
		EXPECT().
		GetProof(gomock.Eq(proofID), gomock.Eq(job.TenantID)).
		Return(resultProof, nil)
	m.
		EXPECT().
		UpdateProof(resultProof).
		Return(resultProof, nil)
	m.
		EXPECT().
		UpdateJob(resultJob).
		Return(resultJob, nil)
	m.
		EXPECT().
		AddEvent(event).
		Return(event, nil)

	l := createListener(m)

	_ = l.UpdateProof(job, proofUpdate)
}
