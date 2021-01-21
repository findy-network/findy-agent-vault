package agency

import (
	"github.com/findy-network/findy-agent-vault/agency/findy"
	"github.com/findy-network/findy-agent-vault/agency/mock"
	"github.com/findy-network/findy-agent-vault/agency/model"
)

const (
	AgencyTypeMock      = "MOCK"
	AgencyTypeFindyGRPC = "FINDY_GRPC"
)

func Create(agencyType string) model.Agency {
	register := make(map[string]model.Agency)

	register[AgencyTypeFindyGRPC] = &findy.Agency{}
	register[AgencyTypeMock] = &mock.Mock{}

	a := register[agencyType]

	if a == nil {
		panic("Invalid agency type: " + agencyType)
	}

	return a
}
