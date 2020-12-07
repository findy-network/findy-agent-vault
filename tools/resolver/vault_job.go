package resolver

import (
	"context"
	"fmt"

	"github.com/findy-network/findy-agent-vault/agency"
	data "github.com/findy-network/findy-agent-vault/tools/data/model"
	"github.com/findy-network/findy-agent-vault/tools/tools"
	"github.com/findy-network/findy-agent-vault/utils"
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

	utils.LogLow().Infof("Jobs: returning jobs between %d and %d", afterIndex, beforeIndex)

	return items.JobConnection(afterIndex, beforeIndex), nil
}

func (r *queryResolver) Job(ctx context.Context, id string) (node *model.Job, err error) {
	utils.LogMed().Info("queryResolver:Job, id: ", id)

	items := state.Jobs
	edge := items.JobForID(id)
	if edge == nil {
		err = fmt.Errorf("job for id %s was not found", id)
	} else {
		node = edge.Node
	}
	return
}

func (r *jobResolver) Output(ctx context.Context, obj *model.Job) (output *model.JobOutput, err error) {
	utils.LogMed().Info("jobResolver:Output, id: ", obj.ID)
	defer err2.Return(&err)

	output = state.OutputForJob(obj.ID)

	return
}

func (r *pairwiseResolver) Jobs(
	ctx context.Context,
	obj *model.Pairwise,
	after, before *string,
	first, last *int,
	completed *bool,
) (c *model.JobConnection, err error) {
	defer err2.Return(&err)
	pagination := &PaginationParams{
		first:  first,
		last:   last,
		after:  after,
		before: before,
	}
	logPaginationRequest("pairwiseResolver:jobs", pagination)

	items := state.Jobs
	includeCompleted := completed != nil && *completed
	items = items.Filter(func(item data.APIObject) data.APIObject {
		j := item.Job()
		if (j.Status != model.JobStatusComplete || includeCompleted) && (j.PairwiseID != nil && *j.PairwiseID == obj.ID) {
			return j.Copy()
		}
		return nil
	})

	afterIndex, beforeIndex, err := pick(items, pagination)
	err2.Check(err)

	utils.LogLow().Infof("Jobs: returning jobs between %d and %d", afterIndex, beforeIndex)

	return items.JobConnection(afterIndex, beforeIndex), nil
}

func addJob(
	id string,
	protocol model.ProtocolType,
	protocolID *string,
	initiatedByUs bool,
	pairwiseID *string,
	description string) {
	addJobWithStatus(
		id,
		protocol,
		protocolID,
		initiatedByUs,
		pairwiseID,
		description,
		model.JobStatusWaiting,
		model.JobResultNone)
}

func addJobWithStatus(
	id string,
	protocol model.ProtocolType,
	protocolID *string,
	initiatedByUs bool,
	pairwiseID *string,
	description string,
	status model.JobStatus,
	result model.JobResult) {
	timeNow := tools.CurrentTimeMs()
	items := state.Jobs
	items.Append(&data.InternalJob{
		BaseObject: &data.BaseObject{
			ID:        id,
			CreatedMs: timeNow,
		},
		ProtocolType:  protocol,
		ProtocolID:    protocolID,
		InitiatedByUs: initiatedByUs,
		PairwiseID:    pairwiseID,
		Status:        status,
		Result:        result,
		UpdatedMs:     timeNow,
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

func (r *mutationResolver) Resume(ctx context.Context, input model.ResumeJobInput) (res *model.Response, err error) {
	defer err2.Return(&err)
	utils.LogMed().Info("mutationResolver:Resume")

	job := state.Jobs.JobDataForID(input.ID)
	if job == nil {
		return nil, fmt.Errorf("job not found with id %s", input.ID)
	}

	desc := "Accepted"
	if !input.Accept {
		desc = "Declined"
	}

	switch job.ProtocolType {
	case model.ProtocolTypeCredential:
		err2.Check(agency.Instance.ResumeCredentialOffer(ctx, *job.ProtocolID, input.Accept))
		desc += " credential"
	case model.ProtocolTypeProof:
		err2.Check(agency.Instance.ResumeProofRequest(ctx, *job.ProtocolID, input.Accept))
		desc += " proof"
	case model.ProtocolTypeBasicMessage:
	case model.ProtocolTypeConnection:
	case model.ProtocolTypeNone:
		// N/A
		break
	}

	res = &model.Response{Ok: true}

	updateJob(
		input.ID,
		job.ProtocolID,
		job.PairwiseID,
		model.JobStatusWaiting,
		model.JobResultNone,
		desc)

	return res, nil
}
