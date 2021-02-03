package findy

import (
	"reflect"
	"testing"

	"github.com/findy-network/findy-agent-api/grpc/agency"
	"github.com/findy-network/findy-agent-vault/agency/model"
	graph "github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/utils"
)

type statusListener struct {
	connJob    *model.JobInfo
	connection *model.Connection

	msgJob  *model.JobInfo
	message *model.Message

	credUpdateJob *model.JobInfo
	credUpdate    *model.CredentialUpdate

	proofUpdateJob *model.JobInfo
	proofUpdate    *model.ProofUpdate

	credentialJob *model.JobInfo
	credential    *model.Credential

	proofJob *model.JobInfo
	proof    *model.Proof
}

func (s *statusListener) AddConnection(job *model.JobInfo, connection *model.Connection) error {
	s.connJob = job
	s.connection = connection
	return nil
}

func (s *statusListener) AddMessage(job *model.JobInfo, message *model.Message) error {
	s.msgJob = job
	s.message = message
	return nil
}

func (s *statusListener) UpdateMessage(job *model.JobInfo, update *model.MessageUpdate) error {
	panic("Not implemented")
}

func (s *statusListener) AddCredential(job *model.JobInfo, credential *model.Credential) error {
	s.credentialJob = job
	s.credential = credential
	return nil
}

func (s *statusListener) UpdateCredential(job *model.JobInfo, update *model.CredentialUpdate) error {
	s.credUpdateJob = job
	s.credUpdate = update
	return nil
}

func (s *statusListener) AddProof(job *model.JobInfo, proof *model.Proof) error {
	s.proofJob = job
	s.proof = proof
	return nil
}

func (s *statusListener) UpdateProof(job *model.JobInfo, update *model.ProofUpdate) error {
	s.proofUpdateJob = job
	s.proofUpdate = update
	return nil
}

func (s *statusListener) FailJob(job *model.JobInfo) error {
	panic("Not implemented")
}

func TestHandleConnectionStatus(t *testing.T) {
	var (
		testJob        = &model.JobInfo{JobID: "conn-job-id"}
		testConnection = &model.Connection{
			OurDID:        "myDID",
			TheirDID:      "theirDID",
			TheirEndpoint: "theirEndpoint",
			TheirLabel:    "theirLabel",
		}
	)
	listener := &statusListener{}
	testFindy := &Agency{vault: listener}
	testFindy.handleStatus(
		&model.JobInfo{JobID: testJob.JobID},
		&agency.Notification{ProtocolType: agency.Protocol_CONNECT},
		&agency.ProtocolStatus{
			State: &agency.ProtocolState{
				State: agency.ProtocolState_OK,
			},
			Status: &agency.ProtocolStatus_Connection_{Connection: &agency.ProtocolStatus_Connection{
				Id:            "pwName",
				MyDid:         testConnection.OurDID,
				TheirDid:      testConnection.TheirDID,
				TheirEndpoint: testConnection.TheirEndpoint,
				TheirLabel:    testConnection.TheirLabel,
			},
			},
		})

	if !reflect.DeepEqual(testJob, listener.connJob) {
		t.Errorf("Mismatch in connection job  expected: %v  got: %v", testJob, listener.connJob)
	}
	if !reflect.DeepEqual(testConnection, listener.connection) {
		t.Errorf("Mismatch in connection  expected: %v  got: %v", testConnection, listener.connection)
	}
}

func TestHandleBasicMessageStatus(t *testing.T) {
	var (
		testJob     = &model.JobInfo{JobID: "msg-job-id"}
		testMessage = &model.Message{Message: "messageContent"}
	)
	listener := &statusListener{}
	testFindy := &Agency{vault: listener}
	testFindy.handleStatus(
		&model.JobInfo{JobID: testJob.JobID},
		&agency.Notification{ProtocolType: agency.Protocol_BASIC_MESSAGE},
		&agency.ProtocolStatus{
			State: &agency.ProtocolState{
				State: agency.ProtocolState_OK,
				ProtocolId: &agency.ProtocolID{
					Role: agency.Protocol_ADDRESSEE,
				},
			},
			Status: &agency.ProtocolStatus_BasicMessage_{BasicMessage: &agency.ProtocolStatus_BasicMessage{
				Content:  testMessage.Message,
				SentByMe: testMessage.SentByMe,
			},
			},
		})

	if !reflect.DeepEqual(testJob, listener.msgJob) {
		t.Errorf("Mismatch in message job  expected: %v  got: %v", testJob, listener.msgJob)
	}
	if !reflect.DeepEqual(testMessage, listener.message) {
		t.Errorf("Mismatch in message  expected: %v  got: %v", testMessage, listener.message)
	}
}

func TestHandleCredentialStatus(t *testing.T) {
	var (
		now        = utils.CurrentTimeMs()
		testJob    = &model.JobInfo{JobID: "issue-job-id"}
		testUpdate = &model.CredentialUpdate{IssuedMs: &now}
	)
	listener := &statusListener{}
	testFindy := &Agency{vault: listener}
	testFindy.handleStatus(
		&model.JobInfo{JobID: testJob.JobID},
		&agency.Notification{ProtocolType: agency.Protocol_ISSUE},
		&agency.ProtocolStatus{
			State: &agency.ProtocolState{
				State: agency.ProtocolState_OK,
			},
			Status: &agency.ProtocolStatus_Issue_{Issue: &agency.ProtocolStatus_Issue{}},
		})

	if !reflect.DeepEqual(testJob, listener.credUpdateJob) {
		t.Errorf("Mismatch in cred update job  expected: %v  got: %v", testJob, listener.credUpdateJob)
	}
	if !reflect.DeepEqual(testUpdate, listener.credUpdate) {
		t.Errorf("Mismatch in cred update  expected: %v  got: %v", testUpdate, listener.credUpdate)
	}
}

func TestHandleProofStatus(t *testing.T) {
	var (
		now        = utils.CurrentTimeMs()
		testJob    = &model.JobInfo{JobID: "proof-job-id"}
		testUpdate = &model.ProofUpdate{VerifiedMs: &now}
	)
	listener := &statusListener{}
	testFindy := &Agency{vault: listener}
	testFindy.handleStatus(
		&model.JobInfo{JobID: testJob.JobID},
		&agency.Notification{ProtocolType: agency.Protocol_PROOF},
		&agency.ProtocolStatus{
			State: &agency.ProtocolState{
				State: agency.ProtocolState_OK,
			},
			Status: &agency.ProtocolStatus_Proof{Proof: &agency.Protocol_Proof{}},
		})

	if !reflect.DeepEqual(testJob, listener.proofUpdateJob) {
		t.Errorf("Mismatch in proof update job  expected: %v  got: %v", testJob, listener.proofUpdateJob)
	}
	if !reflect.DeepEqual(testUpdate, listener.proofUpdate) {
		t.Errorf("Mismatch in proof update  expected: %v  got: %v", testUpdate, listener.proofUpdate)
	}
}

func TestHandleCredentialAction(t *testing.T) {
	var (
		testJob        = &model.JobInfo{JobID: "cred-job-id"}
		testCredential = &model.Credential{
			Role:      graph.CredentialRoleHolder,
			SchemaID:  "schema-id",
			CredDefID: "cred-def-id",
			Attributes: []*graph.CredentialValue{{
				Name:  "attribute-name",
				Value: "attribute-value",
			}},
			InitiatedByUs: false,
		}
	)
	listener := &statusListener{}
	testFindy := &Agency{vault: listener}
	testFindy.handleAction(
		&model.JobInfo{JobID: testJob.JobID},
		&agency.Notification{
			ProtocolType: agency.Protocol_ISSUE,
			Role:         agency.Protocol_ADDRESSEE,
		},
		&agency.ProtocolStatus{
			State: &agency.ProtocolState{
				ProtocolId: &agency.ProtocolID{
					Role: agency.Protocol_ADDRESSEE,
				},
			},
			Status: &agency.ProtocolStatus_Issue_{
				Issue: &agency.ProtocolStatus_Issue{
					SchemaId:  testCredential.SchemaID,
					CredDefId: testCredential.CredDefID,
					Attrs: []*agency.Protocol_Attribute{
						{
							Name:  testCredential.Attributes[0].Name,
							Value: testCredential.Attributes[0].Value,
						},
					},
				},
			},
		})

	if !reflect.DeepEqual(testJob, listener.credentialJob) {
		t.Errorf("Mismatch in cred job  expected: %v  got: %v", testJob, listener.credentialJob)
	}
	if !reflect.DeepEqual(testCredential, listener.credential) {
		t.Errorf("Mismatch in cred update  expected: %v  got: %v", testCredential, listener.credential)
	}
}

func TestHandleProofAction(t *testing.T) {
	var (
		value     = ""
		testJob   = &model.JobInfo{JobID: "cred-job-id"}
		testProof = &model.Proof{
			Role: graph.ProofRoleProver,
			Attributes: []*graph.ProofAttribute{{
				Name:      "attribute-name",
				CredDefID: "cred-def-id",
				Value:     &value,
			}},
			InitiatedByUs: false,
		}
	)
	listener := &statusListener{}
	testFindy := &Agency{vault: listener}
	testFindy.handleAction(
		&model.JobInfo{JobID: testJob.JobID},
		&agency.Notification{
			ProtocolType: agency.Protocol_PROOF,
			Role:         agency.Protocol_ADDRESSEE,
		},
		&agency.ProtocolStatus{
			State: &agency.ProtocolState{
				ProtocolId: &agency.ProtocolID{
					Role: agency.Protocol_ADDRESSEE,
				},
			},
			Status: &agency.ProtocolStatus_Proof{
				Proof: &agency.Protocol_Proof{
					Attrs: []*agency.Protocol_Proof_Attr{
						{
							Name:      testProof.Attributes[0].Name,
							CredDefId: testProof.Attributes[0].CredDefID,
						},
					},
				},
			},
		})

	if !reflect.DeepEqual(testJob, listener.proofJob) {
		t.Errorf("Mismatch in proof job  expected: %v  got: %v", testJob, listener.proofJob)
	}
	if !reflect.DeepEqual(testProof, listener.proof) {
		t.Errorf("Mismatch in proof update  expected: %v  got: %v", testProof, listener.proof)
	}
}

func TestGetStatus(t *testing.T) {
	status, ok := findy.getStatus(findy.userCmdConn(agent), &agency.Notification{
		ProtocolType: agency.Protocol_PROOF,
		Role:         agency.Protocol_ADDRESSEE,
	})
	if !ok {
		t.Errorf("Failed to fetch status")
	}
	if status == nil {
		t.Errorf("Received nil status")
	}
}
