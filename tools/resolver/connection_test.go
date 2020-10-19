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

type ConnTestRes struct {
	Connections *model.PairwiseConnection
	Error       error
}

func TestGetConnections(t *testing.T) {
	t.Run("get connections", func(t *testing.T) {
		state := data.State.Connections
		var (
			valid              = 1
			tooLow             = 0
			tooHigh            = 101
			invalidCursor      = "1"
			missingError       = errors.New(resolver.ErrorFirstLastMissing)
			invalidCountError  = errors.New(resolver.ErrorFirstLastInvalid)
			invalidCursorError = errors.New(resolver.ErrorCursorInvalid)
			first              = state.PairwiseConnection(0, 1)
			second             = state.PairwiseConnection(1, 2)
			last               = state.PairwiseConnection(state.Count()-1, state.Count())
		)
		tests := []struct {
			name   string
			args   PaginationParams
			result ConnTestRes
		}{
			{"connections, pagination missing", PaginationParams{}, ConnTestRes{Error: missingError}},
			{"connections, pagination first too low", PaginationParams{first: &tooLow}, ConnTestRes{Error: invalidCountError}},
			{"connections, pagination first too high", PaginationParams{first: &tooHigh}, ConnTestRes{Error: invalidCountError}},
			{"connections, pagination last too low", PaginationParams{last: &tooLow}, ConnTestRes{Error: invalidCountError}},
			{"connections, pagination last too high", PaginationParams{last: &tooHigh}, ConnTestRes{Error: invalidCountError}},
			{"connections, after cursor invalid", PaginationParams{first: &valid, after: &invalidCursor}, ConnTestRes{Error: invalidCursorError}},
			{"first connection", PaginationParams{first: &valid}, ConnTestRes{Connections: first}},
			{"last connection", PaginationParams{last: &valid}, ConnTestRes{Connections: last}},
			{"second connection", PaginationParams{first: &valid, after: &first.Edges[0].Cursor}, ConnTestRes{Connections: second}},
			{"previous to second connection", PaginationParams{first: &valid, before: &second.Edges[0].Cursor}, ConnTestRes{Connections: first}},
		}

		r := Resolver{}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := r.Query().Connections(context.TODO(), tt.args.after, tt.args.before, tt.args.first, tt.args.last)
				result := ConnTestRes{
					Connections: got,
					Error:       err,
				}
				if !reflect.DeepEqual(result, tt.result) {
					t.Errorf("%s = %v, want %v", tt.name, result, tt.result)
				}
			})
		}

	})
}
