package resolver

import (
	"context"
	"os"
	"testing"

	"github.com/findy-network/findy-agent-vault/utils"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

func firstPairwise() *model.Pairwise {
	return state.Connections().PairwiseConnection(0, 1).Nodes[0]
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func resetResolver(skipFake bool) {
	r := InitResolver(skipFake)

	// Generate some jobs data
	res, err := r.Mutation().Invite(context.TODO())
	if err != nil {
		panic("Invitation failed")
	}

	_, err = r.Mutation().Connect(context.TODO(), model.ConnectInput{
		Invitation: res.Invitation,
	})
	if err != nil {
		panic("Connect request failed")
	}
}

func setup() {
	utils.SetLogDefaults()

	resetResolver(false)
}

func teardown() {
}
