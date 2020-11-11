package resolver

import (
	"context"
	"fmt"

	"github.com/findy-network/findy-agent-vault/tools/utils"
	"github.com/golang/glog"

	"github.com/lainio/err2"

	"github.com/findy-network/findy-agent-vault/graph/model"
	data "github.com/findy-network/findy-agent-vault/tools/data/model"
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

	items := state.Connections

	afterIndex, beforeIndex, err := pick(items, pagination)
	err2.Check(err)

	glog.V(logLevelLow).Infof("Connections: returning connections between %d and %d", afterIndex, beforeIndex)
	c = items.PairwiseConnection(afterIndex, beforeIndex)

	return
}

func (r *queryResolver) Connection(_ context.Context, id string) (edge *model.PairwiseEdge, err error) {
	glog.V(logLevelMedium).Info("queryResolver:Connection, id: ", id)

	items := state.Connections
	edge = items.PairwiseForID(id)
	if edge == nil {
		err = fmt.Errorf("connection for id %s was not found", id)
	}
	return
}

func doAddConnection(connection *data.InternalPairwise) {
	items := state.Connections
	connection.CreatedMs = utils.CurrentTimeMs()
	initiatedByUs := state.Jobs.IsJobInitiatedByUs(connection.ID)
	if initiatedByUs != nil {
		connection.InitiatedByUs = *initiatedByUs
	}
	items.Append(connection)
	glog.Infof("Added connection %s", connection.ID)
	updateJob(
		connection.ID,
		&connection.ID,
		&connection.ID,
		model.JobStatusComplete,
		model.JobResultSuccess,
		"Established connection to "+connection.TheirLabel)
}
