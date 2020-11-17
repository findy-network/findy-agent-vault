package resolver

import (
	"context"
	"reflect"
	"testing"

	"github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/google/uuid"
)

type MsgTestRes struct {
	Message      *model.BasicMessage
	Error        error
	ExpectsError bool
}

func checkMessage(r *Resolver, pw *model.Pairwise, t *testing.T, name string, got, expected *MsgTestRes) {
	if got.Error == nil && !expected.ExpectsError {
		if expected.Message.ID != got.Message.ID ||
			expected.Message.Message != got.Message.Message ||
			expected.Message.SentByMe != got.Message.SentByMe ||
			got.Message.CreatedMs == "" {
			t.Errorf("%s = %v, want %v", name, got.Message, expected.Message)
		}
		gotConn, err := r.BasicMessage().Connection(context.TODO(), got.Message)
		if err != nil || !reflect.DeepEqual(gotConn, pw) {
			t.Errorf("%s = %v, want %v", "get message connection", got, pw)
		}

	} else if got.Error != nil && !expected.ExpectsError {
		t.Errorf("%s failed, expecting for error", name)
	}

}

func TestGetMessage(t *testing.T) {
	t.Run("get message", func(t *testing.T) {
		// add new message
		listener := &agencyListener{}
		id := uuid.New().String()
		msg := "Hello world"
		pw := firstPairwise()
		listener.AddMessage(pw.ID, id, msg, true)

		tests := []struct {
			name   string
			ID     string
			result *MsgTestRes
		}{
			{"added message", id, &MsgTestRes{Message: &model.BasicMessage{ID: id, Message: msg, SentByMe: true}}},
			{"invalid connection", "", &MsgTestRes{ExpectsError: true}},
		}

		r := Resolver{}
		for _, testCase := range tests {
			tc := testCase
			t.Run(tc.name, func(t *testing.T) {
				got, err := r.Query().Message(context.TODO(), tc.ID)
				checkMessage(&r, pw, t, tc.name, &MsgTestRes{got, err, false}, tc.result)
			})
		}
	})
}

func TestPaginationErrorsGetMessages(t *testing.T) {
	testPaginationErrors(t, "messages", func(ctx context.Context, after, before *string, first, last *int) error {
		r := Resolver{}
		pw := firstPairwise()
		_, err := r.Pairwise().Messages(context.TODO(), pw, after, before, first, last)
		return err
	})
}

func TestGetMessages(t *testing.T) {
	r := Resolver{}

	// add new connection and messages
	listener := &agencyListener{}

	connId := uuid.New().String()
	listener.AddConnection(connId, "ourDID", "theirDID", "theirEndpoint", "theirLabel")
	conn, _ := r.Query().Connection(context.TODO(), connId)

	msgID1 := uuid.New().String()
	listener.AddMessage(connId, msgID1, msgID1, true)

	msgID2 := uuid.New().String()
	listener.AddMessage(connId, msgID2, msgID2, true)

	msgID3 := uuid.New().String()
	listener.AddMessage(connId, msgID3, msgID3, true)

	t.Run("get messages", func(t *testing.T) {
		var (
			valid = 1
		)
		tests := []struct {
			name   string
			args   PaginationParams
			result *MsgTestRes
		}{
			{"first message", PaginationParams{first: &valid}, &MsgTestRes{Message: &model.BasicMessage{ID: msgID1, Message: msgID1, SentByMe: true}}},
			{"last message", PaginationParams{last: &valid}, &MsgTestRes{Message: &model.BasicMessage{ID: msgID3, Message: msgID3, SentByMe: true}}},
		}

		for _, testCase := range tests {
			tc := testCase
			t.Run(tc.name, func(t *testing.T) {
				got, err := r.Pairwise().Messages(context.TODO(), conn, tc.args.after, tc.args.before, tc.args.first, tc.args.last)
				checkMessage(&r, conn, t, tc.name, &MsgTestRes{got.Nodes[0], err, false}, tc.result)
			})
		}
	})
}
