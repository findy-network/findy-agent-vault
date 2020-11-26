// +build findy_grpc

package agency

import (
	"context"
	"fmt"

	"github.com/dgrijalva/jwt-go"
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

func getToken(ctx context.Context) string {
	defer err2.Catch(func(err error) {
		panic(err)
	})

	user := ctx.Value("user")
	if user == nil {
		err2.Check(fmt.Errorf("no authenticated user found"))
	}
	jwtToken, ok := user.(*jwt.Token)
	if !ok {
		err2.Check(fmt.Errorf("no jwt token found for user"))
	}
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		err2.Check(fmt.Errorf("no claims found for token"))
	}
	caDID, ok := claims["un"].(string)
	if !ok || caDID == "" {
		err2.Check(fmt.Errorf("no cloud agent DID found for token"))
	}
	if jwtToken.Raw == "" {
		err2.Check(fmt.Errorf("no raw token found for user %s", caDID))
	}

	return jwtToken.Raw
}

func grpcClient(ctx context.Context) client.Conn {
	baseCfg := client.BuildClientConnBase("", agencyHost, agencyPort, nil)
	return client.TryAuthOpen(getToken(ctx), baseCfg)
}

var Instance Agency = &FindyGrpc{}

func (f *FindyGrpc) Init(l Listener) {
}

func (f *FindyGrpc) Invite(ctx context.Context) (invitation, id string, err error) {
	return
}

func (f *FindyGrpc) Connect(ctx context.Context, invitation string) (id string, err error) {
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

func (f *FindyGrpc) SendMessage(ctx context.Context, connectionID, message string) (id string, err error) {
	return
}

func (f *FindyGrpc) ResumeCredentialOffer(ctx context.Context, id string, accept bool) (err error) {
	return
}

func (f *FindyGrpc) ResumeProofRequest(ctx context.Context, id string, accept bool) (err error) {
	return
}
