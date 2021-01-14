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
		testHost   = "test-host"
	)

	strPort := fmt.Sprintf("%d", testPort)

	os.Setenv("FAV_SERVER_PORT", strPort)
	os.Setenv("FAV_JWT_KEY", testSecret)
	os.Setenv("FAV_DB_HOST", testHost)
	os.Setenv("FAV_DB_PORT", strPort)
	os.Setenv("FAV_DB_PASSWORD", testSecret)

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
	if config.DBHost != testHost {
		t.Errorf("db host differs")
	}
	if config.DBPort != testPort {
		t.Errorf("db port differs")
	}
	if config.DBPassword != testSecret {
		t.Errorf("db password differs")
	}
}
