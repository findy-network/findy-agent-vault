package resolver

import (
	"context"
	"reflect"
	"testing"

	"github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/tools/utils"
	"github.com/google/uuid"
)

type CredentialTestRes struct {
	Credential   *model.Credential
	Error        error
	ExpectsError bool
}

func checkCredential(r *Resolver, pw *model.Pairwise, t *testing.T, name string, got, expected *CredentialTestRes) {
	if got.Error == nil && !expected.ExpectsError {
		if expected.Credential.ID != got.Credential.ID ||
			expected.Credential.Role != got.Credential.Role ||
			expected.Credential.SchemaID != got.Credential.SchemaID ||
			expected.Credential.CredDefID != got.Credential.CredDefID ||
			expected.Credential.InitiatedByUs != got.Credential.InitiatedByUs ||
			!reflect.DeepEqual(expected.Credential.Attributes, got.Credential.Attributes) ||
			got.Credential.CreatedMs == "" {
			t.Errorf("%s = %v, want %v", name, got.Credential, expected.Credential)
		}
		gotConn, err := r.Credential().Connection(context.TODO(), got.Credential)
		if err != nil || !reflect.DeepEqual(gotConn, pw) {
			t.Errorf("%s = %v, want %v", "get credential connection", got, pw)
		}
	} else if got.Error != nil && !expected.ExpectsError {
		t.Errorf("%s failed, expecting for error", name)
	}
}

func TestGetCredential(t *testing.T) {
	t.Run("get credential", func(t *testing.T) {
		// add new credential
		listener := &agencyListener{}
		id := uuid.New().String()
		pw := firstPairwise()
		credential := &model.Credential{
			ID:            id,
			Role:          model.CredentialRoleHolder,
			SchemaID:      "schemaID",
			CredDefID:     "credDefID",
			Attributes:    []*model.CredentialValue{{Name: "email", Value: "emailValue"}},
			InitiatedByUs: false,
		}
		listener.AddCredential(
			pw.ID,
			credential.ID,
			credential.Role,
			credential.SchemaID,
			credential.CredDefID,
			credential.Attributes,
			credential.InitiatedByUs)

		tests := []struct {
			name   string
			ID     string
			result *CredentialTestRes
		}{
			{"added message", id, &CredentialTestRes{Credential: credential}},
			{"invalid connection", "", &CredentialTestRes{ExpectsError: true}},
		}

		r := Resolver{}
		for _, testCase := range tests {
			tc := testCase
			t.Run(tc.name, func(t *testing.T) {
				got, err := r.Query().Credential(context.TODO(), tc.ID)
				checkCredential(&r, pw, t, tc.name, &CredentialTestRes{got, err, false}, tc.result)
			})
		}
	})
}

func TestPaginationErrorsGetCredentials(t *testing.T) {
	testPaginationErrors(t, "credentials", func(ctx context.Context, after, before *string, first, last *int) error {
		r := Resolver{}
		pw := firstPairwise()
		_, err := r.Pairwise().Credentials(context.TODO(), pw, after, before, first, last)
		return err
	})
}

func TestGetCredentialsForPairwise(t *testing.T) {
	r := Resolver{}

	// add new connection
	listener := &agencyListener{}
	currentTime := utils.CurrentTimeMs()

	connID := uuid.New().String()
	listener.AddConnection(connID, "ourDID", "theirDID", "theirEndpoint", "theirLabel")
	conn, _ := r.Query().Connection(context.TODO(), connID)

	id := uuid.New().String()
	credential := &model.Credential{
		ID:            id,
		Role:          model.CredentialRoleHolder,
		SchemaID:      "schemaID",
		CredDefID:     "credDefID",
		Attributes:    []*model.CredentialValue{{Name: "email", Value: "emailValue"}},
		InitiatedByUs: false,
	}
	listener.AddCredential(
		conn.ID,
		credential.ID,
		credential.Role,
		credential.SchemaID,
		credential.CredDefID,
		credential.Attributes,
		credential.InitiatedByUs)
	listener.UpdateCredential(connID, credential.ID, &currentTime, &currentTime, nil)

	// add incomplete credential
	listener.AddCredential(
		connID,
		uuid.New().String(),
		credential.Role,
		credential.SchemaID,
		credential.CredDefID,
		credential.Attributes,
		credential.InitiatedByUs)

	t.Run("get credentials", func(t *testing.T) {
		var (
			valid = 1
		)
		tests := []struct {
			name   string
			args   PaginationParams
			result *CredentialTestRes
		}{
			{"first credential", PaginationParams{first: &valid}, &CredentialTestRes{Credential: credential}},
			{"last credential", PaginationParams{last: &valid}, &CredentialTestRes{Credential: credential}},
		}

		for _, testCase := range tests {
			tc := testCase
			t.Run(tc.name, func(t *testing.T) {
				got, err := r.Pairwise().Credentials(context.TODO(), conn, tc.args.after, tc.args.before, tc.args.first, tc.args.last)
				checkCredential(&r, conn, t, tc.name, &CredentialTestRes{got.Nodes[0], err, false}, tc.result)
			})
		}
	})
}

func TestGetAllCredentials(t *testing.T) {
	resetResolver(true)
	r := Resolver{}

	// add new connections
	listener := &agencyListener{}
	currentTime := utils.CurrentTimeMs()

	connID := uuid.New().String()
	listener.AddConnection(connID, "ourDID", "theirDID", "theirEndpoint", "theirLabel")
	conn, _ := r.Query().Connection(context.TODO(), connID)

	connID2 := uuid.New().String()
	listener.AddConnection(connID2, "ourDID", "theirDID", "theirEndpoint", "theirLabel")
	conn2, _ := r.Query().Connection(context.TODO(), connID2)

	id := uuid.New().String()
	credential := &model.Credential{
		ID:            id,
		Role:          model.CredentialRoleHolder,
		SchemaID:      "schemaID",
		CredDefID:     "credDefID",
		Attributes:    []*model.CredentialValue{{Name: "email", Value: "emailValue"}},
		InitiatedByUs: false,
	}
	listener.AddCredential(
		conn.ID,
		credential.ID,
		credential.Role,
		credential.SchemaID,
		credential.CredDefID,
		credential.Attributes,
		credential.InitiatedByUs)
	listener.UpdateCredential(connID, credential.ID, &currentTime, &currentTime, nil)

	cred2ID := uuid.New().String()
	cred2 := &model.Credential{
		ID:            cred2ID,
		Role:          model.CredentialRoleHolder,
		SchemaID:      "schemaID",
		CredDefID:     "credDefID",
		Attributes:    []*model.CredentialValue{{Name: "email", Value: "emailValue"}},
		InitiatedByUs: false,
	}
	listener.AddCredential(
		conn2.ID,
		cred2ID,
		cred2.Role,
		cred2.SchemaID,
		cred2.CredDefID,
		cred2.Attributes,
		cred2.InitiatedByUs)
	listener.UpdateCredential(connID2, cred2ID, &currentTime, &currentTime, nil)

	// add incomplete credential
	listener.AddCredential(
		connID,
		uuid.New().String(),
		credential.Role,
		credential.SchemaID,
		credential.CredDefID,
		credential.Attributes,
		credential.InitiatedByUs)

	t.Run("get credentials", func(t *testing.T) {
		var (
			first = 2
			last  = 1
		)
		tests := []struct {
			name   string
			args   PaginationParams
			result *CredentialTestRes
			conn   *model.Pairwise
			count  int
		}{
			{"first credential", PaginationParams{first: &first}, &CredentialTestRes{Credential: credential}, conn, first},
			{"last credential", PaginationParams{last: &last}, &CredentialTestRes{Credential: cred2}, conn2, last},
		}

		for _, testCase := range tests {
			tc := testCase
			t.Run(tc.name, func(t *testing.T) {
				got, err := r.Query().Credentials(context.TODO(), tc.args.after, tc.args.before, tc.args.first, tc.args.last)
				if tc.count != len(got.Nodes) {
					t.Errorf("Credentials count is not matching got %d expected %d", len(got.Nodes), tc.count)
				}
				checkCredential(&r, tc.conn, t, tc.name, &CredentialTestRes{got.Nodes[0], err, false}, tc.result)
			})
		}
	})
}
