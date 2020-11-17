package resolver

import (
	"github.com/findy-network/findy-agent-vault/graph/model"
	data "github.com/findy-network/findy-agent-vault/tools/data/model"
	"github.com/findy-network/findy-agent-vault/tools/utils"
	"github.com/golang/glog"
)

func proofStatus(verifiedMs, approvedMs *int64, role model.ProofRole) string {
	if verifiedMs != nil {
		switch role {
		case model.ProofRoleVerifier:
			return "Verified credential"
		case model.ProofRoleProver:
			return "Proved credential"
		}
	} else if approvedMs != nil {
		return "Approved proof"
	}
	switch role {
	case model.ProofRoleVerifier:
		return "Received proof offer"
	case model.ProofRoleProver:
		return "Received proof request"
	}
	return ""
}

func proofStatusForProof(proof *data.InternalProof) string {
	return proofStatus(proof.VerifiedMs, proof.ApprovedMs, proof.Role)
}

func (l *agencyListener) AddProof(connectionID, id string, role model.ProofRole, attributes []*model.ProofAttribute, initiatedByUs bool) {
	proof := &data.InternalProof{
		BaseObject: &data.BaseObject{
			ID:        id,
			CreatedMs: utils.CurrentTimeMs(),
		},
		Role:          role,
		Attributes:    attributes,
		InitiatedByUs: initiatedByUs,
		Result:        false,
		VerifiedMs:    nil,
		ApprovedMs:    nil,
		PairwiseID:    connectionID,
	}
	state.Proofs().Objects().Append(proof)

	glog.Infof("Added proof %s", proof.ID)
	addJob(
		id,
		model.ProtocolTypeProof,
		&id,
		initiatedByUs,
		&connectionID,
		proofStatusForProof(proof))
}

func (l *agencyListener) UpdateProof(connectionID, id string, approvedMs, verifiedMs *int64, failed bool) {
	var result *bool
	if verifiedMs != nil {
		r := !failed
		result = &r
	}
	role := state.Proofs().UpdateProof(id, result, verifiedMs, approvedMs)
	glog.Infof("Updated proof %s", id)

	status := model.JobStatusWaiting
	jobResult := model.JobResultNone
	if failed {
		status = model.JobStatusComplete
		jobResult = model.JobResultFailure
	} else if approvedMs == nil && verifiedMs == nil {
		status = model.JobStatusPending
	} else if verifiedMs != nil {
		status = model.JobStatusComplete
		jobResult = model.JobResultSuccess
	}

	// TODO: handle not found error properly
	desc := "ERROR"
	if role != nil {
		desc = proofStatus(verifiedMs, approvedMs, *role)
	}

	updateJob(id, &id, &connectionID, status, jobResult, desc)
}
