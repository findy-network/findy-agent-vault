// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type BasicMessage struct {
	ID         string    `json:"id"`
	Message    string    `json:"message"`
	SentByMe   bool      `json:"sentByMe"`
	Delivered  *bool     `json:"delivered"`
	CreatedMs  string    `json:"createdMs"`
	Connection *Pairwise `json:"connection"`
}

type BasicMessageConnection struct {
	ConnectionID *string             `json:"ConnectionId"`
	Edges        []*BasicMessageEdge `json:"edges"`
	Nodes        []*BasicMessage     `json:"nodes"`
	PageInfo     *PageInfo           `json:"pageInfo"`
	TotalCount   int                 `json:"totalCount"`
}

type BasicMessageEdge struct {
	Cursor string        `json:"cursor"`
	Node   *BasicMessage `json:"node"`
}

type ConnectInput struct {
	Invitation string `json:"invitation"`
}

type Credential struct {
	ID            string             `json:"id"`
	Role          CredentialRole     `json:"role"`
	SchemaID      string             `json:"schemaId"`
	CredDefID     string             `json:"credDefId"`
	Attributes    []*CredentialValue `json:"attributes"`
	InitiatedByUs bool               `json:"initiatedByUs"`
	ApprovedMs    *string            `json:"approvedMs"`
	IssuedMs      *string            `json:"issuedMs"`
	CreatedMs     string             `json:"createdMs"`
	Connection    *Pairwise          `json:"connection"`
}

type CredentialConnection struct {
	ConnectionID *string           `json:"connectionId"`
	Edges        []*CredentialEdge `json:"edges"`
	Nodes        []*Credential     `json:"nodes"`
	PageInfo     *PageInfo         `json:"pageInfo"`
	TotalCount   int               `json:"totalCount"`
}

type CredentialEdge struct {
	Cursor string      `json:"cursor"`
	Node   *Credential `json:"node"`
}

type CredentialMatch struct {
	ID           string `json:"id"`
	CredentialID string `json:"credentialId"`
	Value        string `json:"value"`
}

type CredentialValue struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Event struct {
	ID          string    `json:"id"`
	Read        bool      `json:"read"`
	Description string    `json:"description"`
	CreatedMs   string    `json:"createdMs"`
	Job         *JobEdge  `json:"job"`
	Connection  *Pairwise `json:"connection"`
}

type EventConnection struct {
	ConnectionID *string      `json:"connectionId"`
	Edges        []*EventEdge `json:"edges"`
	Nodes        []*Event     `json:"nodes"`
	PageInfo     *PageInfo    `json:"pageInfo"`
	TotalCount   int          `json:"totalCount"`
}

type EventEdge struct {
	Cursor string `json:"cursor"`
	Node   *Event `json:"node"`
}

type InvitationResponse struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	Endpoint string `json:"endpoint"`
	Raw      string `json:"raw"`
	ImageB64 string `json:"imageB64"`
}

type Job struct {
	ID            string       `json:"id"`
	Protocol      ProtocolType `json:"protocol"`
	InitiatedByUs bool         `json:"initiatedByUs"`
	Status        JobStatus    `json:"status"`
	Result        JobResult    `json:"result"`
	CreatedMs     string       `json:"createdMs"`
	UpdatedMs     string       `json:"updatedMs"`
	Output        *JobOutput   `json:"output"`
}

type JobConnection struct {
	ConnectionID *string    `json:"connectionId"`
	Completed    *bool      `json:"completed"`
	Edges        []*JobEdge `json:"edges"`
	Nodes        []*Job     `json:"nodes"`
	PageInfo     *PageInfo  `json:"pageInfo"`
	TotalCount   int        `json:"totalCount"`
}

type JobEdge struct {
	Cursor string `json:"cursor"`
	Node   *Job   `json:"node"`
}

type JobOutput struct {
	Connection *PairwiseEdge     `json:"connection"`
	Message    *BasicMessageEdge `json:"message"`
	Credential *CredentialEdge   `json:"credential"`
	Proof      *ProofEdge        `json:"proof"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type MarkReadInput struct {
	ID string `json:"id"`
}

type MessageInput struct {
	ConnectionID string `json:"connectionId"`
	Message      string `json:"message"`
}

type PageInfo struct {
	EndCursor       *string `json:"endCursor"`
	HasNextPage     bool    `json:"hasNextPage"`
	HasPreviousPage bool    `json:"hasPreviousPage"`
	StartCursor     *string `json:"startCursor"`
}

type Pairwise struct {
	ID            string                  `json:"id"`
	OurDid        string                  `json:"ourDid"`
	TheirDid      string                  `json:"theirDid"`
	TheirEndpoint string                  `json:"theirEndpoint"`
	TheirLabel    string                  `json:"theirLabel"`
	CreatedMs     string                  `json:"createdMs"`
	ApprovedMs    string                  `json:"approvedMs"`
	Invited       bool                    `json:"invited"`
	Messages      *BasicMessageConnection `json:"messages"`
	Credentials   *CredentialConnection   `json:"credentials"`
	Proofs        *ProofConnection        `json:"proofs"`
	Jobs          *JobConnection          `json:"jobs"`
	Events        *EventConnection        `json:"events"`
}

type PairwiseConnection struct {
	Edges      []*PairwiseEdge `json:"edges"`
	Nodes      []*Pairwise     `json:"nodes"`
	PageInfo   *PageInfo       `json:"pageInfo"`
	TotalCount int             `json:"totalCount"`
}

type PairwiseEdge struct {
	Cursor string    `json:"cursor"`
	Node   *Pairwise `json:"node"`
}

type Proof struct {
	ID            string            `json:"id"`
	Role          ProofRole         `json:"role"`
	Attributes    []*ProofAttribute `json:"attributes"`
	Values        []*ProofValue     `json:"values"`
	Provable      *Provable         `json:"provable"`
	InitiatedByUs bool              `json:"initiatedByUs"`
	Result        bool              `json:"result"`
	VerifiedMs    *string           `json:"verifiedMs"`
	ApprovedMs    *string           `json:"approvedMs"`
	CreatedMs     string            `json:"createdMs"`
	Connection    *Pairwise         `json:"connection"`
}

type ProofAttribute struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CredDefID string `json:"credDefId"`
}

type ProofConnection struct {
	ConnectionID *string      `json:"connectionId"`
	Edges        []*ProofEdge `json:"edges"`
	Nodes        []*Proof     `json:"nodes"`
	PageInfo     *PageInfo    `json:"pageInfo"`
	TotalCount   int          `json:"totalCount"`
}

type ProofEdge struct {
	Cursor string `json:"cursor"`
	Node   *Proof `json:"node"`
}

type ProofRequestAttribute struct {
	Name      string `json:"name"`
	CredDefID string `json:"credDefId"`
}

type ProofRequestInput struct {
	ConnectionID string                   `json:"connectionId"`
	Attributes   []*ProofRequestAttribute `json:"attributes"`
}

type ProofValue struct {
	ID          string `json:"id"`
	AttributeID string `json:"attributeId"`
	Value       string `json:"value"`
}

type Provable struct {
	ID         string               `json:"id"`
	Provable   bool                 `json:"provable"`
	Attributes []*ProvableAttribute `json:"attributes"`
}

type ProvableAttribute struct {
	ID          string             `json:"id"`
	Attribute   *ProofAttribute    `json:"attribute"`
	Credentials []*CredentialMatch `json:"credentials"`
}

type Response struct {
	Ok bool `json:"ok"`
}

type ResumeJobInput struct {
	ID     string `json:"id"`
	Accept bool   `json:"accept"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CredentialRole string

const (
	CredentialRoleIssuer CredentialRole = "ISSUER"
	CredentialRoleHolder CredentialRole = "HOLDER"
)

var AllCredentialRole = []CredentialRole{
	CredentialRoleIssuer,
	CredentialRoleHolder,
}

func (e CredentialRole) IsValid() bool {
	switch e {
	case CredentialRoleIssuer, CredentialRoleHolder:
		return true
	}
	return false
}

func (e CredentialRole) String() string {
	return string(e)
}

func (e *CredentialRole) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = CredentialRole(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid CredentialRole", str)
	}
	return nil
}

func (e CredentialRole) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type JobResult string

const (
	JobResultNone    JobResult = "NONE"
	JobResultSuccess JobResult = "SUCCESS"
	JobResultFailure JobResult = "FAILURE"
)

var AllJobResult = []JobResult{
	JobResultNone,
	JobResultSuccess,
	JobResultFailure,
}

func (e JobResult) IsValid() bool {
	switch e {
	case JobResultNone, JobResultSuccess, JobResultFailure:
		return true
	}
	return false
}

func (e JobResult) String() string {
	return string(e)
}

func (e *JobResult) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = JobResult(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid JobResult", str)
	}
	return nil
}

func (e JobResult) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type JobStatus string

const (
	JobStatusWaiting  JobStatus = "WAITING"
	JobStatusPending  JobStatus = "PENDING"
	JobStatusBlocked  JobStatus = "BLOCKED"
	JobStatusComplete JobStatus = "COMPLETE"
)

var AllJobStatus = []JobStatus{
	JobStatusWaiting,
	JobStatusPending,
	JobStatusBlocked,
	JobStatusComplete,
}

func (e JobStatus) IsValid() bool {
	switch e {
	case JobStatusWaiting, JobStatusPending, JobStatusBlocked, JobStatusComplete:
		return true
	}
	return false
}

func (e JobStatus) String() string {
	return string(e)
}

func (e *JobStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = JobStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid JobStatus", str)
	}
	return nil
}

func (e JobStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ProofRole string

const (
	ProofRoleVerifier ProofRole = "VERIFIER"
	ProofRoleProver   ProofRole = "PROVER"
)

var AllProofRole = []ProofRole{
	ProofRoleVerifier,
	ProofRoleProver,
}

func (e ProofRole) IsValid() bool {
	switch e {
	case ProofRoleVerifier, ProofRoleProver:
		return true
	}
	return false
}

func (e ProofRole) String() string {
	return string(e)
}

func (e *ProofRole) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ProofRole(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ProofRole", str)
	}
	return nil
}

func (e ProofRole) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ProtocolType string

const (
	ProtocolTypeNone         ProtocolType = "NONE"
	ProtocolTypeConnection   ProtocolType = "CONNECTION"
	ProtocolTypeCredential   ProtocolType = "CREDENTIAL"
	ProtocolTypeProof        ProtocolType = "PROOF"
	ProtocolTypeBasicMessage ProtocolType = "BASIC_MESSAGE"
)

var AllProtocolType = []ProtocolType{
	ProtocolTypeNone,
	ProtocolTypeConnection,
	ProtocolTypeCredential,
	ProtocolTypeProof,
	ProtocolTypeBasicMessage,
}

func (e ProtocolType) IsValid() bool {
	switch e {
	case ProtocolTypeNone, ProtocolTypeConnection, ProtocolTypeCredential, ProtocolTypeProof, ProtocolTypeBasicMessage:
		return true
	}
	return false
}

func (e ProtocolType) String() string {
	return string(e)
}

func (e *ProtocolType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ProtocolType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ProtocolType", str)
	}
	return nil
}

func (e ProtocolType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
