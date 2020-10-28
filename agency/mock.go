// +build !findy

package agency

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

type Mock struct{}

var Instance Agency = &Mock{}

func (m *Mock) Init() {
}

func (m *Mock) Invite() (invitation string, err error) {
	return "", nil
}

func (m *Mock) Connect() (string, error) {
	return "", nil
}
