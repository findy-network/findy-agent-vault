package resolver

import (
	"context"
	"fmt"

	data "github.com/findy-network/findy-agent-vault/tools/data/model"
	"github.com/findy-network/findy-agent-vault/tools/faker"
	"github.com/findy-network/findy-agent-vault/utils"

	"github.com/lainio/err2"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

func (r *basicMessageResolver) Connection(ctx context.Context, obj *model.BasicMessage) (pw *model.Pairwise, err error) {
	utils.LogMed().Info("basicMessageResolver:Connection, id: ", obj.ID)
	defer err2.Return(&err)

	if connectionID := state.Messages.MessagePairwiseID(obj.ID); connectionID != nil {
		return r.Query().Connection(ctx, *connectionID)
	}

	err = fmt.Errorf("pairwise for message id %s was not found", obj.ID)
	return
}

func (r *pairwiseResolver) Messages(
	ctx context.Context,
	pw *model.Pairwise,
	after, before *string,
	first, last *int) (c *model.BasicMessageConnection, err error) {
	defer err2.Return(&err)
	pagination := &PaginationParams{
		first:  first,
		last:   last,
		after:  after,
		before: before,
	}
	logPaginationRequest("queryResolver:messages", pagination)

	items := state.Messages
	items = items.Filter(func(item data.APIObject) data.APIObject {
		m := item.BasicMessage()
		if m.PairwiseID == pw.ID {
			return m.Copy()
		}
		return nil
	})

	afterIndex, beforeIndex, err := pick(items, pagination)
	err2.Check(err)

	utils.LogLow().Infof("Messages: returning messages between %d and %d", afterIndex, beforeIndex)

	return items.MessageConnection(afterIndex, beforeIndex), nil
}

func (r *queryResolver) Message(ctx context.Context, id string) (node *model.BasicMessage, err error) {
	utils.LogMed().Info("queryResolver:Message, id: ", id)

	items := state.Messages
	edge := items.MessageForID(id)
	if edge == nil {
		err = fmt.Errorf("connection for id %s was not found", id)
	} else {
		node = edge.Node
	}
	return
}

func (r *mutationResolver) AddRandomMessage(ctx context.Context) (ok bool, err error) {
	utils.LogMed().Info("mutationResolver:AddRandomMessage ")
	defer err2.Return(&err)

	msgs, err := faker.FakeMessages(1)
	err2.Check(err)

	msg := msgs[0]
	r.listener.AddMessage(msg.PairwiseID, msg.ID, msg.Message, msg.SentByMe)

	ok = true

	return
}
