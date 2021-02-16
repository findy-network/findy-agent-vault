package archive

import (
	"fmt"

	agency "github.com/findy-network/findy-agent-vault/agency/model"
	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/db/store"
	graph "github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/golang/glog"
	"github.com/lainio/err2"
	"github.com/lainio/err2/assert"
)

type Archiver struct {
	db store.DB
}

func NewArchiver(db store.DB) *Archiver {
	return &Archiver{db}
}

// TODO: write unit tests
// TODO: should the archived flag be attached to the job instead of
// protocol object itself?
// TODO: should archived flag be visible somehow to the clients?
// e.g. removing/arhiving of incomplete jobs from the UI

func (a *Archiver) matchProtocol(job *model.Job) (idToUpdate **string, archive func(string, string) error, err error) {
	switch job.ProtocolType {
	case graph.ProtocolTypeConnection:
		idToUpdate = &job.ProtocolConnectionID
		archive = a.db.ArchiveConnection
	case graph.ProtocolTypeBasicMessage:
		idToUpdate = &job.ProtocolMessageID
		archive = a.db.ArchiveMessage
	case graph.ProtocolTypeCredential:
		idToUpdate = &job.ProtocolCredentialID
		archive = a.db.ArchiveCredential
	case graph.ProtocolTypeProof:
		idToUpdate = &job.ProtocolProofID
		archive = a.db.ArchiveProof
	default:
		return nil, nil, fmt.Errorf("invalid protocol type for job: %s", job.ProtocolType)
	}

	return idToUpdate, archive, nil
}

func (a *Archiver) archiveExisting(
	info *agency.ArchiveInfo,
	agent *model.Agent,
	job *model.Job,
) (err error) {
	defer err2.Return(&err)

	var (
		idToUpdate **string
		archive    func(string, string) error
	)

	utils.LogLow().Infof("Archive for existing job %s (%s)", job.ID, job.ProtocolType)

	idToUpdate, archive, err = a.matchProtocol(job)
	err2.Check(err)

	assert.P.True(*idToUpdate != nil, "existing job to archive should have a valid protocol id")

	// TODO: update data also?
	err2.Check(archive(**idToUpdate, agent.TenantID))

	if job.Status != graph.JobStatusComplete || job.Result != graph.JobResultSuccess {
		utils.LogLow().Infof("Update existing job %s (%s) when archiving", job.ID, job.ProtocolType)

		// update job
		job.Status = graph.JobStatusComplete
		job.Result = graph.JobResultSuccess
		job.ConnectionID = &info.ConnectionID
		_, err = a.db.UpdateJob(job)
		err2.Check(err)
	}
	return nil
}

func (a *Archiver) archiveNew(
	info *agency.ArchiveInfo,
	agent *model.Agent,
	protocolType graph.ProtocolType,
	addToStore func(*model.Agent, bool) (string, error),
) (err error) {
	defer err2.Return(&err)

	id, err := addToStore(agent, info.InitiatedByUs)
	err2.Check(err)

	job := model.NewJob(info.JobID, agent.TenantID, &model.Job{
		ConnectionID:  &info.ConnectionID,
		ProtocolType:  protocolType,
		InitiatedByUs: info.InitiatedByUs,
		Status:        graph.JobStatusComplete,
		Result:        graph.JobResultSuccess,
	})
	idToUpdate, _, err := a.matchProtocol(job)
	*idToUpdate = &id

	utils.LogLow().Infof("Create new job %s (%s) when archiving", job.ID, job.ProtocolType)

	// add job
	_, err = a.db.AddJob(job)
	err2.Check(err)

	return
}

func (a *Archiver) archive(
	info *agency.ArchiveInfo,
	protocolType graph.ProtocolType,
	addToStore func(*model.Agent, bool) (string, error),
) (err error) {
	defer err2.Return(&err)

	agent, err := a.db.GetAgent(nil, &info.AgentID)
	err2.Check(err)

	job, err := a.db.GetJob(info.JobID, agent.TenantID)

	if err == nil { // job exists
		err2.Check(a.archiveExisting(info, agent, job))
	} else if store.ErrorCode(err) == store.ErrCodeNotFound { // job is new
		err2.Check(a.archiveNew(info, agent, protocolType, addToStore))
	} else {
		err2.Check(err) // some other error
	}

	return
}

func (a *Archiver) ArchiveConnection(info *agency.ArchiveInfo, data *agency.Connection) {
	defer err2.Catch(func(err error) {
		glog.Errorf("Encountered error when archiving connection %s", err)
	})

	err2.Check(a.archive(info, graph.ProtocolTypeConnection, func(agent *model.Agent, initiatedByUs bool) (id string, err error) {
		defer err2.Return(&err)

		now := utils.CurrentTime()
		connection, err := a.db.AddConnection(model.NewConnection(info.ConnectionID, agent.TenantID, &model.Connection{
			OurDid:        data.OurDID,
			TheirDid:      data.TheirDID,
			TheirEndpoint: data.TheirEndpoint,
			TheirLabel:    data.TheirLabel,
			Approved:      &now, // TODO: get approved from agency
			Invited:       initiatedByUs,
			Archived:      &now,
		}))
		err2.Check(err)

		return connection.ID, nil
	}))
}

func (a *Archiver) ArchiveMessage(info *agency.ArchiveInfo, data *agency.Message) {
	defer err2.Catch(func(err error) {
		glog.Errorf("Encountered error when archiving message %s", err)
	})
	err2.Check(a.archive(info, graph.ProtocolTypeBasicMessage, func(agent *model.Agent, initiatedByUs bool) (id string, err error) {
		defer err2.Return(&err)

		now := utils.CurrentTime()
		message, err := a.db.AddMessage(model.NewMessage(agent.TenantID, &model.Message{
			ConnectionID: info.ConnectionID,
			Message:      data.Message,
			SentByMe:     data.SentByMe, // TODO: sent time
			Archived:     &now,
		}))
		err2.Check(err)

		return message.ID, nil
	}))
}

func (a *Archiver) ArchiveCredential(info *agency.ArchiveInfo, data *agency.Credential) {
	defer err2.Catch(func(err error) {
		glog.Errorf("Encountered error when archiving credential %s", err)
	})

	err2.Check(a.archive(info, graph.ProtocolTypeCredential, func(agent *model.Agent, initiatedByUs bool) (id string, err error) {
		defer err2.Return(&err)

		now := utils.CurrentTime()
		credential, err := a.db.AddCredential(model.NewCredential(agent.TenantID, &model.Credential{
			ConnectionID:  info.ConnectionID,
			Role:          data.Role,
			SchemaID:      data.SchemaID,
			CredDefID:     data.CredDefID,
			Attributes:    data.Attributes,
			InitiatedByUs: data.InitiatedByUs,
			Issued:        &now, // TODO: get actual issued time
			Archived:      &now,
		}))
		err2.Check(err)

		return credential.ID, nil
	}))
}

func (a *Archiver) ArchiveProof(info *agency.ArchiveInfo, data *agency.Proof) {
	defer err2.Catch(func(err error) {
		glog.Errorf("Encountered error when archiving proof %s", err)
	})

	err2.Check(a.archive(info, graph.ProtocolTypeProof, func(agent *model.Agent, initiatedByUs bool) (id string, err error) {
		defer err2.Return(&err)

		now := utils.CurrentTime()
		proof, err := a.db.AddProof(model.NewProof(agent.TenantID, &model.Proof{
			ConnectionID:  info.ConnectionID,
			Role:          data.Role,
			Attributes:    data.Attributes,
			Result:        true,
			InitiatedByUs: data.InitiatedByUs,
			Verified:      &now, // TODO: get actual verified time

		}))
		err2.Check(err)

		return proof.ID, nil
	}))
}
