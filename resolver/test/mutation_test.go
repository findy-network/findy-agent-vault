package test

import (
	"testing"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

const (
	testInvitation = `{"serviceEndpoint":` +
		`"http://findy-agent.op-ai.fi/a2a/Xmjk7cFr8TT2j5kWLWyhDB/Xmjk7cFr8TT2j5kWLWyhDB/GqmnSTxevze48yio5m2fUE",` +
		`"recipientKeys":["Hmk4756ry7fqBCKPf634SRvaM3xss1QBhoFC1uAbwkVL"],"@id":"d679e4c6-b8db-4c39-99ca-783034b51bd4"` +
		`,"label":"findy-issuer","@type":"did:sov:BzCbsNYhMrjHiqZDTUASHg;spec/connections/1.0/invitation"}`
)

func TestMarkEventRead(t *testing.T) {
	event, err := r.Mutation().MarkEventRead(testContext(), model.MarkReadInput{ID: testEventID})
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if event == nil {
		t.Errorf("Expecting result, received %v", event)
	}
}

func TestInvite(t *testing.T) {
	resp, err := r.Mutation().Invite(testContext())
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if resp == nil {
		t.Errorf("Expecting result, received %v", resp)
	}
}

func TestConnect(t *testing.T) {
	resp, err := r.Mutation().Connect(testContext(), model.ConnectInput{Invitation: testInvitation})
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if resp == nil {
		t.Errorf("Expecting result, received %v", resp)
	}
}

func TestSendMessage(t *testing.T) {
	resp, err := r.Mutation().SendMessage(testContext(), model.MessageInput{})
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if resp == nil {
		t.Errorf("Expecting result, received %v", resp)
	}
}

func TestResume(t *testing.T) {
	resp, err := r.Mutation().Resume(testContext(), model.ResumeJobInput{ID: testJobID})
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if resp == nil {
		t.Errorf("Expecting result, received %v", resp)
	}
}
