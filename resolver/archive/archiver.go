package archive

import (
	agency "github.com/findy-network/findy-agent-vault/agency/model"
	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/db/store"
	graph "github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/golang/glog"
	"github.com/lainio/err2"
)

type Archiver struct {
	db store.DB
}

func NewArchiver(db store.DB) *Archiver {
	return &Archiver{db}
}

func (a *Archiver) ArchiveConnection(info *agency.ArchiveInfo, data *agency.Connection) {
	defer err2.Catch(func(err error) {
		glog.Errorf("Encountered error when archiving connection %s", err)
	})

	agent, err := a.db.GetAgent(nil, &info.AgentID)
	err2.Check(err)

	utils.LogMed().Infof("Archiving connection %+v for tenant %s", data, agent.TenantID)

	job, err := a.db.GetJob(info.JobID, agent.TenantID)
	jobIsIncomplete := err == nil &&
		(job.Status != graph.JobStatusComplete ||
			job.Result != graph.JobResultSuccess ||
			job.ProtocolConnectionID == nil)

	if jobIsIncomplete {
		// update connection
		// TODO: update data also?
		err = a.db.ArchiveConnection(info.ConnectionID, agent.TenantID)
		err2.Check(err)

		// update job
		job.Status = graph.JobStatusComplete
		job.Result = graph.JobResultSuccess
		job.ConnectionID = &info.ConnectionID
		job.ProtocolConnectionID = &info.ConnectionID
		_, err = a.db.UpdateJob(job)
		err2.Check(err)
	} else if store.ErrorCode(err) == store.NotExists {
		now := utils.CurrentTime()

		// create connection
		var connection *model.Connection
		connection, err = a.db.AddConnection(model.NewConnection(info.ConnectionID, agent.TenantID, &model.Connection{
			OurDid:        data.OurDID,
			TheirDid:      data.TheirDID,
			TheirEndpoint: data.TheirEndpoint,
			TheirLabel:    data.TheirLabel,
			Approved:      &now, // TODO: get approved from agency
			Invited:       job.InitiatedByUs,
			Archived:      &now,
		}))
		err2.Check(err)

		// add job
		_, err = a.db.AddJob(model.NewJob(info.JobID, agent.TenantID, &model.Job{
			ConnectionID:         &connection.ID,
			ProtocolType:         graph.ProtocolTypeConnection,
			ProtocolConnectionID: &connection.ID,
			InitiatedByUs:        info.InitiatedByUs,
			Status:               graph.JobStatusComplete,
			Result:               graph.JobResultSuccess,
		}))
		err2.Check(err)
	} else {
		err2.Check(err)
	}
}

func (a *Archiver) ArchiveMessage(info *agency.ArchiveInfo, data *agency.Message) {
	defer err2.Catch(func(err error) {
		glog.Errorf("Encountered error when archiving message %s", err)
	})

	agent, err := a.db.GetAgent(nil, &info.AgentID)
	err2.Check(err)

	utils.LogMed().Infof("Archiving message with connection id %s for tenant %s", info.ConnectionID, agent.TenantID)

	/*	_, err = a.db.ArchiveMessage(model.NewMessage(agent.TenantID, &model.Message{
			ConnectionID: info.ConnectionID,
			Message:      data.Message,
			SentByMe:     data.SentByMe, // TODO: sent time
		}))
		err2.Check(err)*/
}

func (a *Archiver) ArchiveCredential(info *agency.ArchiveInfo, data *agency.Credential) {
	defer err2.Catch(func(err error) {
		glog.Errorf("Encountered error when archiving credential %s", err)
	})

	agent, err := a.db.GetAgent(nil, &info.AgentID)
	err2.Check(err)

	utils.LogMed().Infof("Archiving credential with connection id %s for tenant %s", info.ConnectionID, agent.TenantID)

	/*	now := utils.CurrentTime()

		_, err = a.db.ArchiveCredential(model.NewCredential(agent.TenantID, &model.Credential{
			ConnectionID:  info.ConnectionID,
			Role:          data.Role,
			SchemaID:      data.SchemaID,
			CredDefID:     data.CredDefID,
			Attributes:    data.Attributes,
			InitiatedByUs: data.InitiatedByUs,
			Issued:        &now, // TODO: get actual issued time
		}))
		err2.Check(err)*/
}

func (a *Archiver) ArchiveProof(info *agency.ArchiveInfo, data *agency.Proof) {
	defer err2.Catch(func(err error) {
		glog.Errorf("Encountered error when archiving proof %s", err)
	})

	agent, err := a.db.GetAgent(nil, &info.AgentID)
	err2.Check(err)

	utils.LogMed().Infof("Archiving proof with connection id %s for tenant %s", info.ConnectionID, agent.TenantID)

	/*now := utils.CurrentTime()

	_, err = a.db.ArchiveProof(model.NewProof(agent.TenantID, &model.Proof{
		ConnectionID:  info.ConnectionID,
		Role:          data.Role,
		Attributes:    data.Attributes,
		Result:        true,
		InitiatedByUs: data.InitiatedByUs,
		Verified:      &now, // TODO: get actual verified time

	}))
	err2.Check(err)*/
}
