package utils

import (
	"fmt"
	"os"
	"testing"
)

func TestConfigFromEnv(t *testing.T) {
	const (
		testPort     = 1234
		testSecret   = "test-secret"
		testHost     = "test-host"
		testPath     = "test-path"
		testInsecure = "true"
	)

	strPort := fmt.Sprintf("%d", testPort)

	os.Setenv("FAV_SERVER_PORT", strPort)
	os.Setenv("FAV_JWT_KEY", testSecret)
	os.Setenv("FAV_DB_HOST", testHost)
	os.Setenv("FAV_DB_PORT", strPort)
	os.Setenv("FAV_DB_PASSWORD", testSecret)
	os.Setenv("FAV_DB_MIGRATIONS_PATH", testPath)
	os.Setenv("FAV_DB_NAME", testHost)
	os.Setenv("FAV_AGENCY_HOST", testHost)
	os.Setenv("FAV_AGENCY_PORT", strPort)
	os.Setenv("FAV_AGENCY_ADMIN_ID", testSecret)
	os.Setenv("FAV_AGENCY_CERT_PATH", testPath)
	os.Setenv("FAV_AGENCY_INSECURE", testInsecure)

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
	if config.DBMigrationsPath != testPath {
		t.Errorf("db migrations path differs")
	}
	if config.DBName != testHost {
		t.Errorf("db name differs")
	}
	if config.AgencyHost != testHost {
		t.Errorf("agency host differs")
	}
	if config.AgencyPort != testPort {
		t.Errorf("agency port differs")
	}
	if config.AgencyAdminID != testSecret {
		t.Errorf("agency admin id differs")
	}
	if config.AgencyCertPath != testPath {
		t.Errorf("agency cert path differs")
	}
	if !config.AgencyInsecure {
		t.Errorf("agency insecure differs")
	}
}
