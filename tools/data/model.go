package data

import (
	"encoding/base64"
	"reflect"
	"strconv"

	"github.com/findy-network/findy-agent-api/graph/model"
)

func CreateCursor(created int64, object interface{}) string {
	typeName := reflect.TypeOf(object).Name()
	return base64.StdEncoding.EncodeToString([]byte(typeName + ":" + strconv.FormatInt(created, 10)))
}

type APIObject interface {
	Identifier() string
	Created() int64
	Pairwise() *InternalPairwise
	Event() *InternalEvent
}

type InternalPairwise struct {
	ID            string `faker:"uuid_hyphenated"`
	OurDid        string
	TheirDid      string
	TheirEndpoint string `faker:"url"`
	TheirLabel    string `faker:"organisationLabel"`
	InitiatedByUs bool
	ApprovedMs    int64 `faker:"unix_time"`
	CreatedMs     int64 `faker:"unix_time"`
}

func (p *InternalPairwise) Created() int64 {
	return p.CreatedMs
}

func (p *InternalPairwise) Identifier() string {
	return p.ID
}

func (p *InternalPairwise) Pairwise() *InternalPairwise {
	return p
}

func (p *InternalPairwise) Event() *InternalEvent {
	panic("Pairwise is not event")
}

func (p *InternalPairwise) ToNode() *model.Pairwise {
	return &model.Pairwise{
		ID:            p.ID,
		OurDid:        p.OurDid,
		TheirDid:      p.TheirDid,
		TheirEndpoint: p.TheirEndpoint,
		TheirLabel:    p.TheirLabel,
		CreatedMs:     strconv.FormatInt(p.CreatedMs, 10),
		ApprovedMs:    strconv.FormatInt(p.ApprovedMs, 10),
		InitiatedByUs: p.InitiatedByUs,
	}
}

type InternalEvent struct {
	ID           string             `faker:"uuid_hyphenated"`
	Read         bool               `faker:"-"`
	Description  string             `faker:"sentence"`
	ProtocolType model.ProtocolType `faker:"oneof: model.ProtocolTypeNone, model.ProtocolTypeConnection, model.ProtocolTypeCredential, model.ProtocolTypeProof, model.ProtocolTypeBasicMessage"`
	Type         model.EventType    `faker:"oneof: model.EventTypeNotification, model.EventTypeAction"`
	PairwiseID   string             `faker:"eventPairwiseId"`
	CreatedMs    int64              `faker:"unix_time"`
}

func (e *InternalEvent) Created() int64 {
	return e.CreatedMs
}

func (e *InternalEvent) Identifier() string {
	return e.ID
}

func (e *InternalEvent) Pairwise() *InternalPairwise {
	panic("Event is not pairwise")
}

func (e *InternalEvent) Event() *InternalEvent {
	return e
}

func (e *InternalEvent) ToEdge() *model.EventEdge {
	cursor := CreateCursor(e.CreatedMs, model.Event{})
	return &model.EventEdge{
		Cursor: cursor,
		Node:   e.ToNode(),
	}
}

func (e *InternalEvent) ToNode() *model.Event {
	createdStr := strconv.FormatInt(e.CreatedMs, 10)
	var node *model.Pairwise
	for _, c := range State.Connections.items {
		if c.Identifier() == e.PairwiseID {
			node = c.Pairwise().ToNode()
			break
		}
	}
	return &model.Event{
		ID:          e.ID,
		Read:        e.Read,
		Description: e.Description,
		Protocol:    e.ProtocolType,
		Type:        e.Type,
		CreatedMs:   createdStr,
		Connection:  node,
	}
}

type InternalUser struct {
	ID   string `faker:"uuid_hyphenated"`
	Name string `faker:"first_name"`
}

func (u *InternalUser) ToNode() *model.User {
	return &model.User{
		ID:   u.ID,
		Name: u.Name,
	}
}
