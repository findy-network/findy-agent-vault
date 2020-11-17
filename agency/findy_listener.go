// +build findy

package agency

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

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

// TODO: use IDL/findy-agent types
func (f *Findy) findyCallback(pl *mesg.Payload) (while bool, err error) {
	defer err2.Return(&err)

	switch pl.Type {
	case pltype.CANotifyStatus:
		var status prot.TaskStatus

		err = mapstructure.Decode(pl.Message.Body, &status)
		err2.Check(err)

		switch status.Type {
		case pltype.AriesProtocolConnection:
			var c statusPairwise
			err = mapstructure.Decode(status.Payload, &c)
			err2.Check(err)

			f.listener.AddConnection(c.Name, c.MyDID, c.TheirDID, c.TheirEndpoint, c.TheirLabel)
		case pltype.ProtocolBasicMessage:
			var m statusBasicMessage
			err = mapstructure.Decode(status.Payload, &m)
			err2.Check(err)

			id := f.taskMapper.read(status.ID)
			f.listener.AddMessage(m.PwName, id, m.Message, m.SentByMe)
		}
	default:
		fmt.Println(dto.ToJSON(pl))
	}
	return true, nil
}
