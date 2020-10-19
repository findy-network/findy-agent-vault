package resolver

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/findy-network/findy-agent-api/tools/data"

	"github.com/findy-network/findy-agent-api/graph/model"
	"github.com/findy-network/findy-agent-api/resolver"
)

type EvTestRes struct {
	Events *model.EventConnection
	Error  error
}

func TestGetEvents(t *testing.T) {
	t.Run("get events", func(t *testing.T) {
		state := data.State.Events
		var (
			valid              = 1
			tooLow             = 0
			tooHigh            = 101
			invalidCursor      = "1"
			missingError       = errors.New(resolver.ErrorFirstLastMissing)
			invalidCountError  = errors.New(resolver.ErrorFirstLastInvalid)
			invalidCursorError = errors.New(resolver.ErrorCursorInvalid)
			first              = state.EventConnection(0, 1)
			second             = state.EventConnection(1, 2)
			last               = state.EventConnection(state.Count()-1, state.Count())
		)
		tests := []struct {
			name   string
			args   PaginationParams
			result EvTestRes
		}{
			{"events, pagination missing", PaginationParams{}, EvTestRes{Error: missingError}},
			{"events, pagination first too low", PaginationParams{first: &tooLow}, EvTestRes{Error: invalidCountError}},
			{"events, pagination first too high", PaginationParams{first: &tooHigh}, EvTestRes{Error: invalidCountError}},
			{"events, pagination last too low", PaginationParams{last: &tooLow}, EvTestRes{Error: invalidCountError}},
			{"events, pagination last too high", PaginationParams{last: &tooHigh}, EvTestRes{Error: invalidCountError}},
			{"events, after cursor invalid", PaginationParams{first: &valid, after: &invalidCursor}, EvTestRes{Error: invalidCursorError}},
			{"first event", PaginationParams{first: &valid}, EvTestRes{Events: first}},
			{"last event", PaginationParams{last: &valid}, EvTestRes{Events: last}},
			{"second event", PaginationParams{first: &valid, after: &first.Edges[0].Cursor}, EvTestRes{Events: second}},
			{"previous to second event", PaginationParams{first: &valid, before: &second.Edges[0].Cursor}, EvTestRes{Events: first}},
		}

		r := Resolver{}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := r.Query().Events(context.TODO(), tt.args.after, tt.args.before, tt.args.first, tt.args.last)
				result := EvTestRes{
					Events: got,
					Error:  err,
				}
				if !reflect.DeepEqual(result, tt.result) {
					t.Errorf("%s = %v, want %v", tt.name, result, tt.result)
				}
			})
		}

	})
}
