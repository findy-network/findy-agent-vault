package model

import "github.com/findy-network/findy-agent-vault/graph/model"

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
	AddConnection(job *JobInfo, connection *Connection)

	AddMessage(job *JobInfo, message *Message)
	UpdateMessage(job *JobInfo, update *MessageUpdate)

	AddCredential(job *JobInfo, credential *Credential)
	UpdateCredential(job *JobInfo, update *CredentialUpdate)

	AddProof(job *JobInfo, proof *Proof)
	UpdateProof(job *JobInfo, update *ProofUpdate)

	FailJob(job *JobInfo)
}

type ArchiveInfo struct {
	AgentID      string
	ConnectionID string
}

type Archiver interface {
	ArchiveConnection(info *ArchiveInfo, connection *Connection)
	ArchiveMessage(info *ArchiveInfo, message *Message)
	ArchiveCredential(info *ArchiveInfo, credential *Credential)
	ArchiveProof(info *ArchiveInfo, proof *Proof)
}
