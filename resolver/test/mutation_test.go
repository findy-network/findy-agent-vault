package test

import (
	"testing"

	agency "github.com/findy-network/findy-agent-vault/agency/model"
	"github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/golang/mock/gomock"
)

const (
	testInvitation = `{"serviceEndpoint":` +
		`"http://url",` +
		`"recipientKeys":["Hmk4756ry7fqBCKPf634SRvaM3xss1QBhoFC1uAbwkVL"],"@id":"d679e4c6-b8db-4c39-99ca-783034b51bd4"` +
		`,"label":"findy-issuer","@type":"did:sov:BzCbsNYhMrjHiqZDTUASHg;spec/connections/1.0/invitation"}`
)

func TestMarkEventRead(t *testing.T) {
	beforeEach(t)

	event, err := r.Mutation().MarkEventRead(testContext(), model.MarkReadInput{ID: testEventID})
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if event == nil {
		t.Errorf("Expecting result, received %v", event)
	}
}

func TestInvite(t *testing.T) {
	const user = "TestInvite"
	m := beforeEachWithID(t, user)

	data := agency.InvitationData{
		ID:  "d679e4c6-b8db-4c39-99ca-783034b51bd4",
		Raw: "didcomm://aries_connection_invitation?c_i=eyJzZXJ2aWNlRW5kcG9pbnQiOiJodHRwOi8vbG9jYWxob3N0OjgwODAvYTJhLzNKY3NUYW9tR2NQdlRBMmlpdDdSZ2YvM0pjc1Rhb21HY1B2VEEyaWl0N1JnZi9MVkRZa1ZMZXpyREhUa0VGSG1vQ216L2Q2OTAzZjIxLTFjOTItNDVkMi04MmFkLTM0ZTg5NGJmYjAxYiIsInJlY2lwaWVudEtleXMiOlsiQmQxTkZnSDlQeW5MakJzSmlqNFFIc2U5WXlxbjJTcFB1RHNSaVBVZFdveXUiXSwiQGlkIjoiZDY5MDNmMjEtMWM5Mi00NWQyLTgyYWQtMzRlODk0YmZiMDFiIiwibGFiZWwiOiJlbXB0eS1sYWJlbCIsIkB0eXBlIjoiZGlkOnNvdjpCekNic05ZaE1yakhpcVpEVFVBU0hnO3NwZWMvY29ubmVjdGlvbnMvMS4wL2ludml0YXRpb24ifQ",
	}

	m.
		EXPECT().
		Invite(gomock.Any()).Return(&data, nil)

	resp, err := r.Mutation().Invite(testContextForUser(user))
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if resp == nil {
		t.Errorf("Expecting result, received %v", resp)
	}
}

func TestConnect(t *testing.T) {
	const user = "TestConnect"
	m := beforeEachWithID(t, user)

	m.
		EXPECT().
		Connect(gomock.Any(), gomock.Any()).
		Return("d679e4c6-b8db-4c39-99ca-783034b51bd4", nil)

	resp, err := r.Mutation().Connect(testContextForUser(user), model.ConnectInput{Invitation: testInvitation})
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if resp == nil {
		t.Errorf("Expecting result, received %v", resp)
	}
}

func TestSendMessage(t *testing.T) {
	m := beforeEach(t)

	m.
		EXPECT().
		SendMessage(gomock.Any(), gomock.Any(), gomock.Any())

	resp, err := r.Mutation().SendMessage(testContext(), model.MessageInput{})
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if resp == nil {
		t.Errorf("Expecting result, received %v", resp)
	}
}

func TestResume(t *testing.T) {
	m := beforeEach(t)

	m.
		EXPECT().
		ResumeCredentialOffer(gomock.Any(), gomock.Any(), gomock.Any())

	resp, err := r.Mutation().Resume(testContext(), model.ResumeJobInput{ID: testJobID})
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if resp == nil {
		t.Errorf("Expecting result, received %v", resp)
	}
}
