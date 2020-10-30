package resolver

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/findy-network/findy-agent-vault/resolver"
)

type PaginationExecutor interface {
	Request(ctx context.Context, after *string, before *string, first *int, last *int) error
}

func testPaginationErrors(t *testing.T, objName string, ex PaginationExecutor) {
	t.Run(fmt.Sprintf("get %s", objName), func(t *testing.T) {
		var (
			valid              = 1
			tooLow             = 0
			tooHigh            = 101
			invalidCursor      = "1"
			missingError       = errors.New(resolver.ErrorFirstLastMissing)
			invalidCountError  = errors.New(resolver.ErrorFirstLastInvalid)
			invalidCursorError = errors.New(resolver.ErrorCursorInvalid)
		)
		tests := []struct {
			name string
			args PaginationParams
			err  error
		}{
			{fmt.Sprintf("%s, pagination missing", objName), PaginationParams{}, missingError},
			{fmt.Sprintf("%s, pagination first too low", objName), PaginationParams{first: &tooLow}, invalidCountError},
			{fmt.Sprintf("%s, pagination first too high", objName), PaginationParams{first: &tooHigh}, invalidCountError},
			{fmt.Sprintf("%s, pagination last too low", objName), PaginationParams{last: &tooLow}, invalidCountError},
			{fmt.Sprintf("%s, pagination last too high", objName), PaginationParams{last: &tooHigh}, invalidCountError},
			{fmt.Sprintf("%s, after cursor invalid", objName), PaginationParams{first: &valid, after: &invalidCursor}, invalidCursorError},
		}

		for _, testCase := range tests {
			tc := testCase
			t.Run(tc.name, func(t *testing.T) {
				err := ex.Request(context.TODO(), tc.args.after, tc.args.before, tc.args.first, tc.args.last)
				if !reflect.DeepEqual(err, tc.err) {
					t.Errorf("%s = err (%v)\n want (%v)", tc.name, err, tc.err)
				}
			})
		}
	})
}
