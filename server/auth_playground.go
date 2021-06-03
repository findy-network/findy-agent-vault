package server

import (
	"time"

	"github.com/findy-network/findy-common-go/jwt"
)

const (
	hoursForTest = 2
)

func (v *VaultServer) CreateTestToken(userName, validationKey string) string {
	jwt.SetJWTSecret(validationKey) // for test token generation
	return jwt.BuildJWTWithTime(userName, "minnie mouse", time.Hour*hoursForTest)
}
