package findy

import (
	"testing"

	"github.com/findy-network/findy-agent-api/grpc/agency"
	"github.com/findy-network/findy-agent-vault/agency/model"
	graph "github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/utils"
)

type statusListener struct {
	// connection
	connJobID     string
	myDID         string
	theirDID      string
	theirEndpoint string
	theirLabel    string

	// message
	msgJobID       string
	messageContent string
	sentByMe       bool

	// cred
	issueJobID string
	issued     int64

	// proof
	proofJobID string
	verified   int64
}

func (s *statusListener) AddConnection(job *model.JobInfo, ourDID, theirDID, theirEndpoint, theirLabel string) {
	s.connJobID = job.JobID
	s.myDID = ourDID
	s.theirDID = theirDID
	s.theirEndpoint = theirEndpoint
	s.theirLabel = theirLabel
}

func (s *statusListener) AddMessage(job *model.JobInfo, message string, sentByMe bool) {
	s.msgJobID = job.JobID
	s.messageContent = message
	s.sentByMe = sentByMe
}

func (s *statusListener) UpdateMessage(job *model.JobInfo, delivered bool) { panic("Not implemented") }

func (s *statusListener) AddCredential(
	job *model.JobInfo,
	role graph.CredentialRole,
	schemaID, credDefID string,
	attributes []*graph.CredentialValue,
	initiatedByUs bool,
) {
	panic("Not implemented")
}
func (s *statusListener) UpdateCredential(job *model.JobInfo, approvedMs, issuedMs, failedMs *int64) {
	s.issueJobID = job.JobID
	s.issued = *issuedMs
}

func (s *statusListener) AddProof(job *model.JobInfo, role graph.ProofRole, attributes []*graph.ProofAttribute, initiatedByUs bool) {
	panic("Not implemented")
}

func (s *statusListener) UpdateProof(job *model.JobInfo, approvedMs, verifiedMs, failedMs *int64) {
	s.proofJobID = job.JobID
	s.verified = *verifiedMs
}

func TestHandleConnectionStatus(t *testing.T) {
	const (
		jobID         = "conn-job-id"
		myDID         = "myDID"
		theirDID      = "theirDID"
		theirEndpoint = "theirEndpoint"
		theirLabel    = "theirLabel"
	)
	listener := &statusListener{}
	testFindy := &Agency{vault: listener}
	testFindy.handleStatus(
		&model.JobInfo{JobID: jobID},
		&agency.Notification{ProtocolType: agency.Protocol_CONNECT},
		&agency.ProtocolStatus{
			Status: &agency.ProtocolStatus_Connection_{Connection: &agency.ProtocolStatus_Connection{
				Id:            "pwName",
				MyDid:         myDID,
				TheirDid:      theirDID,
				TheirEndpoint: theirEndpoint,
				TheirLabel:    theirLabel,
			},
			},
		})

	if jobID != listener.connJobID {
		t.Errorf("Mismatch on connection status job id, expected %v, got %v", jobID, listener.connJobID)
	}
	if myDID != listener.myDID {
		t.Errorf("Mismatch on connection status my did, expected %v, got %v", myDID, listener.myDID)
	}
	if theirDID != listener.theirDID {
		t.Errorf("Mismatch on connection status their did, expected %v, got %v", theirDID, listener.theirDID)
	}
	if theirEndpoint != listener.theirEndpoint {
		t.Errorf("Mismatch on connection status their endpoint, expected %v, got %v", theirEndpoint, listener.theirEndpoint)
	}
	if theirLabel != listener.theirLabel {
		t.Errorf("Mismatch on connection status their label, expected %v, got %v", theirLabel, listener.theirLabel)
	}
}

func TestHandleBasicMessageStatus(t *testing.T) {
	const (
		jobID          = "msg-job-id"
		messageContent = "messageContent"
		sentByMe       = false
	)
	listener := &statusListener{}
	testFindy := &Agency{vault: listener}
	testFindy.handleStatus(
		&model.JobInfo{JobID: jobID},
		&agency.Notification{ProtocolType: agency.Protocol_BASIC_MESSAGE},
		&agency.ProtocolStatus{
			Status: &agency.ProtocolStatus_BasicMessage_{BasicMessage: &agency.ProtocolStatus_BasicMessage{
				Content:  messageContent,
				SentByMe: sentByMe,
			},
			},
		})

	if jobID != listener.msgJobID {
		t.Errorf("Mismatch on message status job id, expected %v, got %v", jobID, listener.msgJobID)
	}
	if messageContent != listener.messageContent {
		t.Errorf("Mismatch on message status job id, expected %v, got %v", messageContent, listener.messageContent)
	}
	if sentByMe != listener.sentByMe {
		t.Errorf("Mismatch on message status sent by me, expected %v, got %v", sentByMe, listener.sentByMe)
	}
}

func TestHandleCredentialStatus(t *testing.T) {
	var (
		jobID = "issue-job-id"
		now   = utils.CurrentTimeMs()
	)
	listener := &statusListener{}
	testFindy := &Agency{vault: listener}
	testFindy.handleStatus(
		&model.JobInfo{JobID: jobID},
		&agency.Notification{ProtocolType: agency.Protocol_ISSUE},
		&agency.ProtocolStatus{
			Status: &agency.ProtocolStatus_Issue_{Issue: &agency.ProtocolStatus_Issue{}},
		})

	if jobID != listener.issueJobID {
		t.Errorf("Mismatch on issue status job id, expected %v, got %v", jobID, listener.issueJobID)
	}
	if now != listener.issued {
		t.Errorf("Mismatch on issue status issued ts, expected %v, got %v", now, listener.issued)
	}
}

func TestHandleProofStatus(t *testing.T) {
	var (
		jobID = "proof-job-id"
		now   = utils.CurrentTimeMs()
	)
	listener := &statusListener{}
	testFindy := &Agency{vault: listener}
	testFindy.handleStatus(
		&model.JobInfo{JobID: jobID},
		&agency.Notification{ProtocolType: agency.Protocol_PROOF},
		&agency.ProtocolStatus{
			Status: &agency.ProtocolStatus_Proof{Proof: &agency.Protocol_Proof{}},
		})

	if jobID != listener.proofJobID {
		t.Errorf("Mismatch on proof status job id, expected %v, got %v", jobID, listener.proofJobID)
	}
	if now != listener.verified {
		t.Errorf("Mismatch on proof status issued ts, expected %v, got %v", now, listener.verified)
	}
}
