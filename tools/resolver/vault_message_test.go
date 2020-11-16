package resolver

import (
	"context"
	"reflect"
	"testing"

	"github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/google/uuid"
)

func TestGetMessage(t *testing.T) {
	t.Run("get message", func(t *testing.T) {
		// add new message
		listener := &agencyListener{}
		id := uuid.New().String()
		msg := "Hello world"
		pw := firstPairwise()
		listener.AddMessage(pw.ID, id, msg, true)

		type TestRes struct {
			Message *model.BasicMessage
			Error   bool
		}
		tests := []struct {
			name   string
			ID     string
			result TestRes
		}{
			{"added message", id, TestRes{Message: &model.BasicMessage{ID: id, Message: msg, SentByMe: true}}},
			{"invalid connection", "", TestRes{Error: true}},
		}

		r := Resolver{}
		for _, testCase := range tests {
			tc := testCase
			t.Run(tc.name, func(t *testing.T) {
				got, err := r.Query().Message(context.TODO(), tc.ID)
				if err == nil && !tc.result.Error {
					if tc.result.Message.ID != got.ID ||
						tc.result.Message.Message != got.Message ||
						tc.result.Message.SentByMe != got.SentByMe ||
						got.CreatedMs == "" {
						t.Errorf("%s = %v, want %v", tc.name, got, tc.result.Message)
					}
					gotConn, err := r.BasicMessage().Connection(context.TODO(), got)
					if err != nil || !reflect.DeepEqual(gotConn, pw) {
						t.Errorf("%s = %v, want %v", "get message connection", got, pw)
					}

				} else if err != nil && !tc.result.Error {
					t.Errorf("%s failed, expecting for error", tc.name)
				}
			})
		}
	})
}
