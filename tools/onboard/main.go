package main

import (
	"fmt"

	clock "github.com/findy-network/findy-agent-vault/utils"
	"github.com/findy-network/findy-grpc/enclave"
	"github.com/findy-network/findy-grpc/utils"
)

func main() {
	utils.ParseLoggingArgs("")

	enclave.Init("localhost", 50052)

	u := &enclave.User{Name: fmt.Sprintf("%d-minnie@example.com", clock.CurrentTimeMs())}
	//u := &enclave.User{Name: "minnie@example.com"}
	if err := u.AllocateCloudAgent(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("DID:", u.DID)
		fmt.Println("JWT:", u.JWT())
	}

}
