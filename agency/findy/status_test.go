package findy

import (
	"reflect"
	"testing"

	"github.com/findy-network/findy-agent-api/grpc/agency"
	"github.com/findy-network/findy-agent-vault/agency/model"
	graph "github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/utils"
)

var (
	testConnection = &model.Connection{
		OurDID:        "myDID",
		TheirDID:      "theirDID",
		TheirEndpoint: "theirEndpoint",
		TheirLabel:    "theirLabel",
	}
	testConnectionStatus = func(jobID string) *agency.ProtocolStatus {
		return &agency.ProtocolStatus{
			State: &agency.ProtocolState{
				State: agency.ProtocolState_OK,
				ProtocolId: &agency.ProtocolID{
					Id:     jobID,
					TypeId: agency.Protocol_CONNECT,
				},
			},
			Status: &agency.ProtocolStatus_Connection_{
				Connection: &agency.ProtocolStatus_Connection{
					Id:            "pwName",
					MyDid:         testConnection.OurDID,
					TheirDid:      testConnection.TheirDID,
					TheirEndpoint: testConnection.TheirEndpoint,
					TheirLabel:    testConnection.TheirLabel,
				},
			},
		}
	}

	testMessage       = &model.Message{Message: "messageContent"}
	testMessageStatus = func(jobID string, state agency.ProtocolState_State) *agency.ProtocolStatus {
		return &agency.ProtocolStatus{
			State: &agency.ProtocolState{
				State: state,
				ProtocolId: &agency.ProtocolID{
					Role:   agency.Protocol_ADDRESSEE,
					Id:     jobID,
					TypeId: agency.Protocol_BASIC_MESSAGE,
				},
			},
			Status: &agency.ProtocolStatus_BasicMessage_{
				BasicMessage: &agency.ProtocolStatus_BasicMessage{
					Content:  testMessage.Message,
					SentByMe: testMessage.SentByMe,
				},
			},
		}
	}

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
	testCredentialStatus = func(jobID string, state agency.ProtocolState_State) *agency.ProtocolStatus {
		return &agency.ProtocolStatus{
			State: &agency.ProtocolState{
				State: state,
				ProtocolId: &agency.ProtocolID{
					Role:   agency.Protocol_ADDRESSEE,
					Id:     jobID,
					TypeId: agency.Protocol_ISSUE,
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
		}
	}

	testProof = &model.Proof{
		Role: graph.ProofRoleProver,
		Attributes: []*graph.ProofAttribute{{
			Name:      "attribute-name",
			CredDefID: "cred-def-id",
		}},
		InitiatedByUs: false,
	}
	testProofStatus = func(jobID string, state agency.ProtocolState_State) *agency.ProtocolStatus {
		return &agency.ProtocolStatus{
			State: &agency.ProtocolState{
				State: state,
				ProtocolId: &agency.ProtocolID{
					Role:   agency.Protocol_ADDRESSEE,
					Id:     jobID,
					TypeId: agency.Protocol_PROOF,
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
		}
	}
)

type mockStorage struct {
	info        *model.JobInfo
	connection  *model.Connection
	message     *model.Message
	credUpdate  *model.CredentialUpdate
	proofUpdate *model.ProofUpdate
	credential  *model.Credential
	proof       *model.Proof
	failedJob   *model.JobInfo
}

type statusListener struct {
	conn     *mockStorage
	msg      *mockStorage
	cred     *mockStorage
	credUpt  *mockStorage
	proof    *mockStorage
	proofUpt *mockStorage
	failed   *mockStorage
}

func (s *statusListener) AddConnection(job *model.JobInfo, connection *model.Connection) error {
	s.conn = &mockStorage{info: job, connection: connection}
	return nil
}

func (s *statusListener) AddMessage(job *model.JobInfo, message *model.Message) error {
	s.msg = &mockStorage{info: job, message: message}
	return nil
}

func (s *statusListener) UpdateMessage(job *model.JobInfo, update *model.MessageUpdate) error {
	panic("Not implemented")
}

func (s *statusListener) AddCredential(job *model.JobInfo, credential *model.Credential) error {
	s.cred = &mockStorage{info: job, credential: credential}
	return nil
}

func (s *statusListener) UpdateCredential(job *model.JobInfo, update *model.CredentialUpdate) error {
	s.credUpt = &mockStorage{info: job, credUpdate: update}
	return nil
}

func (s *statusListener) AddProof(job *model.JobInfo, proof *model.Proof) error {
	s.proof = &mockStorage{info: job, proof: proof}
	return nil
}

func (s *statusListener) UpdateProof(job *model.JobInfo, update *model.ProofUpdate) error {
	s.proofUpt = &mockStorage{info: job, proofUpdate: update}
	return nil
}

func (s *statusListener) FailJob(job *model.JobInfo) error {
	s.failed = &mockStorage{failedJob: job}
	return nil
}

func (s *statusListener) connectionStorage() *mockStorage       { return s.conn }
func (s *statusListener) messageStorage() *mockStorage          { return s.msg }
func (s *statusListener) credentialStorage() *mockStorage       { return s.cred }
func (s *statusListener) credentialUpdateStorage() *mockStorage { return s.credUpt }
func (s *statusListener) proofStorage() *mockStorage            { return s.proof }
func (s *statusListener) proofUpdateStorage() *mockStorage      { return s.proofUpt }
func (s *statusListener) failedStorage() *mockStorage           { return s.failed }

func TestHandleNotification(t *testing.T) {
	now := utils.CurrentTimeMs()

	listener := &statusListener{}
	testFindy := &Agency{vault: listener, currentTimeMs: func() int64 { return now }}

	var createJob = func(id string) *model.JobInfo { return &model.JobInfo{JobID: id} }
	const (
		connName        = "connection"
		msgName         = "message"
		credName        = "credential"
		credUpdateName  = "cred_update"
		proofName       = "proof"
		proofUpdateName = "proof_update"
		failedName      = "failed"
		failedCredName  = "failed_cred"
		failedProofName = "failed_proof"
	)
	tests := []struct {
		name         string
		job          *model.JobInfo
		notification *agency.Notification
		status       *agency.ProtocolStatus
		exp          *mockStorage
		got          func() *mockStorage
	}{
		{
			connName,
			createJob(connName),
			&agency.Notification{ProtocolType: agency.Protocol_CONNECT},
			testConnectionStatus(connName),
			&mockStorage{info: createJob(connName), connection: testConnection},
			listener.connectionStorage,
		},
		{
			msgName,
			createJob(msgName),
			&agency.Notification{ProtocolType: agency.Protocol_BASIC_MESSAGE},
			testMessageStatus(msgName, agency.ProtocolState_OK),
			&mockStorage{info: createJob(msgName), message: testMessage},
			listener.messageStorage,
		},
		{
			credName,
			createJob(credName),
			&agency.Notification{
				TypeId:       agency.Notification_ACTION_NEEDED,
				ProtocolType: agency.Protocol_ISSUE,
				Role:         agency.Protocol_ADDRESSEE,
			},
			testCredentialStatus(credName, agency.ProtocolState_OK),
			&mockStorage{info: createJob(credName), credential: testCredential},
			listener.credentialStorage,
		},
		{
			credUpdateName,
			createJob(credUpdateName),
			&agency.Notification{
				ProtocolType: agency.Protocol_ISSUE,
				Role:         agency.Protocol_ADDRESSEE,
			},
			testCredentialStatus(credUpdateName, agency.ProtocolState_OK),
			&mockStorage{info: createJob(credUpdateName), credUpdate: &model.CredentialUpdate{IssuedMs: &now}},
			listener.credentialUpdateStorage,
		},
		{
			proofName,
			createJob(proofName),
			&agency.Notification{
				TypeId:       agency.Notification_ACTION_NEEDED,
				ProtocolType: agency.Protocol_PROOF,
				Role:         agency.Protocol_ADDRESSEE,
			},
			testProofStatus(proofName, agency.ProtocolState_OK),
			&mockStorage{info: createJob(proofName), proof: testProof},
			listener.proofStorage,
		},
		{
			proofUpdateName,
			createJob(proofUpdateName),
			&agency.Notification{
				ProtocolType: agency.Protocol_PROOF,
				Role:         agency.Protocol_ADDRESSEE,
			},
			testProofStatus(proofUpdateName, agency.ProtocolState_OK),
			&mockStorage{info: createJob(proofUpdateName), proofUpdate: &model.ProofUpdate{VerifiedMs: &now}},
			listener.proofUpdateStorage,
		},
		{
			failedName,
			createJob(failedName),
			&agency.Notification{ProtocolType: agency.Protocol_BASIC_MESSAGE},
			testMessageStatus(failedName, agency.ProtocolState_ERR),
			&mockStorage{failedJob: createJob(failedName)},
			listener.failedStorage,
		},
		{
			failedCredName,
			createJob(failedCredName),
			&agency.Notification{ProtocolType: agency.Protocol_ISSUE},
			testCredentialStatus(failedCredName, agency.ProtocolState_ERR),
			&mockStorage{info: createJob(failedCredName), credUpdate: &model.CredentialUpdate{FailedMs: &now}},
			listener.credentialUpdateStorage,
		},
		{
			failedProofName,
			createJob(failedProofName),
			&agency.Notification{ProtocolType: agency.Protocol_PROOF},
			testProofStatus(failedProofName, agency.ProtocolState_ERR),
			&mockStorage{info: createJob(failedProofName), proofUpdate: &model.ProofUpdate{FailedMs: &now}},
			listener.proofUpdateStorage,
		},
	}

	for _, testCase := range tests {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			testFindy.handleNotification(&model.Agent{}, tc.job, tc.notification, tc.status)
			if !reflect.DeepEqual(tc.exp, tc.got()) {
				t.Errorf("Mismatch in status %s, expected: %+v  got: %+v", tc.name, tc.exp, tc.got())
			}
		})
	}
}

func TestGetStatus(t *testing.T) {
	status, ok := findy.getStatus(&model.Agent{}, &agency.Notification{ProtocolId: "id", ProtocolType: agency.Protocol_ISSUE})
	if !ok {
		t.Errorf("Failure when getting status")
	}
	if status == nil {
		t.Errorf("Received nil status")
	}
}
