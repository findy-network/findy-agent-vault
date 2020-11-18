package model

import (
	"encoding/base64"
	"reflect"
	"strconv"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

func CreateCursor(created int64, object interface{}) string {
	typeName := reflect.TypeOf(object).Name()
	return base64.StdEncoding.EncodeToString([]byte(typeName + ":" + strconv.FormatInt(created, 10)))
}

type ProtocolStatus struct {
	Status      model.JobStatus
	Result      model.JobResult
	Description string
}

type APIObject interface {
	Identifier() string
	Created() int64
	Pairwise() *InternalPairwise
	BasicMessage() *InternalMessage
	Credential() *InternalCredential
	Proof() *InternalProof
	Event() *InternalEvent
	Job() *InternalJob
}

type BaseObject struct {
	ID        string `faker:"uuid_hyphenated"`
	CreatedMs int64  `faker:"created"`
}

func (o *BaseObject) Created() int64 {
	return o.CreatedMs
}

func (o *BaseObject) Identifier() string {
	return o.ID
}

func (o *BaseObject) Pairwise() *InternalPairwise {
	panic("Object is not pairwise")
}

func (o *BaseObject) Credential() *InternalCredential {
	panic("Object is not credential")
}

func (o *BaseObject) BasicMessage() *InternalMessage {
	panic("Object is not credential")
}

func (o *BaseObject) Proof() *InternalProof {
	panic("Object is not proof")
}

func (o *BaseObject) Event() *InternalEvent {
	panic("Object is not event")
}

func (o *BaseObject) Job() *InternalJob {
	panic("Object is not job")
}
