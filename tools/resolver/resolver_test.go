package resolver

import (
	"context"
	"os"
	"testing"

	"github.com/findy-network/findy-agent-vault/tools/utils"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

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
		Invitation: "",
	})
	if err != nil {
		panic("Connect request failed")
	}
}

func teardown() {
}
