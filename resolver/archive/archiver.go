package archive

import (
	agency "github.com/findy-network/findy-agent-vault/agency/model"
	"github.com/findy-network/findy-agent-vault/db/store"
	"github.com/golang/glog"
)

type Archiver struct {
	db store.DB
}

func NewArchiver(db store.DB) *Archiver {
	return &Archiver{db}
}

func (a *Archiver) ArchiveConnection(info *agency.ArchiveInfo, connection *agency.Connection) {
	glog.Infof("ArchiveConnection NOT IMPLEMENTED %s", info.ConnectionID)
}

func (a *Archiver) ArchiveMessage(info *agency.ArchiveInfo, message *agency.Message) {
	glog.Infof("ArchiveMessage NOT IMPLEMENTED %s", info.ConnectionID)
}

func (a *Archiver) ArchiveCredential(info *agency.ArchiveInfo, credential *agency.Credential) {
	glog.Infof("ArchiveCredential NOT IMPLEMENTED %s", info.ConnectionID)
}

func (a *Archiver) ArchiveProof(info *agency.ArchiveInfo, proof *agency.Proof) {
	glog.Infof("ArchiveProof NOT IMPLEMENTED %s", info.ConnectionID)
}
