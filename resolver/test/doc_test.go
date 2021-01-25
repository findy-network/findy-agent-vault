package test

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
	"github.com/findy-network/findy-agent-vault/resolver"
	"github.com/findy-network/findy-agent-vault/server"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/findy-network/findy-grpc/jwt"
)

var (
	r                *resolver.Resolver
	testConnectionID string
	testCredentialID string
	testProofID      string
	testMessageID    string
	testEventID      string
	testJobID        string
	totalCount       = 5
)

func testContext() context.Context {
	const testValidationKey = "test-secret"
	uToken := server.NewServer(nil, testValidationKey).CreateTestToken(testValidationKey)
	ctx := jwt.TokenToContext(context.Background(), "user", &jwt.Token{Raw: uToken})

	return ctx
}

func setup() {
	utils.SetLogDefaults()

	r = resolver.InitResolver(&utils.Configuration{UseMockDB: true, UseMockAgency: true})
	db := r.Store()

	size := totalCount
	a, c := test.AddAgentAndConnections(db, fake.FakeCloudDID, size)
	testConnectionID = c[0].ID

	cr := fake.AddCredentials(db, a.ID, c[0].ID, size)
	testCredentialID = cr[0].ID

	pr := fake.AddProofs(db, a.ID, c[0].ID, size)
	testProofID = pr[0].ID

	msg := fake.AddMessages(db, a.ID, c[0].ID, size)
	testMessageID = msg[0].ID

	jb := fake.AddCredentialJobs(db, a.ID, c[0].ID, testCredentialID, size)
	testJobID = jb[0].ID

	ev := fake.AddEvents(db, a.ID, c[0].ID, &testJobID, size)
	testEventID = ev[0].ID
}

func teardown() {
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

type executor func(ctx context.Context, after *string, before *string, first *int, last *int) error

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
