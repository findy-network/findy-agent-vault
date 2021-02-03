package archive

import (
	agency "github.com/findy-network/findy-agent-vault/agency/model"
	"github.com/findy-network/findy-agent-vault/db/store"
)

type Archiver struct {
	db store.DB
}

func NewArchiver(db store.DB) *Archiver {
	return &Archiver{db}
}

func (a *Archiver) ArchiveConnection(info *agency.ArchiveInfo, connection *agency.Connection) {

}

func (a *Archiver) ArchiveMessage(info *agency.ArchiveInfo, message *agency.Message) {

}

func (a *Archiver) ArchiveCredential(info *agency.ArchiveInfo, credential *agency.Credential) {

}

func (a *Archiver) ArchiveProof(info *agency.ArchiveInfo, proof *agency.Proof) {

}
