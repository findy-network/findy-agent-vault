package resolver

import (
	"context"
	"flag"
	"os"
	"testing"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	flag.Parse()

	InitResolver(nil)
	r := Resolver{}

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
