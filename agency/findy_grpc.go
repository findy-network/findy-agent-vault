// +build findy_grpc

package agency

import (
	"fmt"

	"github.com/findy-network/findy-agent/grpc/client"
	"github.com/lainio/err2"
)

const (
	agencyHost = "localhost"
	agencyPort = 50051
)

type FindyGrpc struct {
	listener Listener
}

func grpcClient(a *Agent) client.Conn {
	baseCfg := client.BuildClientConnBase("", agencyHost, agencyPort, nil)
	return client.TryAuthOpen(a.JwtRaw, baseCfg)
}

var Instance Agency = &FindyGrpc{}

func (f *FindyGrpc) Init(l Listener) {
}

func (f *FindyGrpc) Invite(a *Agent) (invitation, id string, err error) {
	return
}

func (f *FindyGrpc) Connect(a *Agent, invitation string) (id string, err error) {
	defer err2.Return(&err) // TODO: do not leak internal errors to client

	conn := grpcClient(ctx)
	defer conn.Close()

	/* TODO:
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()*/

	connID, ch, err := conn.Connection(ctx, invitation)
	err2.Check(err)
	for status := range ch {
		fmt.Printf("Connection status: %s|%s: %s\n", connID, status.ProtocolId, status.State)
	}

	return
}

func (f *FindyGrpc) SendMessage(a *Agent, connectionID, message string) (id string, err error) {
	return
}

func (f *FindyGrpc) ResumeCredentialOffer(a *Agent, id string, accept bool) (err error) {
	return
}

func (f *FindyGrpc) ResumeProofRequest(a *Agent, id string, accept bool) (err error) {
	return
}
