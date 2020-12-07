package resolver

import (
	"github.com/findy-network/findy-agent-vault/graph/model"
	data "github.com/findy-network/findy-agent-vault/tools/data/model"
	"github.com/findy-network/findy-agent-vault/tools/tools"
	"github.com/golang/glog"
)

func (l *agencyListener) AddProof(connectionID, id string, role model.ProofRole, attributes []*model.ProofAttribute, initiatedByUs bool) {
	proof := &data.InternalProof{
		BaseObject: &data.BaseObject{
			ID:        id,
			CreatedMs: tools.CurrentTimeMs(),
		},
		Role:          role,
		Attributes:    attributes,
		InitiatedByUs: initiatedByUs,
		Result:        false,
		VerifiedMs:    nil,
		ApprovedMs:    nil,
		PairwiseID:    connectionID,
	}
	desc := proof.Description()
	status := model.JobStatusWaiting
	if !initiatedByUs {
		status = model.JobStatusPending
	}
	state.Proofs().Objects().Append(proof)

	glog.Infof("Added proof %s for connection %s", proof.ID, connectionID)
	addJobWithStatus(
		id,
		model.ProtocolTypeProof,
		&id,
		initiatedByUs,
		&connectionID,
		desc,
		status,
		model.JobResultNone)
}

func (l *agencyListener) UpdateProof(connectionID, id string, approvedMs, verifiedMs, failedMs *int64) {
	var result *bool
	if verifiedMs != nil || failedMs != nil {
		r := verifiedMs != nil && failedMs == nil
		result = &r
	}
	status := state.Proofs().UpdateProof(id, result, verifiedMs, approvedMs, failedMs)
	glog.Infof("Updated proof %s for connection %s", id, connectionID)

	// TODO: handle not found
	if status != nil {
		updateJob(id, &id, &connectionID, status.Status, status.Result, status.Description)
	}
}
