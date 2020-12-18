// +build findy

package legacy

import (
	"github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"

	"github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/tools/tools"
	"github.com/findy-network/findy-agent/agent/didcomm"
	"github.com/findy-network/findy-agent/agent/prot"

	"github.com/findy-network/findy-agent/agent/mesg"
	"github.com/findy-network/findy-agent/agent/pltype"
	"github.com/findy-network/findy-wrapper-go/dto"
	"github.com/lainio/err2"
)

// TODO: use from IDL/findy-agent
type statusPairwise struct {
	Name          string `json:"name"`
	MyDID         string `json:"myDid"`
	TheirDID      string `json:"theirDid"`
	TheirEndpoint string `json:"theirEndpoint"`
	TheirLabel    string `json:"theirLabel"`
}

type statusBasicMessage struct {
	PwName    string `json:"pairwise"`
	Message   string `json:"message"`
	SentByMe  bool   `json:"sentByMe"`
	Delivered bool   `json:"delivered"`
}

type statusIssueCredential struct {
	CredDefID  string                        `json:"credDefId"`
	SchemaID   string                        `json:"schemaId"`
	Attributes []didcomm.CredentialAttribute `json:"attributes"`
}

type statusPresentProof struct {
	Attributes []didcomm.ProofAttribute `json:"attributes"`
}

// TODO: use IDL/findy-agent types
func (f *Findy) findyCallback(pl *mesg.Payload) (while bool, err error) {
	defer err2.Return(&err) // TODO

	glog.Infof("Received findy callback %s %s", pl.Type)

	currentTime := utils.CurrentTimeMs()

	switch pl.Type {
	case pltype.CANotifyUserAction:
		fallthrough
	case pltype.CANotifyStatus:
		var status prot.TaskStatus

		err2.Check(mapstructure.Decode(pl.Message.Body, &status))

		glog.Infof("Callback payload status %s", status.Type)

		switch status.Type {
		case pltype.AriesProtocolConnection:
			var c statusPairwise
			err = mapstructure.Decode(status.Payload, &c)
			err2.Check(err)

			f.listener.AddConnection(c.Name, c.MyDID, c.TheirDID, c.TheirEndpoint, c.TheirLabel)

		case pltype.ProtocolBasicMessage:
			var m statusBasicMessage
			err2.Check(mapstructure.Decode(status.Payload, &m))

			id := f.taskMapper.readID(status.ID)
			sentByMe := true
			if id == "" {
				sentByMe = false
				id = uuid.New().String()
			}
			f.listener.AddMessage(status.Name, id, m.Message, sentByMe)

		case pltype.ProtocolIssueCredential:
			var c statusIssueCredential
			err2.Check(mapstructure.Decode(status.Payload, &c))

			// TODO: credential issuance initiated by holder
			if status.PendingUserAction {
				values := make([]*model.CredentialValue, 0)
				for _, v := range c.Attributes {
					values = append(values, &model.CredentialValue{
						Name:  v.Name,
						Value: v.Value,
					})
				}
				id := uuid.New().String()
				f.taskMapper.write(status.ID, id)
				f.listener.AddCredential(status.Name, id, model.CredentialRoleHolder, c.SchemaID, c.CredDefID, values, false)
			} else {
				// if ready -> what if fails
				id := f.taskMapper.readID(status.ID)
				f.listener.UpdateCredential(status.Name, id, nil, &currentTime, nil)
			}

		case pltype.ProtocolPresentProof:
			var p statusPresentProof
			err2.Check(mapstructure.Decode(status.Payload, &p))

			// TODO: proof initiated by prover
			if status.PendingUserAction {
				attributes := make([]*model.ProofAttribute, 0)
				for _, v := range p.Attributes {
					attributes = append(attributes, &model.ProofAttribute{
						Name:      v.Name,
						CredDefID: v.CredDefID,
					})
				}
				id := uuid.New().String()
				f.taskMapper.write(status.ID, id)
				f.listener.AddProof(status.Name, id, model.ProofRoleProver, attributes, false)
			} else {
				id := f.taskMapper.readID(status.ID)
				f.listener.UpdateProof(status.Name, id, nil, &currentTime, nil)
			}
		default:
			glog.Warning(dto.ToJSON(pl))

		}
	default:
		glog.Warning(dto.ToJSON(pl))
	}
	return true, nil
}
