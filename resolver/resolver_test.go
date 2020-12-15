package resolver

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/findy-network/findy-agent-vault/db/fake"

	"github.com/findy-network/findy-agent-vault/db/store/test"
	"github.com/findy-network/findy-agent-vault/paginator"
	"github.com/findy-network/findy-agent-vault/server"
	"github.com/findy-network/findy-agent-vault/utils"
)

type executor func(ctx context.Context, after *string, before *string, first *int, last *int) error

func testContext() context.Context {
	u := server.CreateTestToken("test")
	ctx := context.WithValue(context.Background(), "user", u)
	return ctx
}

func testPaginationErrors(t *testing.T, objName string, ex executor) {
	t.Run(fmt.Sprintf("get %s", objName), func(t *testing.T) {
		var (
			valid              = 1
			tooLow             = 0
			tooHigh            = 101
			invalidCursor      = "1"
			missingError       = errors.New(paginator.ErrorFirstLastMissing)
			invalidCountError  = errors.New(paginator.ErrorFirstLastInvalid)
			invalidCursorError = errors.New(paginator.ErrorCursorInvalid)
		)
		tests := []struct {
			name string
			args paginator.Params
			err  error
		}{
			{fmt.Sprintf("%s, pagination missing", objName), paginator.Params{}, missingError},
			{fmt.Sprintf("%s, pagination first too low", objName), paginator.Params{First: &tooLow}, invalidCountError},
			{fmt.Sprintf("%s, pagination first too high", objName), paginator.Params{First: &tooHigh}, invalidCountError},
			{fmt.Sprintf("%s, pagination last too low", objName), paginator.Params{Last: &tooLow}, invalidCountError},
			{fmt.Sprintf("%s, pagination last too high", objName), paginator.Params{Last: &tooHigh}, invalidCountError},
			{fmt.Sprintf("%s, after cursor invalid", objName), paginator.Params{First: &valid, After: &invalidCursor}, invalidCursorError},
		}

		for _, testCase := range tests {
			tc := testCase
			t.Run(tc.name, func(t *testing.T) {
				err := ex(testContext(), tc.args.After, tc.args.Before, tc.args.First, tc.args.Last)
				if !reflect.DeepEqual(err, tc.err) {
					t.Errorf("%s = err (%v)\n want (%v)", tc.name, err, tc.err)
				}
			})
		}
	})
}

func TestPaginationErrorsGetConnections(t *testing.T) {
	testPaginationErrors(t, "connections", func(ctx context.Context, after, before *string, first, last *int) error {
		r := InitResolver(true)
		_, err := r.Query().Connections(ctx, after, before, first, last)
		return err
	})
}

func TestResolverGetConnections(t *testing.T) {
	first := 1
	c, err := r.Query().Connections(testContext(), nil, nil, &first, nil)
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if c == nil || len(c.Edges) == 0 {
		t.Errorf("Expecting result, received %v", c)
	}
}

func TestGetConnection(t *testing.T) {
	c, err := r.Query().Connection(testContext(), testConnectionID)
	if err != nil {
		t.Errorf("Received unexpected error %s", err)
	}
	if c == nil {
		t.Errorf("Expecting result, received %v", c)
	}
}

var (
	r                *Resolver
	testConnectionID string
)

func setup() {
	utils.SetLogDefaults()
	r = InitResolver(true)
	size := 5
	_, c := test.AddAgentAndConnections(r.db, fake.FakeCloudDID, size)
	testConnectionID = c[0].ID
}

func teardown() {
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}
