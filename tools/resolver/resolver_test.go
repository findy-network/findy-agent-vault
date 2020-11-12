package resolver

import (
	"context"
	"os"
	"testing"

	"github.com/findy-network/findy-agent-vault/tools/utils"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

const invitation = "{\"serviceEndpoint\":\"https://www.ufwTCAB.info/ZUavCJk\"," +
	"\"recipientKeys\":[\"CDdVp7CyP9Ued38FpFd8rqxF3eEKhrnjAsPWf6LEeLJC\"]," +
	"\"@id\":\"5c103f67-7b46-4561-972f-34a8047bad96\",\"label\":\"Sabina\"," +
	"\"@type\":\"did:sov:BzCbsNYhMrjHiqZDTUASHg;spec/connections/1.0/invitation\"}"

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	utils.SetLogDefaults()

	r := InitResolver()

	// Generate some jobs data
	_, err := r.Mutation().Invite(context.TODO())
	if err != nil {
		panic("Invitation failed")
	}

	_, err = r.Mutation().Connect(context.TODO(), model.ConnectInput{
		Invitation: invitation,
	})
	if err != nil {
		panic("Connect request failed")
	}
}

func teardown() {
}
