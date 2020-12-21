package agency

import (
	"github.com/findy-network/findy-agent-vault/agency/model"
)

const (
	AgencyTypeMock      = "MOCK"
	AgencyTypeFindyGRPC = "FINDY_GRPC"
	// TODO: is legacy needed?
)

var (
	Register map[string]model.Agency = make(map[string]model.Agency)
)

func InitAgency(agencyType string, listener model.Listener, agents []*model.Agent) model.Agency {
	a := Register[agencyType]

	if a == nil {
		panic("Invalid agency type: " + agencyType)
	}

	a.Init(listener, agents)
	return a
}
