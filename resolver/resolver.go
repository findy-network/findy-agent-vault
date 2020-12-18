package resolver

import (
	agencys "github.com/findy-network/findy-agent-vault/agency"
	agency "github.com/findy-network/findy-agent-vault/agency/model"
	"github.com/findy-network/findy-agent-vault/db/fake"
	dbModel "github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/db/store"
	"github.com/findy-network/findy-agent-vault/db/store/mock"
	"github.com/findy-network/findy-agent-vault/db/store/pg"
	"github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/paginator"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/lainio/err2"

	"github.com/golang/glog"
)

//go:generate go run github.com/99designs/gqlgen

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	db               store.DB
	agency           agency.Agency
	eventSubscribers *subscriberRegister
}

func InitResolver(mockDB, fakeData bool) *Resolver {
	var db store.DB
	if mockDB {
		db = mock.InitState()
	} else {
		db = pg.InitDB("file://db/migrations", "5432", false)
	}

	nextPage := true
	after := uint64(0)
	allAgents := make([]*dbModel.Agent, 0)
	for nextPage {
		agents, err := db.GetListenerAgents(&paginator.BatchInfo{Count: 50, After: after})
		if err != nil {
			panic(err)
		}
		allAgents = append(allAgents, agents.Agents...)
		nextPage = agents.HasNextPage
		after = agents.Agents[len(agents.Agents)-1].Cursor
	}

	listenerAgents := make([]*agency.Agent, len(allAgents))
	for index := range allAgents {
		listenerAgents[index] = agencyAuth(allAgents[index])
	}

	r := &Resolver{
		db:               db,
		eventSubscribers: newSubscriberRegister(),
	}
	r.agency = agencys.InitAgency(agencys.AgencyTypeFindyGRPC, r, listenerAgents)

	if fakeData {
		fake.AddData(db)
	}

	return r
}

func (r *Resolver) addEvent(tenantID string, job *dbModel.Job, description string) (err error) {
	defer err2.Return(&err)
	var connectionID, jobID *string
	if job != nil {
		connectionID = job.ConnectionID
		jobID = &job.ID
	}
	event, err := r.db.AddEvent(dbModel.NewEvent(tenantID, &dbModel.Event{
		Read:         false,
		Description:  description,
		ConnectionID: connectionID,
		JobID:        jobID,
	}))
	err2.Check(err)

	r.eventSubscribers.notify(tenantID, event)
	return err
}

func (r *Resolver) addJob(job *dbModel.Job, description string) (err error) {
	defer err2.Return(&err)

	utils.LogLow().Infof("Add job with ID %s for tenant %s", job.ID, job.TenantID)

	job, err = r.db.AddJob(job)
	err2.Check(err)

	err2.Check(r.addEvent(job.TenantID, job, description))

	return
}

func (r *Resolver) updateJob(job *dbModel.Job, description string) (err error) {
	defer err2.Return(&err)
	job, err = r.db.UpdateJob(job)
	err2.Check(err)

	err2.Check(r.addEvent(job.TenantID, job, description))

	return
}

func (r *Resolver) AddConnection(info *agency.JobInfo, ourDID, theirDID, theirEndpoint, theirLabel string) {
	defer err2.Catch(func(err error) {
		glog.Errorf("Encountered error when adding connection %s", err.Error())
	})

	utils.LogLow().Infof("Add connection %s for tenant %s", info.ConnectionID, info.TenantID)

	job, err := r.db.GetJob(info.JobID, info.TenantID)
	err2.Check(err)

	now := utils.CurrentTime()

	connection, err := r.db.AddConnection(
		dbModel.NewConnection(info.ConnectionID, info.TenantID, &dbModel.Connection{
			OurDid:        ourDID,
			TheirDid:      theirDID,
			TheirEndpoint: theirEndpoint,
			TheirLabel:    theirLabel,
			Approved:      &now, // TODO: get approved from agency
			Invited:       job.InitiatedByUs,
		}))
	err2.Check(err)

	job.ConnectionID = &connection.ID
	job.ProtocolConnectionID = &connection.ID
	job.Status = model.JobStatusComplete
	job.Result = model.JobResultSuccess

	err2.Check(r.updateJob(
		job,
		"Established connection to "+connection.TheirLabel,
	))
}

func (r *Resolver) AddMessage(info *agency.JobInfo, message string, sentByMe bool) {
	defer err2.Catch(func(err error) {
		glog.Errorf("Encountered error when adding message %s", err.Error())
	})
	msg, err := r.db.AddMessage(dbModel.NewMessage(info.TenantID, &dbModel.Message{
		ConnectionID: info.ConnectionID,
		Message:      message,
		SentByMe:     sentByMe,
	}))
	err2.Check(err)

	err2.Check(r.addJob(dbModel.NewJob(info.JobID, info.TenantID, &dbModel.Job{
		ConnectionID:      &info.ConnectionID,
		ProtocolType:      model.ProtocolTypeBasicMessage,
		ProtocolMessageID: &msg.ID,
		InitiatedByUs:     sentByMe,
		Status:            model.JobStatusComplete,
		Result:            model.JobResultSuccess,
	}), msg.Description()))
}

func (r *Resolver) UpdateMessage(info *agency.JobInfo, delivered bool) {
	// TODO
}

func (r *Resolver) AddCredential(
	info *agency.JobInfo,
	role model.CredentialRole,
	schemaID, credDefID string,
	attributes []*model.CredentialValue,
	initiatedByUs bool,
) {
	defer err2.Catch(func(err error) {
		glog.Errorf("Encountered error when adding credential %s", err.Error())
	})
	credential, err := r.db.AddCredential(dbModel.NewCredential(info.TenantID, &dbModel.Credential{
		ConnectionID:  info.ConnectionID,
		Role:          role,
		SchemaID:      schemaID,
		CredDefID:     credDefID,
		Attributes:    attributes,
		InitiatedByUs: initiatedByUs,
	}))
	err2.Check(err)

	status := model.JobStatusWaiting
	if !initiatedByUs {
		status = model.JobStatusPending
	}

	err2.Check(r.addJob(dbModel.NewJob(info.JobID, info.TenantID, &dbModel.Job{
		ConnectionID:         &info.ConnectionID,
		ProtocolType:         model.ProtocolTypeCredential,
		ProtocolCredentialID: &credential.ID,
		InitiatedByUs:        initiatedByUs,
		Status:               status,
		Result:               model.JobResultNone,
	}), credential.Description()))
}

func (r *Resolver) UpdateCredential(info *agency.JobInfo, approvedMs, issuedMs, failedMs *int64) {
	defer err2.Catch(func(err error) {
		glog.Errorf("Encountered error when updating credential %s", err.Error())
	})
	job, err := r.db.GetJob(info.JobID, info.TenantID)
	err2.Check(err)

	// TODO: is this needed - can we just update directly
	credential, err := r.db.GetCredential(*job.ProtocolCredentialID, job.TenantID)
	err2.Check(err)

	credential.Approved = utils.TimestampToTime(approvedMs)
	credential.Issued = utils.TimestampToTime(issuedMs)
	credential.Failed = utils.TimestampToTime(failedMs)

	credential, err = r.db.UpdateCredential(credential)
	err2.Check(err)

	status := model.JobStatusWaiting
	result := model.JobResultNone
	if failedMs != nil {
		status = model.JobStatusComplete
		result = model.JobResultFailure
	} else if approvedMs == nil && issuedMs == nil {
		status = model.JobStatusPending
	} else if issuedMs != nil {
		status = model.JobStatusComplete
		result = model.JobResultSuccess
	}

	job.Status = status
	job.Result = result

	err2.Check(r.updateJob(job, credential.Description()))
}

func (r *Resolver) AddProof(info *agency.JobInfo, role model.ProofRole, attributes []*model.ProofAttribute, initiatedByUs bool) {
	defer err2.Catch(func(err error) {
		glog.Errorf("Encountered error when adding proof %s", err.Error())
	})

	proof, err := r.db.AddProof(dbModel.NewProof(info.TenantID, &dbModel.Proof{
		ConnectionID:  info.ConnectionID,
		Role:          role,
		Attributes:    attributes,
		Result:        false,
		InitiatedByUs: initiatedByUs,
	}))
	err2.Check(err)

	status := model.JobStatusWaiting
	if !initiatedByUs {
		status = model.JobStatusPending
	}

	err2.Check(r.addJob(dbModel.NewJob(info.JobID, info.TenantID, &dbModel.Job{
		ConnectionID:    &info.ConnectionID,
		ProtocolType:    model.ProtocolTypeProof,
		ProtocolProofID: &proof.ID,
		InitiatedByUs:   initiatedByUs,
		Status:          status,
		Result:          model.JobResultNone,
	}), proof.Description()))
}

func (r *Resolver) UpdateProof(info *agency.JobInfo, approvedMs, verifiedMs, failedMs *int64) {
	defer err2.Catch(func(err error) {
		glog.Errorf("Encountered error when updating proof %s", err.Error())
	})
	job, err := r.db.GetJob(info.JobID, info.TenantID)
	err2.Check(err)

	// TODO: is this needed - can we just update directly
	proof, err := r.db.GetProof(*job.ProtocolProofID, job.TenantID)
	err2.Check(err)

	proof.Approved = utils.TimestampToTime(approvedMs)
	proof.Verified = utils.TimestampToTime(verifiedMs)
	proof.Failed = utils.TimestampToTime(failedMs)

	proof, err = r.db.UpdateProof(proof)
	err2.Check(err)

	status := model.JobStatusWaiting
	result := model.JobResultNone
	if failedMs != nil {
		status = model.JobStatusComplete
		result = model.JobResultFailure
	} else if approvedMs == nil && verifiedMs == nil {
		status = model.JobStatusPending
	} else if verifiedMs != nil {
		status = model.JobStatusComplete
		result = model.JobResultSuccess
	}

	job.Status = status
	job.Result = result

	err2.Check(r.updateJob(job, proof.Description()))
}
