package resolver

import (
	"context"
	"fmt"

	"github.com/golang/glog"

	"github.com/lainio/err2"

	"github.com/findy-network/findy-agent-api/graph/model"
	"github.com/findy-network/findy-agent-api/tools/data"
)

/*

if first, last missing, return error

Start from the greedy query: SELECT * FROM table ORDER BY created
If the after argument is provided, add id > parsed_cursor to the WHERE clause
If the before argument is provided, add id < parsed_cursor to the WHERE clause
If the first argument is provided, add ORDER BY id DESC LIMIT first+1 to the query
If the last argument is provided, add ORDER BY id ASC LIMIT last+1 to the query
If the last argument is provided, I reverse the order of the results
If the first argument is provided then I set hasPreviousPage: false (see spec for a description of this behavior).
If no less than first+1 results are returned, I set hasNextPage: true, otherwise I set it to false.
If the last argument is provided then I set hasNextPage: false (see spec for a description of this behavior).
If no less last+1 results are returned, I set hasPreviousPage: true, otherwise I set it to false.
*/

func (r *queryResolver) Connections(
	_ context.Context,
	after *string, before *string,
	first *int, last *int) (c *model.PairwiseConnection, err error) {
	defer err2.Return(&err)

	pagination := &PaginationParams{
		first:  first,
		last:   last,
		after:  after,
		before: before,
	}
	logPaginationRequest("queryResolver:conns", pagination)

	state := data.State.Connections
	afterIndex, beforeIndex, err := pick(state, pagination)
	err2.Check(err)

	return state.PairwiseConnection(afterIndex, beforeIndex), nil
}

func (r *queryResolver) Connection(_ context.Context, id string) (node *model.Pairwise, err error) {
	glog.V(logLevelMedium).Info("queryResolver:Connection, id: ", id)

	state := data.State.Connections
	node = state.PairwiseForID(id)
	if node == nil {
		err = fmt.Errorf("connection for id %s was not found", id)
	}
	return
}
