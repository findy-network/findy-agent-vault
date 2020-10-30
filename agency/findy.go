// +build findy

package agency

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/golang/glog"
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
	"github.com/findy-network/findy-wrapper-go/dto"
)

const (
	walletName = "findy-agent-vault"
	walletKey  = "9C5qFG3grXfU9LodHdMop7CNVb3HtKddjgRc7oK5KhWY"
	agencyURL  = "http://localhost:8080"
)

type Findy struct {
	agent    *cloud.Agent
	client   *client.Client
	endpoint string
}

var Instance Agency = &Findy{}

func walletCmd() *cmds.Cmd {
	return &cmds.Cmd{
		WalletName: walletName,
		WalletKey:  walletKey,
	}
}

func (f *Findy) listenWs() {

	glog.Info("starting listening loop")

	go func() {
		err2.Check(f.client.Listen(func(pl *mesg.Payload) (while bool, err error) {
			defer err2.Return(&err)

			fmt.Println(dto.ToJSON(pl))
			return true, nil
		}))
	}()
}

// TODO: do not onboard here, instead use JWT for authentication to agency
func (f *Findy) Init() {
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

	f.listenWs()
}

// TODO: fetch constructed JSON from CA
func (f *Findy) Invite() (invitation string, err error) {
	defer err2.Return(&err)

	inv := didexchange.Invitation{
		ID:              utils.UUID(),
		Type:            pltype.AriesConnectionInvitation,
		ServiceEndpoint: f.endpoint,
		RecipientKeys:   []string{f.agent.Tr.PayloadPipe().Out.VerKey()},
		Label:           walletName,
	}

	jsonBytes := err2.Bytes.Try(json.Marshal(&inv))
	invitation = string(jsonBytes)
	return
}

func (f *Findy) Connect() (string, error) {
	return "", nil
}
