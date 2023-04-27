package model

import (
	dbModel "github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/graph/model"
)

type Connection struct {
	OurDID, TheirDID, TheirEndpoint, TheirLabel string
}

type Message struct {
	Message  string
	SentByMe bool
}

type MessageUpdate struct {
	Delivered bool
}

type Credential struct {
	Role                model.CredentialRole
	SchemaID, CredDefID string
	Attributes          []*model.CredentialValue
	InitiatedByUs       bool
}

type CredentialUpdate struct {
	ApprovedMs, IssuedMs, FailedMs *int64
}

type Proof struct {
	Role          model.ProofRole
	Attributes    []*model.ProofAttribute
	InitiatedByUs bool
}

type ProofUpdate struct {
	ApprovedMs, VerifiedMs, FailedMs *int64
}

type Listener interface {
	AddConnection(job *JobInfo, connection *Connection) error

	AddMessage(job *JobInfo, message *Message) error
	UpdateMessage(job *JobInfo, update *MessageUpdate) error

	AddCredential(job *JobInfo, credential *Credential) (*dbModel.Job, error)
	UpdateCredential(job *JobInfo, credential *Credential, update *CredentialUpdate) error

	AddProof(job *JobInfo, proof *Proof) (*dbModel.Job, error)
	UpdateProof(job *JobInfo, proof *Proof, update *ProofUpdate) error

	FailJob(job *JobInfo) error
}

type ArchiveInfo struct {
	AgentID       string
	JobID         string
	ConnectionID  string
	InitiatedByUs bool
}

type Archiver interface {
	ArchiveConnection(info *ArchiveInfo, connection *Connection)
	ArchiveMessage(info *ArchiveInfo, message *Message)
	ArchiveCredential(info *ArchiveInfo, credential *Credential)
	ArchiveProof(info *ArchiveInfo, proof *Proof)
}
