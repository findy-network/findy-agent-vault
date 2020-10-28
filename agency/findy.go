// +build findy

package agency

import (
	"bytes"
	"os"

	"github.com/lainio/err2"

	"github.com/findy-network/findy-agent/cmds/onboard"

	"github.com/findy-network/findy-agent/cmds"
	"github.com/findy-network/findy-agent/cmds/agent"
)

const (
	walletName = "findy-agent-vault"
	walletKey  = "9C5qFG3grXfU9LodHdMop7CNVb3HtKddjgRc7oK5KhWY"
	agencyURL  = "http://localhost:8080"
)

type Findy struct{}

var Instance Agency = &Findy{}

func walletCmd() *cmds.Cmd {
	return &cmds.Cmd{
		WalletName: walletName,
		WalletKey:  walletKey,
	}
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

	// ping
	_, err = cmd.Exec(os.Stdout)
	if err != nil {
		panic(err)
	}
}

func (f *Findy) Invite() (invitation string, err error) {
	defer err2.Return(&err)

	cmd := agent.InvitationCmd{
		Cmd:  *walletCmd(),
		Name: walletName,
	}
	err = cmd.Validate()
	err2.Check(err)

	buf := new(bytes.Buffer)
	_, err = cmd.Exec(buf)
	err2.Check(err)

	invitation = buf.String()
	return
}

func (f *Findy) Connect() (string, error) {
	return "", nil
}
