package resolver

import (
	"context"
	"reflect"
	"testing"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

type ConnectionsExecutor struct{}

func (*ConnectionsExecutor) Request(ctx context.Context, after, before *string, first, last *int) error {
	r := Resolver{}
	_, err := r.Query().Connections(context.TODO(), after, before, first, last)
	return err
}

func TestPaginationErrorsGetConnections(t *testing.T) {
	testPaginationErrors(t, "connections", &ConnectionsExecutor{})
}

func TestGetConnections(t *testing.T) {
	t.Run("get connections", func(t *testing.T) {
		s := state.Connections()
		var (
			valid  = 1
			first  = s.PairwiseConnection(0, 1)
			second = s.PairwiseConnection(1, 2)
			last   = s.PairwiseConnection(s.Objects().Count()-1, s.Objects().Count())
		)
		type ConnTestRes struct {
			Connections *model.PairwiseConnection
			Error       error
		}
		tests := []struct {
			name   string
			args   PaginationParams
			result ConnTestRes
		}{
			{"first connection", PaginationParams{first: &valid}, ConnTestRes{Connections: first}},
			{"last connection", PaginationParams{last: &valid}, ConnTestRes{Connections: last}},
			{"second connection", PaginationParams{first: &valid, after: &first.Edges[0].Cursor}, ConnTestRes{Connections: second}},
			{"previous to second connection", PaginationParams{first: &valid, before: &second.Edges[0].Cursor}, ConnTestRes{Connections: first}},
		}

		r := Resolver{}
		for _, testCase := range tests {
			tc := testCase
			t.Run(tc.name, func(t *testing.T) {
				got, err := r.Query().Connections(
					context.TODO(),
					tc.args.after, tc.args.before, tc.args.first, tc.args.last)
				result := ConnTestRes{
					Connections: got,
					Error:       err,
				}
				if !reflect.DeepEqual(result, tc.result) {
					t.Errorf("%s = %v, want %v", tc.name, result, tc.result)
				}
			})
		}
	})
}

func TestGetConnection(t *testing.T) {
	t.Run("get connection", func(t *testing.T) {
		var (
			first = firstPairwise()
		)
		type ConnTestRes struct {
			Connection *model.Pairwise
			Error      bool
		}
		tests := []struct {
			name   string
			ID     string
			result ConnTestRes
		}{
			{"first connection", first.ID, ConnTestRes{Connection: first}},
			{"invalid connection", "", ConnTestRes{Error: true}},
		}

		r := Resolver{}
		for _, testCase := range tests {
			tc := testCase
			t.Run(tc.name, func(t *testing.T) {
				got, err := r.Query().Connection(context.TODO(), tc.ID)
				if err == nil && !tc.result.Error {
					if !reflect.DeepEqual(ConnTestRes{Connection: got}, tc.result) {
						t.Errorf("%s = %v, want %v", tc.name, got, tc.result.Connection)
					}
				} else if err != nil && !tc.result.Error {
					t.Errorf("%s failed, expecting for error", tc.name)
				}
			})
		}
	})
}
