package pg

import (
	"testing"
	"time"

	"github.com/findy-network/findy-agent-vault/db/model"
)

func TestAddConnection(t *testing.T) {
	testConnection := &model.Connection{
		TenantID:      testTenantId,
		OurDid:        "ourDid",
		TheirDid:      "theirDid",
		TheirEndpoint: "theirEndpoint",
		TheirLabel:    "theirLabel",
		Invited:       false,
	}

	var validateConnection = func(c *model.Connection) {
		if c == nil {
			t.Errorf("Expecting result, connection is nil")
			return
		}
		if c.ID == "" {
			t.Errorf("Connection id invalid got %s", c.ID)
		}
		if c.TenantID != testConnection.TenantID {
			t.Errorf("Connection tenant id mismatch expected %s got %s", testConnection.TenantID, c.TenantID)
		}
		if c.OurDid != testConnection.OurDid {
			t.Errorf("Connection our did mismatch expected %s got %s", testConnection.OurDid, c.OurDid)
		}
		if c.TheirDid != testConnection.TheirDid {
			t.Errorf("Connection their did mismatch expected %s got %s", testConnection.TheirDid, c.TheirDid)
		}
		if c.TheirEndpoint != testConnection.TheirEndpoint {
			t.Errorf("Connection their endpoint mismatch expected %s got %s", testConnection.TheirEndpoint, c.TheirEndpoint)
		}
		if c.TheirLabel != testConnection.TheirLabel {
			t.Errorf("Connection their label mismatch expected %s got %s", testConnection.TheirLabel, c.TheirLabel)
		}
		if c.Invited != testConnection.Invited {
			t.Errorf("Connection invited mismatch expected %v got %v", testConnection.Invited, c.Invited)
		}
		if time.Since(c.Created) > time.Second {
			t.Errorf("Timestamp not in threshold %v", c.Created)
		}
		if c.Cursor != uint64(c.Created.Unix()) {
			t.Errorf("Cursor mismatch %v %v", c.Cursor, c.Created.Unix())
		}
	}

	// Add data
	c, err := pgDB.AddConnection(testConnection)
	if err != nil {
		t.Errorf("Failed to add connection %s", err.Error())
	} else {
		validateConnection(c)
	}

	// Error with duplicate id
	/*err := pgDB.AddAgent(testConnection)
	if err == nil {
		t.Errorf("Expecting duplicate key error")
	}

	if pgErr, ok := err.(*PgError); ok {
		if pgErr.code != PgErrorUniqueViolation {
			t.Errorf("Expecting duplicate key error %s", pgErr.code)
		}
	} else {
		t.Errorf("Expecting pg error %v", err)
	}*/

	/*var validateConnection = func(c *model.Connection) {
		if c == nil {
			t.Errorf("Expecting result, connection is nil")
			return
		}
		if c.TenantID != testConnection.TenantID {
			t.Errorf("Connection tenant id mismatch expected %s got %s", testConnection.TenantID, c.TenantID)
		}
		if c.OurDid != testConnection.OurDid {
			t.Errorf("Connection our did mismatch expected %s got %s", testConnection.OurDid, c.OurDid)
		}
		if c.TheirDid != testConnection.TheirDid {
			t.Errorf("Connection their did mismatch expected %s got %s", testConnection.TheirDid, c.TheirDid)
		}
		if c.TheirEndpoint != testConnection.TheirEndpoint {
			t.Errorf("Connection their endpoint mismatch expected %s got %s", testConnection.TheirEndpoint, c.TheirEndpoint)
		}
		if c.TheirLabel != testConnection.TheirLabel {
			t.Errorf("Connection their label mismatch expected %s got %s", testConnection.TheirLabel, c.TheirLabel)
		}
		if c.Invited != testConnection.Invited {
			t.Errorf("Connection invited mismatch expected %s got %s", testConnection.Invited, c.Invited)
		}
		if time.Since(c.Created) > time.Second {
			t.Errorf("Timestamp not in threshold %v", c.Created)
		}
	}

	// Get data for id
	c, err := pgDB.GetAgent(&t.ID, nil)
	if err != nil {
		t.Errorf("Error fetching agent %s", err.Error())
	} else {
		validateAgent(a2)
	}*/
}
