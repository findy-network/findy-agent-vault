package agency

import (
	"github.com/findy-network/findy-agent-vault/agency/findy"
	"github.com/findy-network/findy-agent-vault/agency/mock"
	"github.com/findy-network/findy-agent-vault/agency/model"
)

const (
	AgencyTypeMock      = "MOCK"
	AgencyTypeFindyGRPC = "FINDY_GRPC"
	// TODO: is legacy needed?
)

func InitAgency(agencyType string, listener model.Listener, agents []*model.Agent) model.Agency {
	var a model.Agency
	switch agencyType {
	case AgencyTypeMock:
		a = &mock.Mock{}
	case AgencyTypeFindyGRPC:
		a = &findy.Agency{}
	}
	if a == nil {
		panic("Invalid agency type: " + agencyType)
	}

	a.Init(listener, agents)
	return a
}
