package utils

import (
	"fmt"
	"testing"

	"github.com/lainio/err2/assert"
)

func TestConfigFromEnv(t *testing.T) {
	assert.PushTester(t)
	defer assert.PopTester()

	const (
		testPort     = 1234
		testSecret   = "test-secret"
		testHost     = "test-host"
		testPath     = "test-path"
		testInsecure = "true"
	)

	strPort := fmt.Sprintf("%d", testPort)

	t.Setenv("FAV_SERVER_PORT", strPort)
	t.Setenv("FAV_JWT_KEY", testSecret)
	t.Setenv("FAV_DB_HOST", testHost)
	t.Setenv("FAV_DB_PORT", strPort)
	t.Setenv("FAV_DB_PASSWORD", testSecret)
	t.Setenv("FAV_DB_MIGRATIONS_PATH", testPath)
	t.Setenv("FAV_DB_NAME", testHost)
	t.Setenv("FAV_AGENCY_HOST", testHost)
	t.Setenv("FAV_AGENCY_PORT", strPort)
	t.Setenv("FAV_AGENCY_ADMIN_ID", testSecret)
	t.Setenv("FAV_AGENCY_CERT_PATH", testPath)
	t.Setenv("FAV_AGENCY_INSECURE", testInsecure)

	config := LoadConfig()
	assert.Equal(config.ServerPort, testPort, "config port differs")
	assert.Equal(config.JWTKey, testSecret, "config jwt key differs")
	assert.Equal(config.Address, fmt.Sprintf(":%d", testPort), "config address differs")
	assert.Equal(config.DBHost, testHost, "db host differs")
	assert.Equal(config.DBPort, testPort, "db port differs")
	assert.Equal(config.DBPassword, testSecret, "db password differs")
	assert.Equal(config.DBMigrationsPath, testPath, "db migrations path differs")
	assert.Equal(config.DBName, testHost, "db name differs")
	assert.Equal(config.AgencyHost, testHost, "agency host differs")
	assert.Equal(config.AgencyPort, testPort, "agency port differs")
	assert.Equal(config.AgencyAdminID, testSecret, "agency admin id differs")
	assert.Equal(config.AgencyCertPath, testPath, "agency cert path differs")
	assert.Equal(config.AgencyCertPath, testPath, "agency cert path differs")
	assert.That(config.AgencyInsecure, "agency insecure differs")
}
