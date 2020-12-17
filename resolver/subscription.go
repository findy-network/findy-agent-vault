package resolver

import (
	"context"
	"strconv"

	"github.com/findy-network/findy-agent-vault/db/store"
	"github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/lainio/err2"
)

func (r *subscriptionResolver) eventAdded(ctx context.Context) (ch <-chan *model.EventEdge, err error) {
	defer err2.Return(&err)

	agent, err := store.GetAgent(ctx, r.db)
	err2.Check(err)

	id := agent.ID + "-" + strconv.FormatInt(utils.CurrentTimeMs(), 10)
	utils.LogMed().Info("subscriptionResolver:EventAdded, id: ", id)

	events := make(chan *model.EventEdge, 1)

	go func() {
		<-ctx.Done()
		utils.LogMed().Info("subscriptionResolver: event observer removed, id: ", id)
		delete(r.eventObservers, id)
	}()

	r.eventObservers[id] = events

	return events, err
}
