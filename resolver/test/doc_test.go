package test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/findy-network/findy-agent-vault/agency/mock"
	"github.com/findy-network/findy-agent-vault/db/fake"
	"github.com/findy-network/findy-agent-vault/db/store"
	"github.com/findy-network/findy-agent-vault/db/store/pg"
	"github.com/findy-network/findy-agent-vault/db/store/test"
	"github.com/findy-network/findy-agent-vault/paginator"
	"github.com/findy-network/findy-agent-vault/resolver"
	"github.com/findy-network/findy-agent-vault/server"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/findy-network/findy-common-go/jwt"
	"github.com/golang/mock/gomock"
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

	config = &utils.Configuration{
		DBHost:           "localhost",
		DBPassword:       os.Getenv("FAV_DB_PASSWORD"),
		DBPort:           5433,
		DBMigrationsPath: "file://../../db/migrations",
		DBName:           "resolver",
	}
	resolverDB store.DB
)

func testContextForUser(userName string) context.Context {
	const testValidationKey = "test-secret"
	uToken := server.NewServer(nil, testValidationKey).CreateTestToken(userName, testValidationKey)
	ctx := jwt.TokenToContext(context.Background(), "user", &jwt.Token{Raw: uToken})

	return ctx
}

func testContext() context.Context {
	return testContextForUser(fake.FakeCloudDID)
}

func setup() {
	utils.SetLogDefaults()

	resolverDB = pg.InitDB(config, true, true)
}

func teardown() {
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func beforeEachWithID(t *testing.T, id string) (m *mock.MockAgency) {
	ctrl := gomock.NewController(t)

	m = mock.NewMockAgency(ctrl)

	m.
		EXPECT().
		Init(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())

	m.EXPECT().AddAgent(gomock.Any()).AnyTimes()

	r = resolver.InitResolverWithDB(config, m, resolverDB)
	db := r.Store()

	size := totalCount
	a, c := test.AddAgentAndConnections(db, id, size)
	testConnectionID = c[0].ID

	cr := fake.AddCredentials(db, a.ID, c[0].ID, size)
	testCredentialID = cr[0].ID

	pr := fake.AddProofs(db, a.ID, c[0].ID, size, true)
	testProofID = pr[0].ID

	msg := fake.AddMessages(db, a.ID, c[0].ID, size)
	testMessageID = msg[0].ID

	jb := fake.AddCredentialJobs(db, a.ID, c[0].ID, testCredentialID, size)
	testJobID = jb[0].ID

	ev := fake.AddEvents(db, a.ID, c[0].ID, &testJobID, size)
	testEventID = ev[0].ID

	return m
}

func beforeEach(t *testing.T) (m *mock.MockAgency) {
	return beforeEachWithID(t, fake.FakeCloudDID)
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
