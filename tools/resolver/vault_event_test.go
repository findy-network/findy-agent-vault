package resolver

import (
	"context"
	"reflect"
	"testing"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

type EvTestRes struct {
	Events *model.EventConnection
	Error  error
}

type EventsExecutor struct{}

func (*EventsExecutor) Request(ctx context.Context, after, before *string, first, last *int) error {
	r := Resolver{}
	_, err := r.Query().Events(context.TODO(), after, before, first, last)
	return err
}

func TestPaginationErrorsGetEvents(t *testing.T) {
	testPaginationErrors(t, "events", &EventsExecutor{})
}

func TestGetEvents(t *testing.T) {
	t.Run("get events", func(t *testing.T) {
		s := state.Events
		var (
			valid  = 1
			first  = s.EventConnection(0, 1, state.Connections)
			second = s.EventConnection(1, 2, state.Connections)
			last   = s.EventConnection(s.Count()-1, s.Count(), state.Connections)
		)
		tests := []struct {
			name   string
			args   PaginationParams
			result EvTestRes
		}{
			{"first event", PaginationParams{first: &valid}, EvTestRes{Events: first}},
			{"last event", PaginationParams{last: &valid}, EvTestRes{Events: last}},
			{"second event", PaginationParams{first: &valid, after: &first.Edges[0].Cursor}, EvTestRes{Events: second}},
			{"previous to second event", PaginationParams{first: &valid, before: &second.Edges[0].Cursor}, EvTestRes{Events: first}},
		}

		r := Resolver{}
		for _, testCase := range tests {
			tc := testCase
			t.Run(tc.name, func(t *testing.T) {
				got, err := r.Query().Events(context.TODO(), tc.args.after, tc.args.before, tc.args.first, tc.args.last)
				result := EvTestRes{
					Events: got,
					Error:  err,
				}
				if !reflect.DeepEqual(result, tc.result) {
					t.Errorf("%s = %v, want %v", tc.name, result, tc.result)
				}
			})
		}
	})
}
