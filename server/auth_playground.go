package server

import (
	"time"

	"github.com/findy-network/findy-agent-vault/db/fake"
	"github.com/findy-network/findy-grpc/jwt"
)

const (
	hoursInWeek  = 24 * 7
	hoursForTest = 2
)

func (v *VaultServer) CreateToken(validationKey string) string {
	jwt.SetJWTSecret(validationKey) // for test token generation
	return jwt.BuildJWTWithTime(fake.FakeCloudDID, "minnie mouse", time.Hour*hoursInWeek)
}

func (v *VaultServer) CreateTestToken(validationKey string) string {
	jwt.SetJWTSecret(validationKey) // for test token generation
	return jwt.BuildJWTWithTime(fake.FakeCloudDID, "minnie mouse", time.Hour*hoursForTest)
}
