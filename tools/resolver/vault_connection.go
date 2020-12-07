package resolver

import (
	"context"
	"fmt"

	"github.com/lainio/err2"

	"github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/utils"
)

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

	items := state.Connections()

	afterIndex, beforeIndex, err := pick(items.Objects(), pagination)
	err2.Check(err)

	utils.LogLow().Infof("Connections: returning connections between %d and %d", afterIndex, beforeIndex)
	c = items.PairwiseConnection(afterIndex, beforeIndex)

	return
}

func (r *queryResolver) Connection(_ context.Context, id string) (node *model.Pairwise, err error) {
	utils.LogMed().Info("queryResolver:Connection, id: ", id)

	items := state.Connections()
	edge := items.PairwiseForID(id)
	if edge == nil {
		err = fmt.Errorf("connection for id %s was not found", id)
	} else {
		node = edge.Node
	}
	return
}
