package archive

import (
	agency "github.com/findy-network/findy-agent-vault/agency/model"
	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/db/store"
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
	err2.Check(err)

	// 1. Check if we have a job with the id -> if ok, mark done - success
	// if not -> create job
	// 2. If we had the connection -> mark it done
	// if not, create connection
	_, err = a.db.ArchiveConnection(
		model.NewConnection(info.ConnectionID, agent.TenantID, &model.Connection{
			OurDid:        data.OurDID,
			TheirDid:      data.TheirDID,
			TheirEndpoint: data.TheirEndpoint,
			TheirLabel:    data.TheirLabel,
			Invited:       info.InitiatedByUs,
		}))
	err2.Check(err)
}

func (a *Archiver) ArchiveMessage(info *agency.ArchiveInfo, data *agency.Message) {
	defer err2.Catch(func(err error) {
		glog.Errorf("Encountered error when archiving message %s", err)
	})

	agent, err := a.db.GetAgent(nil, &info.AgentID)
	err2.Check(err)

	utils.LogMed().Infof("Archiving message with connection id %s for tenant %s", info.ConnectionID, agent.TenantID)

	_, err = a.db.ArchiveMessage(model.NewMessage(agent.TenantID, &model.Message{
		ConnectionID: info.ConnectionID,
		Message:      data.Message,
		SentByMe:     data.SentByMe, // TODO: sent time
	}))
	err2.Check(err)
}

func (a *Archiver) ArchiveCredential(info *agency.ArchiveInfo, data *agency.Credential) {
	defer err2.Catch(func(err error) {
		glog.Errorf("Encountered error when archiving credential %s", err)
	})

	agent, err := a.db.GetAgent(nil, &info.AgentID)
	err2.Check(err)

	utils.LogMed().Infof("Archiving credential with connection id %s for tenant %s", info.ConnectionID, agent.TenantID)

	now := utils.CurrentTime()

	_, err = a.db.ArchiveCredential(model.NewCredential(agent.TenantID, &model.Credential{
		ConnectionID:  info.ConnectionID,
		Role:          data.Role,
		SchemaID:      data.SchemaID,
		CredDefID:     data.CredDefID,
		Attributes:    data.Attributes,
		InitiatedByUs: data.InitiatedByUs,
		Issued:        &now, // TODO: get actual issued time
	}))
	err2.Check(err)
}

func (a *Archiver) ArchiveProof(info *agency.ArchiveInfo, data *agency.Proof) {
	defer err2.Catch(func(err error) {
		glog.Errorf("Encountered error when archiving proof %s", err)
	})

	agent, err := a.db.GetAgent(nil, &info.AgentID)
	err2.Check(err)

	utils.LogMed().Infof("Archiving proof with connection id %s for tenant %s", info.ConnectionID, agent.TenantID)

	now := utils.CurrentTime()

	_, err = a.db.ArchiveProof(model.NewProof(agent.TenantID, &model.Proof{
		ConnectionID:  info.ConnectionID,
		Role:          data.Role,
		Attributes:    data.Attributes,
		Result:        true,
		InitiatedByUs: data.InitiatedByUs,
		Verified:      &now, // TODO: get actual verified time

	}))
	err2.Check(err)
}
