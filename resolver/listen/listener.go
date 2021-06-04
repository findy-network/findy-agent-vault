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
)

type Listener struct {
	db store.DB
	*update.Updater
}

func NewListener(db store.DB, updater *update.Updater) *Listener {
	return &Listener{db, updater}
}

func (l *Listener) AddConnection(info *agency.JobInfo, data *agency.Connection) (err error) {
	defer err2.Return(&err)

	utils.LogMed().Infof("Add connection %s for tenant %s", info.ConnectionID, info.TenantID)

	job, err := l.db.GetJob(info.JobID, info.TenantID)
	err2.Check(err)

	now := utils.CurrentTime()

	connection, err := l.db.AddConnection(
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
		})
	err2.Check(err)

	job.ConnectionID = &connection.ID
	job.ProtocolConnectionID = &connection.ID
	job.Status = model.JobStatusComplete
	job.Result = model.JobResultSuccess

	err2.Check(l.UpdateJob(
		job,
		"Established connection to "+connection.TheirLabel,
	))
	return nil
}

func (l *Listener) AddMessage(info *agency.JobInfo, data *agency.Message) (err error) {
	defer err2.Return(&err)

	msg, err := l.db.AddMessage(dbModel.NewMessage(info.TenantID, &dbModel.Message{
		ConnectionID: info.ConnectionID,
		Message:      data.Message,
		SentByMe:     data.SentByMe,
	}))
	err2.Check(err)

	err2.Check(l.AddJob(dbModel.NewJob(info.JobID, info.TenantID, &dbModel.Job{
		ConnectionID:      &info.ConnectionID,
		ProtocolType:      model.ProtocolTypeBasicMessage,
		ProtocolMessageID: &msg.ID,
		InitiatedByUs:     data.SentByMe,
		Status:            model.JobStatusComplete,
		Result:            model.JobResultSuccess,
	}), msg.Description()))
	return nil
}

func (l *Listener) UpdateMessage(info *agency.JobInfo, _ *agency.MessageUpdate) (err error) {
	// TODO
	return nil
}

func (l *Listener) AddCredential(info *agency.JobInfo, data *agency.Credential) (err error) {
	defer err2.Return(&err)

	credential, err := l.db.AddCredential(&dbModel.Credential{
		Base:          dbModel.Base{TenantID: info.TenantID},
		ConnectionID:  info.ConnectionID,
		Role:          data.Role,
		SchemaID:      data.SchemaID,
		CredDefID:     data.CredDefID,
		Attributes:    data.Attributes,
		InitiatedByUs: data.InitiatedByUs,
	})
	err2.Check(err)

	utils.LogMed().Infof("Add credential %s for tenant %s", credential.ID, info.TenantID)

	status := model.JobStatusWaiting
	if !data.InitiatedByUs {
		status = model.JobStatusPending
	}

	err2.Check(l.AddJob(dbModel.NewJob(info.JobID, info.TenantID, &dbModel.Job{
		ConnectionID:         &info.ConnectionID,
		ProtocolType:         model.ProtocolTypeCredential,
		ProtocolCredentialID: &credential.ID,
		InitiatedByUs:        data.InitiatedByUs,
		Status:               status,
		Result:               model.JobResultNone,
	}), credential.Description()))
	return nil
}

func (l *Listener) UpdateCredential(info *agency.JobInfo, data *agency.CredentialUpdate) (err error) {
	defer err2.Return(&err)

	job, err := l.db.GetJob(info.JobID, info.TenantID)
	err2.Check(err)

	utils.LogMed().Infof("Update credential %s for tenant %s", *job.ProtocolCredentialID, info.TenantID)

	credential, err := l.db.GetCredential(*job.ProtocolCredentialID, job.TenantID)
	err2.Check(err)

	credential.Approved = utils.TSToTimeIfNotSet(&credential.Approved, data.ApprovedMs)
	credential.Issued = utils.TSToTimeIfNotSet(&credential.Issued, data.IssuedMs)
	credential.Failed = utils.TSToTimeIfNotSet(&credential.Failed, data.FailedMs)

	credential, err = l.db.UpdateCredential(credential)
	err2.Check(err)

	job.Status, job.Result = getJobStatusForTimestamps(&credential.Approved, &credential.Issued, &credential.Failed)

	err2.Check(l.UpdateJob(job, credential.Description()))

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
	defer err2.Return(&err)

	newProof := dbModel.NewProof(info.TenantID, &dbModel.Proof{
		ConnectionID:  info.ConnectionID,
		Role:          data.Role,
		Attributes:    data.Attributes,
		Result:        false,
		InitiatedByUs: data.InitiatedByUs,
	})

	var provableTime *time.Time
	if l.isProvable(info, newProof) {
		now := utils.CurrentTime()
		newProof.Provable = &now
	}

	proof, err := l.db.AddProof(newProof)
	err2.Check(err)

	utils.LogMed().Infof("Add proof %s for tenant %s", proof.ID, info.TenantID)

	status := model.JobStatusWaiting
	if !data.InitiatedByUs {
		status = model.JobStatusPending
	}
	if provableTime == nil {
		status = model.JobStatusBlocked
	}

	err2.Check(l.AddJob(dbModel.NewJob(info.JobID, info.TenantID, &dbModel.Job{
		ConnectionID:    &info.ConnectionID,
		ProtocolType:    model.ProtocolTypeProof,
		ProtocolProofID: &proof.ID,
		InitiatedByUs:   data.InitiatedByUs,
		Status:          status,
		Result:          model.JobResultNone,
	}), proof.Description()))
	return nil
}

func (l *Listener) updateBlockedProof(job *dbModel.Job) (err error) {
	defer err2.Return(&err)

	utils.LogMed().Infof("Update blocked proof %s for tenant %s", *job.ProtocolProofID, job.TenantID)

	proof, err := l.db.GetProof(*job.ProtocolProofID, job.TenantID)
	err2.Check(err)

	if l.isProvable(&agency.JobInfo{TenantID: job.TenantID, JobID: job.ID, ConnectionID: *job.ConnectionID}, proof) {
		now := utils.CurrentTime()
		proof.Provable = &now
		proof, err = l.db.UpdateProof(proof)
		err2.Check(err)

		job.Status, job.Result = getJobStatusForProof(proof)

		err2.Check(l.UpdateJob(job, proof.Description()))
	} else {
		utils.LogMed().Infof("Skipping update for blocked proof %s for tenant %s", *job.ProtocolProofID, job.TenantID)
	}

	return nil
}

func (l *Listener) UpdateProof(info *agency.JobInfo, data *agency.ProofUpdate) (err error) {
	defer err2.Return(&err)

	job, err := l.db.GetJob(info.JobID, info.TenantID)
	err2.Check(err)

	utils.LogMed().Infof("Update proof %s for tenant %s", *job.ProtocolProofID, info.TenantID)

	proof, err := l.db.GetProof(*job.ProtocolProofID, job.TenantID)
	err2.Check(err)

	proof.Approved = utils.TSToTimePtrIfNotSet(proof.Approved, data.ApprovedMs)
	proof.Verified = utils.TSToTimePtrIfNotSet(proof.Verified, data.VerifiedMs)
	proof.Failed = utils.TSToTimePtrIfNotSet(proof.Verified, data.VerifiedMs)

	if proof.Verified != nil {
		// TODO: these values should come from agency
		// now we just pick first found value and actually only guessing what core agency has picked
		var provableAttrs []*model.ProvableAttribute
		provableAttrs, err = l.db.SearchCredentials(proof.TenantID, proof.Attributes)
		err2.Check(err)
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

	proof, err = l.db.UpdateProof(proof)
	err2.Check(err)

	job.Status, job.Result = getJobStatusForProof(proof)

	err2.Check(l.UpdateJob(job, proof.Description()))
	return nil
}

func getJobStatusForTimestamps(approved, completed, failed *time.Time) (status model.JobStatus, result model.JobResult) {
	status = model.JobStatusWaiting
	result = model.JobResultNone
	if failed != nil {
		status = model.JobStatusComplete
		result = model.JobResultFailure
	} else if approved == nil && completed == nil {
		status = model.JobStatusPending
	} else if completed != nil {
		status = model.JobStatusComplete
		result = model.JobResultSuccess
	}
	return
}

func getJobStatusForProof(proof *dbModel.Proof) (status model.JobStatus, result model.JobResult) {
	status, result = getJobStatusForTimestamps(proof.Approved, proof.Verified, proof.Failed)
	if status == model.JobStatusPending && proof.Provable == nil {
		status = model.JobStatusBlocked
	}
	return
}

func (l *Listener) FailJob(info *agency.JobInfo) (err error) {
	defer err2.Return(&err)

	job, err := l.db.GetJob(info.JobID, info.TenantID)
	err2.Check(err)

	utils.LogMed().Infof("Fail job %s for tenant %s", job.ID, info.TenantID)
	job.Status = model.JobStatusComplete
	job.Result = model.JobResultFailure

	err2.Check(l.UpdateJob(job, fmt.Sprintf("Protocol %s failed", job.ProtocolType.String())))
	return nil
}
