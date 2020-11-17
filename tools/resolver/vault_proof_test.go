package resolver

import (
	"github.com/findy-network/findy-agent-vault/graph/model"
)

type ProofTestRes struct {
	Proof        *model.Proof
	Error        error
	ExpectsError bool
}

/*func checkProof(r *Resolver, pw *model.Pairwise, t *testing.T, name string, got, expected *ProofTestRes) {
	if got.Error == nil && !expected.ExpectsError {
		if expected.Proof.ID != got.Proof.ID ||
			expected.Proof.Proof != got.Proof.Proof ||
			expected.Proof.SentByMe != got.Proof.SentByMe ||
			got.Proof.CreatedMs == "" {
			t.Errorf("%s = %v, want %v", name, got.Proof, expected.Proof)
		}
		gotConn, err := r.BasicProof().Connection(context.TODO(), got.Proof)
		if err != nil || !reflect.DeepEqual(gotConn, pw) {
			t.Errorf("%s = %v, want %v", "get message connection", got, pw)
		}

	} else if got.Error != nil && !expected.ExpectsError {
		t.Errorf("%s failed, expecting for error", name)
	}

}

func TestGetProof(t *testing.T) {
	t.Run("get message", func(t *testing.T) {
		// add new message
		listener := &agencyListener{}
		id := uuid.New().String()
		msg := "Hello world"
		pw := firstPairwise()
		listener.AddProof(pw.ID, id, msg, true)

		tests := []struct {
			name   string
			ID     string
			result *MsgTestRes
		}{
			{"added message", id, &MsgTestRes{Proof: &model.BasicProof{ID: id, Proof: msg, SentByMe: true}}},
			{"invalid connection", "", &MsgTestRes{ExpectsError: true}},
		}

		r := Resolver{}
		for _, testCase := range tests {
			tc := testCase
			t.Run(tc.name, func(t *testing.T) {
				got, err := r.Query().Proof(context.TODO(), tc.ID)
				checkProof(&r, pw, t, tc.name, &MsgTestRes{got, err, false}, tc.result)
			})
		}
	})
}

func TestPaginationErrorsGetProofs(t *testing.T) {
	testPaginationErrors(t, "messages", func(ctx context.Context, after, before *string, first, last *int) error {
		r := Resolver{}
		pw := firstPairwise()
		_, err := r.Pairwise().Proofs(context.TODO(), pw, after, before, first, last)
		return err
	})
}

func TestGetProofs(t *testing.T) {
	r := Resolver{}

	// add new connection and messages
	listener := &agencyListener{}

	connId := uuid.New().String()
	listener.AddConnection(connId, "ourDID", "theirDID", "theirEndpoint", "theirLabel")
	conn, _ := r.Query().Connection(context.TODO(), connId)

	msgID1 := uuid.New().String()
	listener.AddProof(connId, msgID1, msgID1, true)

	msgID2 := uuid.New().String()
	listener.AddProof(connId, msgID2, msgID2, true)

	msgID3 := uuid.New().String()
	listener.AddProof(connId, msgID3, msgID3, true)

	t.Run("get messages", func(t *testing.T) {
		var (
			valid = 1
		)
		tests := []struct {
			name   string
			args   PaginationParams
			result *MsgTestRes
		}{
			{"first message", PaginationParams{first: &valid}, &MsgTestRes{Proof: &model.BasicProof{ID: msgID1, Proof: msgID1, SentByMe: true}}},
			{"last message", PaginationParams{last: &valid}, &MsgTestRes{Proof: &model.BasicProof{ID: msgID3, Proof: msgID3, SentByMe: true}}},
		}

		for _, testCase := range tests {
			tc := testCase
			t.Run(tc.name, func(t *testing.T) {
				got, err := r.Pairwise().Proofs(context.TODO(), conn, tc.args.after, tc.args.before, tc.args.first, tc.args.last)
				checkProof(&r, conn, t, tc.name, &MsgTestRes{got.Nodes[0], err, false}, tc.result)
			})
		}
	})
}
*/
