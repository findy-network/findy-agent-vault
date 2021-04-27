package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	didexchange "github.com/findy-network/findy-common-go/std/didexchange/invitation"
	"github.com/findy-network/findy-common-go/jwt"
)

func main() {
	flag.Parse()
	var invitation didexchange.Invitation
	_ = json.Unmarshal([]byte(os.Args[1]), &invitation)

	fmt.Println(jwt.BuildJWT(strings.Split(invitation.ServiceEndpoint, "/")[4]))
}
