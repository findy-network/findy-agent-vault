package pg

import (
	"reflect"
	"testing"
	"time"

	"github.com/findy-network/findy-agent-vault/db/model"
)

func TestAddConnection(t *testing.T) {
	testConnection := &model.Connection{
		TenantID:      testTenantID,
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
		if c.Cursor != ceilTimestamp(&c.Created) {
			t.Errorf("Cursor mismatch %v %v", c.Cursor, ceilTimestamp(&c.Created))
		}
	}

	// Add data
	c, err := pgDB.AddConnection(testConnection)
	if err != nil {
		t.Errorf("Failed to add connection %s", err.Error())
	} else {
		validateConnection(c)
	}

	// Get data for id
	got, err := pgDB.GetConnection(c.ID, testAgentID)
	if err != nil {
		t.Errorf("Error fetching connection %s", err.Error())
	} else if !reflect.DeepEqual(&c, &got) {
		t.Errorf("Mismatch in fetched connection expected: %v  got: %v", c, got)
	}
}
