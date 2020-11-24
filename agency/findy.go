// +build findy

package agency

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/lainio/err2"

	"github.com/findy-network/findy-agent/agent/cloud"
	"github.com/findy-network/findy-agent/agent/mesg"
	"github.com/findy-network/findy-agent/agent/pltype"
	"github.com/findy-network/findy-agent/agent/ssi"
	"github.com/findy-network/findy-agent/agent/utils"
	"github.com/findy-network/findy-agent/client"
	"github.com/findy-network/findy-agent/cmds"
	"github.com/findy-network/findy-agent/cmds/agent"
	"github.com/findy-network/findy-agent/cmds/onboard"
	didexchange "github.com/findy-network/findy-agent/std/didexchange/invitation"
)

const (
	walletName = "findy-agent-vault"
	walletKey  = "9C5qFG3grXfU9LodHdMop7CNVb3HtKddjgRc7oK5KhWY"
	agencyURL  = "http://localhost:8080"
)

// TODO: use db for this
type mapper struct {
	sync.RWMutex
	taskToId map[string]string
	idToTask map[string]string
}

func (m *mapper) write(taskID, id string) {
	m.Lock()
	defer m.Unlock()
	m.taskToId[taskID] = id
	m.idToTask[id] = taskID
}

func (m *mapper) readID(taskID string) (id string) {
	m.RLock()
	defer m.RUnlock()
	id = m.taskToId[taskID]
	return
}

func (m *mapper) readTask(id string) (taskID string) {
	m.RLock()
	defer m.RUnlock()
	taskID = m.idToTask[id]
	return
}

type Findy struct {
	listener Listener
	agent    *cloud.Agent
	client   *client.Client
	endpoint string

	taskMapper *mapper
}

var Instance Agency = &Findy{}

func walletCmd() *cmds.Cmd {
	return &cmds.Cmd{
		WalletName: walletName,
		WalletKey:  walletKey,
	}
}

// TODO: do not onboard here, instead use JWT for authentication to agency
func (f *Findy) Init(l Listener) {
	f.taskMapper = &mapper{
		taskToId: make(map[string]string),
		idToTask: make(map[string]string),
	}
	cmd := agent.PingCmd{Cmd: *walletCmd()}

	err := cmd.Validate()
	// Onboard if wallet is not found
	if err != nil {
		onboardCmd := onboard.Cmd{
			Cmd:        *walletCmd(),
			Email:      walletName + "email",
			AgencyAddr: agencyURL,
		}

		err = onboardCmd.Validate()
		if err != nil {
			panic(err)
		}

		_, err = onboardCmd.Exec(os.Stdout)
		if err != nil {
			panic(err)
		}
	}

	f.listener = l
	f.client = &client.Client{
		Wallet: ssi.NewRawWalletCfg(walletName, walletKey),
	}
	f.agent = cloud.NewTransportReadyEA(f.client.Wallet)

	f.client.SetAgent(f.agent)

	im, err := f.agent.Trans().Call(pltype.CAPingOwnCA, &mesg.Msg{})
	if err != nil {
		panic(err)
	}
	f.endpoint = im.Message.Endpoint

	glog.Info("starting listening loop")

	go func() {
		err2.Check(f.client.Listen(f.findyCallback))
	}()
}

// TODO: fetch constructed JSON from CA
func (f *Findy) Invite() (invitation, id string, err error) {
	defer err2.Return(&err)

	id = utils.UUID()
	inv := didexchange.Invitation{
		ID:              id,
		Type:            pltype.AriesConnectionInvitation,
		ServiceEndpoint: f.endpoint,
		RecipientKeys:   []string{f.agent.Tr.PayloadPipe().Out.VerKey()},
		Label:           walletName,
	}

	jsonBytes := err2.Bytes.Try(json.Marshal(&inv))
	invitation = string(jsonBytes)
	return
}

func (f *Findy) Connect(invitation string) (id string, err error) {
	defer err2.Return(&err)

	inv := didexchange.Invitation{}
	err2.Check(json.Unmarshal([]byte(invitation), &inv))

	_, err = f.agent.Trans().Call(pltype.CAPairwiseCreate, &mesg.Msg{
		Info:       walletName, // our label
		Invitation: &inv,
	})
	err2.Check(err)

	id = inv.ID
	return
}

func (f *Findy) SendMessage(connectionID, message string) (id string, err error) {
	defer err2.Return(&err)

	id = uuid.New().String()

	pl, err := f.agent.Trans().Call(pltype.BasicMessageSend, &mesg.Msg{
		Name: connectionID,
		Info: message,
	})
	err2.Check(err)

	f.taskMapper.write(pl.Message.ID, id)

	return
}

func (f *Findy) ResumeCredentialOffer(id string, accept bool) (err error) {
	defer err2.Return(&err)

	taskID := f.taskMapper.readTask(id)

	if taskID != "" {
		_, err = f.agent.Trans().Call(pltype.CAContinueIssueCredentialProtocol, &mesg.Msg{
			ID:    taskID,
			Ready: accept,
		})
		err2.Check(err)
	} else {
		err = fmt.Errorf("no task found with id %s", id)
	}
	return
}

func (f *Findy) ResumeProofRequest(id string, accept bool) (err error) {
	defer err2.Return(&err)

	taskID := f.taskMapper.readTask(id)

	if taskID != "" {
		_, err = f.agent.Trans().Call(pltype.CAContinuePresentProofProtocol, &mesg.Msg{
			ID:    taskID,
			Ready: accept,
		})
		err2.Check(err)
	} else {
		err = fmt.Errorf("no task found with id %s", id)
	}
	return
}
