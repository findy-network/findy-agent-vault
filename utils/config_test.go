package utils

import (
	"fmt"
	"os"
	"testing"
)

func TestConfigFromEnv(t *testing.T) {
	const (
		testPort   = 1234
		testSecret = "test-secret"
	)
	os.Setenv("FAV_SERVER_PORT", fmt.Sprintf("%d", testPort))
	os.Setenv("FAV_JWT_KEY", testSecret)

	config := LoadConfig()
	if config.ServerPort != testPort {
		t.Errorf("config port differs")
	}
	if config.JWTKey != testSecret {
		t.Errorf("config jwt key differs")
	}
	if config.Address != fmt.Sprintf(":%d", testPort) {
		t.Errorf("config address differs")
	}
}
