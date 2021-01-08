package server

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/findy-network/findy-agent-vault/db/fake"
)

const (
	hoursInWeek  = 24 * 7
	hoursForTest = 2
)

func (v *VaultServer) CreateToken(id string) (string, error) {
	return v.createTokenString(id, time.Hour*hoursInWeek)
}

func (v *VaultServer) CreateTestToken(id string) *jwt.Token {
	token := createClaims(id, time.Hour*hoursForTest)
	str, _ := token.SignedString(v.authChecker.secret)
	token.Raw = str
	return token
}

func (v *VaultServer) createTokenString(id string, duration time.Duration) (string, error) {
	signer := createClaims(id, duration)
	return signer.SignedString(v.authChecker.secret)
}

func createClaims(id string, duration time.Duration) *jwt.Token {
	claims := jwt.MapClaims{}
	claims["id"] = id
	claims["un"] = fake.FakeCloudDID
	claims["label"] = "minnie mouse"
	claims["exp"] = time.Now().Add(duration).Unix()
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
}

func (v *VaultServer) createTestToken() string {
	token, _ := v.createTokenString("test", time.Hour*hoursForTest)
	return token
}
