package resolver

import (
	"github.com/findy-network/findy-agent-vault/graph/model"
	data "github.com/findy-network/findy-agent-vault/tools/data/model"
	"github.com/findy-network/findy-agent-vault/tools/utils"
	"github.com/golang/glog"
)

func (l *agencyListener) AddCredential(
	connectionID, id string,
	role model.CredentialRole,
	schemaID, credDefID string,
	attributes []*model.CredentialValue,
	initiatedByUs bool,
) {
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
	desc := cred.Description()
	state.Connections().Objects().Append(cred)

	glog.Infof("Added credential %s", id)
	addJob(
		id,
		model.ProtocolTypeCredential,
		&id,
		initiatedByUs,
		&connectionID,
		desc,
	)
}

func (l *agencyListener) UpdateCredential(connectionID, id string, approvedMs, issuedMs, failedMs *int64) {
	status := state.Credentials().UpdateCredential(id, approvedMs, issuedMs, failedMs)
	glog.Infof("Updated credential %s", id)

	if status != nil {
		updateJob(id, &id, &connectionID, status.Status, status.Result, status.Description)
	}
}
