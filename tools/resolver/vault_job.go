package resolver

import (
	"context"
	"fmt"
	"time"

	data "github.com/findy-network/findy-agent-vault/tools/data/model"
	"github.com/golang/glog"
	"github.com/lainio/err2"

	"github.com/findy-network/findy-agent-vault/graph/model"
)

func (r *queryResolver) Jobs(
	_ context.Context,
	after, before *string,
	first, last *int,
	completed *bool) (c *model.JobConnection, err error) {
	defer err2.Return(&err)
	pagination := &PaginationParams{
		first:  first,
		last:   last,
		after:  after,
		before: before,
	}
	logPaginationRequest("queryResolver:jobs", pagination)

	items := state.Jobs
	if completed == nil || !*completed {
		items = items.Filter(func(item data.APIObject) data.APIObject {
			j := item.Job()
			if j.Status != model.JobStatusComplete {
				return j.Copy()
			}
			return nil
		})
	}
	afterIndex, beforeIndex, err := pick(items, pagination)
	err2.Check(err)

	glog.V(logLevelLow).Infof("Jobs: returning jobs between %d and %d", afterIndex, beforeIndex)

	return items.JobConnection(afterIndex, beforeIndex, state.Connections), nil
}

func (r *queryResolver) Job(ctx context.Context, id string) (node *model.Job, err error) {
	glog.V(logLevelMedium).Info("queryResolver:Job, id: ", id)

	items := state.Jobs
	node = items.JobForID(id, state.Connections)
	if node == nil {
		err = fmt.Errorf("job for id %s was not found", id)
	}
	return
}

func addJob(id string, protocol model.ProtocolType, protocolID, pairwiseID *string, description string) {
	timeNow := time.Now().Unix()
	items := state.Jobs
	items.Append(&data.InternalJob{
		ID:           id,
		ProtocolType: protocol,
		ProtocolID:   protocolID,
		PairwiseID:   pairwiseID,
		Status:       model.JobStatusWaiting,
		Result:       model.JobResultNone,
		CreatedMs:    timeNow,
		UpdatedMs:    timeNow,
	})
	glog.Infof("Added job %s", id)
	addEvent(description, pairwiseID, &id)
}

func updateJob(id string, protocolID, pairwiseID *string, status model.JobStatus, result model.JobResult, description string) {
	items := state.Jobs
	items.UpdateJob(id, protocolID, pairwiseID, status, result)
	glog.Infof("Updated job %s", id)
	addEvent(description, pairwiseID, &id)
}
