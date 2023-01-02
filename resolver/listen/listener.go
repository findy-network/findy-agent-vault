package listen

import (
	"fmt"
	"time"

	agency "github.com/findy-network/findy-agent-vault/agency/model"
	dbModel "github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/db/store"
	"github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/resolver/update"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/golang/glog"
	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
)

type Listener struct {
	db store.DB
	*update.Updater
}

func NewListener(db store.DB, updater *update.Updater) *Listener {
	return &Listener{db, updater}
}

func (l *Listener) AddConnection(info *agency.JobInfo, data *agency.Connection) (err error) {
	defer err2.Handle(&err)

	utils.LogMed().Infof("Add connection %s for tenant %s", info.ConnectionID, info.TenantID)

	// TODO: job is currently created with Connection ID ->
	// use job ID instead and let agency create the ids, needs API change?
	job := try.To1(l.db.GetJob(info.ConnectionID, info.TenantID))

	now := utils.CurrentTime()

	connection := try.To1(l.db.AddConnection(
		&dbModel.Connection{
			Base: dbModel.Base{
				ID:       info.ConnectionID,
				TenantID: info.TenantID,
			},
			OurDid:        data.OurDID,
			TheirDid:      data.TheirDID,
			TheirEndpoint: data.TheirEndpoint,
			TheirLabel:    data.TheirLabel,
			Approved:      now, // TODO: get approved from agency?
			Invited:       job.InitiatedByUs,
		}))

	job.ConnectionID = &connection.ID
	job.ProtocolConnectionID = &connection.ID
	job.Status = model.JobStatusComplete
	job.Result = model.JobResultSuccess

	try.To(l.UpdateJob(
		job,
		"Established connection to "+connection.TheirLabel,
	))
	return nil
}

func (l *Listener) AddMessage(info *agency.JobInfo, data *agency.Message) (err error) {
	defer err2.Handle(&err)

	msg := try.To1(l.db.AddMessage(&dbModel.Message{
		Base:         dbModel.Base{TenantID: info.TenantID},
		ConnectionID: info.ConnectionID,
		Message:      data.Message,
		SentByMe:     data.SentByMe,
	}))

	try.To(l.AddJob(&dbModel.Job{
		Base:              dbModel.Base{ID: info.JobID, TenantID: info.TenantID},
		ConnectionID:      &info.ConnectionID,
		ProtocolType:      model.ProtocolTypeBasicMessage,
		ProtocolMessageID: &msg.ID,
		InitiatedByUs:     data.SentByMe,
		Status:            model.JobStatusComplete,
		Result:            model.JobResultSuccess,
	}, msg.Description()))
	return nil
}

func (l *Listener) UpdateMessage(info *agency.JobInfo, _ *agency.MessageUpdate) (err error) {
	// TODO: linter needs coment, implement later
	return nil
}

func (l *Listener) AddCredential(info *agency.JobInfo, data *agency.Credential) (err error) {
	defer err2.Handle(&err)

	credential := try.To1(l.db.AddCredential(&dbModel.Credential{
		Base:          dbModel.Base{TenantID: info.TenantID},
		ConnectionID:  info.ConnectionID,
		Role:          data.Role,
		SchemaID:      data.SchemaID,
		CredDefID:     data.CredDefID,
		Attributes:    data.Attributes,
		InitiatedByUs: data.InitiatedByUs,
	}))

	utils.LogMed().Infof("Add credential %s for tenant %s", credential.ID, info.TenantID)

	status := model.JobStatusWaiting
	if !data.InitiatedByUs {
		status = model.JobStatusPending
	}

	try.To(l.AddJob(&dbModel.Job{
		Base:                 dbModel.Base{ID: info.JobID, TenantID: info.TenantID},
		ConnectionID:         &info.ConnectionID,
		ProtocolType:         model.ProtocolTypeCredential,
		ProtocolCredentialID: &credential.ID,
		InitiatedByUs:        data.InitiatedByUs,
		Status:               status,
		Result:               model.JobResultNone,
	}, credential.Description()))
	return nil
}

func (l *Listener) UpdateCredential(info *agency.JobInfo, data *agency.CredentialUpdate) (err error) {
	defer err2.Handle(&err)

	job := try.To1(l.db.GetJob(info.JobID, info.TenantID))

	utils.LogMed().Infof("Update credential %s for tenant %s", *job.ProtocolCredentialID, info.TenantID)

	credential := try.To1(l.db.GetCredential(*job.ProtocolCredentialID, job.TenantID))

	credential.Approved = utils.TSToTimeIfNotSet(&credential.Approved, data.ApprovedMs)
	credential.Issued = utils.TSToTimeIfNotSet(&credential.Issued, data.IssuedMs)
	credential.Failed = utils.TSToTimeIfNotSet(&credential.Failed, data.FailedMs)

	credential = try.To1(l.db.UpdateCredential(credential))

	job.Status, job.Result = getJobStatusForTimestamps(&credential.Approved, &credential.Issued, &credential.Failed)

	try.To(l.UpdateJob(job, credential.Description()))

	// Since we have new credential, check if any of the blocked proofs becomes unblocked
	if credential.IsIssued() {
		proofData := make([]*model.ProofAttribute, 0)
		for _, attribute := range credential.Attributes {
			proofData = append(
				proofData,
				&model.ProofAttribute{Name: attribute.Name, CredDefID: credential.CredDefID},
			)
		}
		if blockedJobs, err := l.db.GetOpenProofJobs(info.TenantID, proofData); err != nil {
			glog.Warningf("Encountered error fetching blocked jobs: %s", err)
		} else {
			for _, blockedJob := range blockedJobs {
				if err = l.updateBlockedProof(blockedJob); err != nil {
					glog.Error(err)
				}
			}
		}
	}
	return nil
}

func (l *Listener) isProvable(info *agency.JobInfo, data *dbModel.Proof) bool {
	attributes, err := l.db.SearchCredentials(info.TenantID, data.Attributes)
	provable := false
	if err == nil {
		provable = true
		for _, attr := range attributes {
			utils.LogMed().Infof("Attribute %s, cred count %d", attr.Attribute.Name, len(attr.Credentials))
			if len(attr.Credentials) == 0 {
				provable = false
				break
			}
		}
	} else {
		glog.Warningf("Encountered error when searching credentials: %s %s", info.TenantID, err)
	}
	return provable
}

func (l *Listener) AddProof(info *agency.JobInfo, data *agency.Proof) (err error) {
	defer err2.Handle(&err)

	newProof := &dbModel.Proof{
		Base:          dbModel.Base{TenantID: info.TenantID},
		ConnectionID:  info.ConnectionID,
		Role:          data.Role,
		Attributes:    data.Attributes,
		Result:        false,
		InitiatedByUs: data.InitiatedByUs,
	}

	if l.isProvable(info, newProof) {
		newProof.Provable = utils.CurrentTime()
	}

	proof := try.To1(l.db.AddProof(newProof))

	utils.LogMed().Infof("Add proof %s for tenant %s", proof.ID, info.TenantID)

	status := model.JobStatusWaiting
	if !data.InitiatedByUs {
		status = model.JobStatusPending
	}
	if newProof.Provable.IsZero() {
		status = model.JobStatusBlocked
	}

	try.To(l.AddJob(&dbModel.Job{
		Base:            dbModel.Base{ID: info.JobID, TenantID: info.TenantID},
		ConnectionID:    &info.ConnectionID,
		ProtocolType:    model.ProtocolTypeProof,
		ProtocolProofID: &proof.ID,
		InitiatedByUs:   data.InitiatedByUs,
		Status:          status,
		Result:          model.JobResultNone,
	}, proof.Description()))
	return nil
}

func (l *Listener) updateBlockedProof(job *dbModel.Job) (err error) {
	defer err2.Handle(&err)

	utils.LogMed().Infof("Update blocked proof %s for tenant %s", *job.ProtocolProofID, job.TenantID)

	proof := try.To1(l.db.GetProof(*job.ProtocolProofID, job.TenantID))

	if l.isProvable(&agency.JobInfo{TenantID: job.TenantID, JobID: job.ID, ConnectionID: *job.ConnectionID}, proof) {
		proof.Provable = utils.CurrentTime()
		proof = try.To1(l.db.UpdateProof(proof))

		job.Status, job.Result = getJobStatusForProof(proof)

		try.To(l.UpdateJob(job, proof.Description()))
	} else {
		utils.LogMed().Infof("Skipping update for blocked proof %s for tenant %s", *job.ProtocolProofID, job.TenantID)
	}

	return nil
}

func (l *Listener) UpdateProof(info *agency.JobInfo, data *agency.ProofUpdate) (err error) {
	defer err2.Handle(&err)

	job := try.To1(l.db.GetJob(info.JobID, info.TenantID))

	utils.LogMed().Infof("Update proof %s for tenant %s", *job.ProtocolProofID, info.TenantID)

	proof := try.To1(l.db.GetProof(*job.ProtocolProofID, job.TenantID))

	proof.Approved = utils.TSToTimeIfNotSet(&proof.Approved, data.ApprovedMs)
	proof.Verified = utils.TSToTimeIfNotSet(&proof.Verified, data.VerifiedMs)
	proof.Failed = utils.TSToTimeIfNotSet(&proof.Verified, data.VerifiedMs)

	if !proof.Verified.IsZero() {
		// TODO: these values should come from agency
		// now we just pick first found value and actually only guessing what core agency has picked
		provableAttrs := try.To1(l.db.SearchCredentials(proof.TenantID, proof.Attributes))
		proof.Values = make([]*model.ProofValue, 0)
		for _, attr := range provableAttrs {
			if len(attr.Credentials) > 0 {
				proof.Values = append(proof.Values, &model.ProofValue{
					ID:          attr.ID,
					AttributeID: attr.ID,
					Value:       attr.Credentials[0].Value,
				})
			}
		}
	}

	proof = try.To1(l.db.UpdateProof(proof))

	job.Status, job.Result = getJobStatusForProof(proof)

	try.To(l.UpdateJob(job, proof.Description()))
	return nil
}

func getJobStatusForTimestamps(approved, completed, failed *time.Time) (status model.JobStatus, result model.JobResult) {
	status = model.JobStatusWaiting
	result = model.JobResultNone
	if failed != nil && !failed.IsZero() {
		status = model.JobStatusComplete
		result = model.JobResultFailure
	} else if (approved == nil || approved.IsZero()) && (completed == nil || completed.IsZero()) {
		status = model.JobStatusPending
	} else if completed != nil && !completed.IsZero() {
		status = model.JobStatusComplete
		result = model.JobResultSuccess
	}
	return
}

func getJobStatusForProof(proof *dbModel.Proof) (status model.JobStatus, result model.JobResult) {
	status, result = getJobStatusForTimestamps(&proof.Approved, &proof.Verified, &proof.Failed)
	if status == model.JobStatusPending && proof.Provable.IsZero() {
		status = model.JobStatusBlocked
	}
	return
}

func (l *Listener) FailJob(info *agency.JobInfo) (err error) {
	defer err2.Handle(&err)

	job := try.To1(l.db.GetJob(info.JobID, info.TenantID))

	utils.LogMed().Infof("Fail job %s for tenant %s", job.ID, info.TenantID)
	job.Status = model.JobStatusComplete
	job.Result = model.JobResultFailure

	try.To(l.UpdateJob(job, fmt.Sprintf("Protocol %s failed", job.ProtocolType.String())))
	return nil
}
