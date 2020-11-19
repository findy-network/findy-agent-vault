package resolver

import (
	"context"
	"reflect"
	"testing"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

type JobTestRes struct {
	Jobs  *model.JobConnection
	Error error
}

func TestPaginationErrorsGetJobs(t *testing.T) {
	testPaginationErrors(t, "jobs", func(ctx context.Context, after, before *string, first, last *int) error {
		r := Resolver{}
		completed := false
		_, err := r.Query().Jobs(context.TODO(), after, before, first, last, &completed)
		return err
	})
}

func TestGetIncompleteJobs(t *testing.T) {
	resetResolver(false)
	t.Run("get incomplete jobs", func(t *testing.T) {
		s := state.Jobs
		var (
			valid  = 1
			first  = s.JobConnection(0, 1)
			second = s.JobConnection(1, 2)
			last   = s.JobConnection(s.Count()-1, s.Count())
		)
		tests := []struct {
			name   string
			args   PaginationParams
			result JobTestRes
		}{
			{"first job", PaginationParams{first: &valid}, JobTestRes{Jobs: first}},
			{"last job", PaginationParams{last: &valid}, JobTestRes{Jobs: last}},
			{"second job", PaginationParams{first: &valid, after: &first.Edges[0].Cursor}, JobTestRes{Jobs: second}},
			{"previous to second job", PaginationParams{first: &valid, before: &second.Edges[0].Cursor}, JobTestRes{Jobs: first}},
		}

		r := Resolver{}
		for _, testCase := range tests {
			tc := testCase
			t.Run(tc.name, func(t *testing.T) {
				completed := false
				got, err := r.Query().Jobs(context.TODO(), tc.args.after, tc.args.before, tc.args.first, tc.args.last, &completed)
				result := JobTestRes{
					Jobs:  got,
					Error: err,
				}
				if !reflect.DeepEqual(result, tc.result) {
					t.Errorf("%s = %v, want %v", tc.name, result, tc.result)
				}
			})
		}
	})
}
