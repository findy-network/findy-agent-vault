package listen

import (
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

func (l *Listener) AddConnection(info *agency.JobInfo, data *agency.Connection) {
	defer err2.Catch(func(err error) {
		glog.Errorf("Encountered error when adding connection %s", err.Error())
	})

	utils.LogMed().Infof("Add connection %s for tenant %s", info.ConnectionID, info.TenantID)

	job, err := l.db.GetJob(info.JobID, info.TenantID)
	err2.Check(err)

	now := utils.CurrentTime()

	connection, err := l.db.AddConnection(
		dbModel.NewConnection(info.ConnectionID, info.TenantID, &dbModel.Connection{
			OurDid:        data.OurDID,
			TheirDid:      data.TheirDID,
			TheirEndpoint: data.TheirEndpoint,
			TheirLabel:    data.TheirLabel,
			Approved:      &now, // TODO: get approved from agency
			Invited:       job.InitiatedByUs,
		}))
	err2.Check(err)

	job.ConnectionID = &connection.ID
	job.ProtocolConnectionID = &connection.ID
	job.Status = model.JobStatusComplete
	job.Result = model.JobResultSuccess

	err2.Check(l.UpdateJob(
		job,
		"Established connection to "+connection.TheirLabel,
	))
}

func (l *Listener) AddMessage(info *agency.JobInfo, data *agency.Message) {
	defer err2.Catch(func(err error) {
		glog.Errorf("Encountered error when adding message %s", err.Error())
	})
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
}

func (l *Listener) UpdateMessage(info *agency.JobInfo, _ *agency.MessageUpdate) {
	// TODO
}

func (l *Listener) AddCredential(info *agency.JobInfo, data *agency.Credential) {
	defer err2.Catch(func(err error) {
		glog.Errorf("Encountered error when adding credential %s", err.Error())
	})
	credential, err := l.db.AddCredential(dbModel.NewCredential(info.TenantID, &dbModel.Credential{
		ConnectionID:  info.ConnectionID,
		Role:          data.Role,
		SchemaID:      data.SchemaID,
		CredDefID:     data.CredDefID,
		Attributes:    data.Attributes,
		InitiatedByUs: data.InitiatedByUs,
	}))
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
}

func (l *Listener) UpdateCredential(info *agency.JobInfo, data *agency.CredentialUpdate) {
	defer err2.Catch(func(err error) {
		glog.Errorf("Encountered error when updating credential %s", err.Error())
	})

	job, err := l.db.GetJob(info.JobID, info.TenantID)
	err2.Check(err)

	utils.LogMed().Infof("Update credential %s for tenant %s", *job.ProtocolCredentialID, info.TenantID)

	credential, err := l.db.GetCredential(*job.ProtocolCredentialID, job.TenantID)
	err2.Check(err)

	if credential.Approved == nil {
		credential.Approved = utils.TimestampToTime(data.ApprovedMs)
	}
	if credential.Issued == nil {
		credential.Issued = utils.TimestampToTime(data.IssuedMs)
	}
	if credential.Failed == nil {
		credential.Failed = utils.TimestampToTime(data.FailedMs)
	}

	credential, err = l.db.UpdateCredential(credential)
	err2.Check(err)

	job.Status, job.Result = getJobStatusForTimestamps(credential.Approved, credential.Issued, credential.Failed)

	err2.Check(l.UpdateJob(job, credential.Description()))
}

func (l *Listener) AddProof(info *agency.JobInfo, data *agency.Proof) {
	defer err2.Catch(func(err error) {
		glog.Errorf("Encountered error when adding proof %s", err.Error())
	})

	proof, err := l.db.AddProof(dbModel.NewProof(info.TenantID, &dbModel.Proof{
		ConnectionID:  info.ConnectionID,
		Role:          data.Role,
		Attributes:    data.Attributes,
		Result:        false,
		InitiatedByUs: data.InitiatedByUs,
	}))
	err2.Check(err)

	utils.LogMed().Infof("Add proof %s for tenant %s", proof.ID, info.TenantID)

	status := model.JobStatusWaiting
	if !data.InitiatedByUs {
		status = model.JobStatusPending
	}

	err2.Check(l.AddJob(dbModel.NewJob(info.JobID, info.TenantID, &dbModel.Job{
		ConnectionID:    &info.ConnectionID,
		ProtocolType:    model.ProtocolTypeProof,
		ProtocolProofID: &proof.ID,
		InitiatedByUs:   data.InitiatedByUs,
		Status:          status,
		Result:          model.JobResultNone,
	}), proof.Description()))
}

func (l *Listener) UpdateProof(info *agency.JobInfo, data *agency.ProofUpdate) {
	defer err2.Catch(func(err error) {
		glog.Errorf("Encountered error when updating proof %s", err.Error())
	})
	job, err := l.db.GetJob(info.JobID, info.TenantID)
	err2.Check(err)

	utils.LogMed().Infof("Update proof %s for tenant %s", *job.ProtocolProofID, info.TenantID)

	proof, err := l.db.GetProof(*job.ProtocolProofID, job.TenantID)
	err2.Check(err)

	if proof.Approved == nil {
		proof.Approved = utils.TimestampToTime(data.ApprovedMs)
	}
	if proof.Verified == nil {
		proof.Verified = utils.TimestampToTime(data.VerifiedMs)
	}
	if proof.Failed == nil {
		proof.Failed = utils.TimestampToTime(data.FailedMs)
	}

	proof, err = l.db.UpdateProof(proof)
	err2.Check(err)

	job.Status, job.Result = getJobStatusForTimestamps(proof.Approved, proof.Verified, proof.Failed)

	err2.Check(l.UpdateJob(job, proof.Description()))
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
