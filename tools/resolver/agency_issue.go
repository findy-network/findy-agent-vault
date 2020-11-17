package resolver

import (
	"github.com/findy-network/findy-agent-vault/graph/model"
	data "github.com/findy-network/findy-agent-vault/tools/data/model"
	"github.com/findy-network/findy-agent-vault/tools/utils"
	"github.com/golang/glog"
)

func credentialStatus(issuedMs, approvedMs *int64, role model.CredentialRole) string {
	if issuedMs != nil {
		switch role {
		case model.CredentialRoleIssuer:
			return "Issued credential"
		case model.CredentialRoleHolder:
			return "Received credential"
		}
	} else if approvedMs != nil {
		return "Approved credential"
	}
	switch role {
	case model.CredentialRoleIssuer:
		return "Received credential request"
	case model.CredentialRoleHolder:
		return "Received credential offer"
	}
	return ""
}

func credentialStatusForCred(cred *data.InternalCredential) string {
	return credentialStatus(cred.IssuedMs, cred.ApprovedMs, cred.Role)
}

func (l *agencyListener) AddCredential(connectionID, id string, role model.CredentialRole, schemaID, credDefID string, attributes []*model.CredentialValue, initiatedByUs bool) {
	cred := &data.InternalCredential{
		BaseObject: &data.BaseObject{
			ID:        id,
			CreatedMs: utils.CurrentTimeMs(),
		},
		Role:          role,
		SchemaID:      schemaID,
		CredDefID:     credDefID,
		Attributes:    attributes,
		InitiatedByUs: initiatedByUs,
		ApprovedMs:    nil,
		IssuedMs:      nil,
		PairwiseID:    connectionID,
	}
	state.Connections().Objects().Append(cred)

	glog.Infof("Added credential %s", cred.ID)
	addJob(
		id,
		model.ProtocolTypeCredential,
		&id,
		initiatedByUs,
		&connectionID,
		credentialStatusForCred(cred))
}

func (l *agencyListener) UpdateCredential(connectionID, id string, approvedMs, issuedMs *int64, failed bool) {
	role := state.Credentials().UpdateCredential(id, approvedMs, issuedMs)
	glog.Infof("Updated credential %s", id)

	status := model.JobStatusWaiting
	result := model.JobResultNone
	if failed {
		status = model.JobStatusComplete
		result = model.JobResultFailure
	} else if approvedMs == nil && issuedMs == nil {
		status = model.JobStatusPending
	} else if issuedMs != nil {
		status = model.JobStatusComplete
		result = model.JobResultSuccess
	}

	// TODO: handle not found error properly
	desc := "ERROR"
	if role != nil {
		desc = credentialStatus(issuedMs, approvedMs, *role)
	}

	updateJob(id, &id, &connectionID, status, result, desc)
}
