// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type ConnectInput struct {
	Invitation string `json:"invitation"`
}

type Event struct {
	ID          string    `json:"id"`
	Read        bool      `json:"read"`
	Description string    `json:"description"`
	CreatedMs   string    `json:"createdMs"`
	Job         *Job      `json:"job"`
	Connection  *Pairwise `json:"connection"`
}

type EventConnection struct {
	Edges      []*EventEdge `json:"edges"`
	Nodes      []*Event     `json:"nodes"`
	PageInfo   *PageInfo    `json:"pageInfo"`
	TotalCount int          `json:"totalCount"`
}

type EventEdge struct {
	Cursor string `json:"cursor"`
	Node   *Event `json:"node"`
}

type InvitationResponse struct {
	Invitation string `json:"invitation"`
	ImageB64   string `json:"imageB64"`
}

type Job struct {
	ID            string       `json:"id"`
	Protocol      ProtocolType `json:"protocol"`
	ProtocolID    *string      `json:"protocolId"`
	InitiatedByUs bool         `json:"initiatedByUs"`
	Connection    *Pairwise    `json:"connection"`
	Status        JobStatus    `json:"status"`
	Result        JobResult    `json:"result"`
	CreatedMs     string       `json:"createdMs"`
	UpdatedMs     string       `json:"updatedMs"`
}

type JobConnection struct {
	Edges      []*JobEdge `json:"edges"`
	Nodes      []*Job     `json:"nodes"`
	PageInfo   *PageInfo  `json:"pageInfo"`
	TotalCount int        `json:"totalCount"`
}

type JobEdge struct {
	Cursor string `json:"cursor"`
	Node   *Job   `json:"node"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type MarkReadInput struct {
	ID string `json:"id"`
}

type Offer struct {
	ID     string `json:"id"`
	Accept bool   `json:"accept"`
}

type PageInfo struct {
	EndCursor       *string `json:"endCursor"`
	HasNextPage     bool    `json:"hasNextPage"`
	HasPreviousPage bool    `json:"hasPreviousPage"`
	StartCursor     *string `json:"startCursor"`
}

type Pairwise struct {
	ID            string `json:"id"`
	OurDid        string `json:"ourDid"`
	TheirDid      string `json:"theirDid"`
	TheirEndpoint string `json:"theirEndpoint"`
	TheirLabel    string `json:"theirLabel"`
	CreatedMs     string `json:"createdMs"`
	ApprovedMs    string `json:"approvedMs"`
	InitiatedByUs bool   `json:"initiatedByUs"`
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

type Request struct {
	ID     string `json:"id"`
	Accept bool   `json:"accept"`
}

type Response struct {
	Ok bool `json:"ok"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
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
	JobStatusComplete JobStatus = "COMPLETE"
)

var AllJobStatus = []JobStatus{
	JobStatusWaiting,
	JobStatusPending,
	JobStatusComplete,
}

func (e JobStatus) IsValid() bool {
	switch e {
	case JobStatusWaiting, JobStatusPending, JobStatusComplete:
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
