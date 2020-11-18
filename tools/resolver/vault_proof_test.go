package resolver

import (
	"context"
	"reflect"
	"testing"

	"github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/tools/utils"
	"github.com/google/uuid"
)

type ProofTestRes struct {
	Proof        *model.Proof
	Error        error
	ExpectsError bool
}

func checkProof(r *Resolver, pw *model.Pairwise, t *testing.T, name string, got, expected *ProofTestRes) {
	if got.Error == nil && !expected.ExpectsError {
		if expected.Proof.ID != got.Proof.ID ||
			expected.Proof.Role != got.Proof.Role ||
			expected.Proof.InitiatedByUs != got.Proof.InitiatedByUs ||
			!reflect.DeepEqual(expected.Proof.Attributes, got.Proof.Attributes) ||
			got.Proof.CreatedMs == "" {
			t.Errorf("%s = %v, want %v", name, got.Proof, expected.Proof)
		}
		gotConn, err := r.Proof().Connection(context.TODO(), got.Proof)
		if err != nil || !reflect.DeepEqual(gotConn, pw) {
			t.Errorf("%s = %v, want %v", "get proof connection", got, pw)
		}
	} else if got.Error != nil && !expected.ExpectsError {
		t.Errorf("%s failed, expecting for error", name)
	}
}

func TestGetProof(t *testing.T) {
	t.Run("get proof", func(t *testing.T) {
		// add new proof
		listener := &agencyListener{}
		id := uuid.New().String()
		pw := firstPairwise()
		proof := &model.Proof{
			ID:            id,
			Role:          model.ProofRoleProver,
			Attributes:    []*model.ProofAttribute{{Name: "email", CredDefID: "credDefID"}},
			InitiatedByUs: false,
		}
		listener.AddProof(pw.ID, proof.ID, proof.Role, proof.Attributes, proof.InitiatedByUs)

		tests := []struct {
			name   string
			ID     string
			result *ProofTestRes
		}{
			{"added message", id, &ProofTestRes{Proof: proof}},
			{"invalid connection", "", &ProofTestRes{ExpectsError: true}},
		}

		r := Resolver{}
		for _, testCase := range tests {
			tc := testCase
			t.Run(tc.name, func(t *testing.T) {
				got, err := r.Query().Proof(context.TODO(), tc.ID)
				checkProof(&r, pw, t, tc.name, &ProofTestRes{got, err, false}, tc.result)
			})
		}
	})
}

func TestPaginationErrorsGetProofs(t *testing.T) {
	testPaginationErrors(t, "proofs", func(ctx context.Context, after, before *string, first, last *int) error {
		r := Resolver{}
		pw := firstPairwise()
		_, err := r.Pairwise().Proofs(context.TODO(), pw, after, before, first, last)
		return err
	})
}

func TestGetProofs(t *testing.T) {
	r := Resolver{}

	// add new connection
	listener := &agencyListener{}
	currentTime := utils.CurrentTimeMs()

	connID := uuid.New().String()
	listener.AddConnection(connID, "ourDID", "theirDID", "theirEndpoint", "theirLabel")
	conn, _ := r.Query().Connection(context.TODO(), connID)

	id := uuid.New().String()
	proof := &model.Proof{
		ID:            id,
		Role:          model.ProofRoleProver,
		Attributes:    []*model.ProofAttribute{{Name: "email", CredDefID: "credDefID"}},
		InitiatedByUs: false,
	}
	listener.AddProof(connID, proof.ID, proof.Role, proof.Attributes, proof.InitiatedByUs)
	listener.UpdateProof(connID, proof.ID, &currentTime, &currentTime, nil)

	// add incomplete proof
	listener.AddProof(connID, uuid.New().String(), proof.Role, proof.Attributes, proof.InitiatedByUs)

	t.Run("get proofs", func(t *testing.T) {
		var (
			valid = 1
		)
		tests := []struct {
			name   string
			args   PaginationParams
			result *ProofTestRes
		}{
			{"first proof", PaginationParams{first: &valid}, &ProofTestRes{Proof: proof}},
			{"last proof", PaginationParams{last: &valid}, &ProofTestRes{Proof: proof}},
		}

		for _, testCase := range tests {
			tc := testCase
			t.Run(tc.name, func(t *testing.T) {
				got, err := r.Pairwise().Proofs(context.TODO(), conn, tc.args.after, tc.args.before, tc.args.first, tc.args.last)
				checkProof(&r, conn, t, tc.name, &ProofTestRes{got.Nodes[0], err, false}, tc.result)
			})
		}
	})
}
