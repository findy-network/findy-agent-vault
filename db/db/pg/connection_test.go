package pg

import (
	"math"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/findy-network/findy-agent-vault/paginator"

	"github.com/findy-network/findy-agent-vault/db/fake"
	"github.com/findy-network/findy-agent-vault/db/model"
)

func validateConnection(t *testing.T, exp, got *model.Connection) {
	if got == nil {
		t.Errorf("Expecting result, connection is nil")
		return
	}
	if got.ID == "" {
		t.Errorf("Connection id invalid got %s", got.ID)
	}
	if got.TenantID != exp.TenantID {
		t.Errorf("Connection tenant id mismatch expected %s got %s", exp.TenantID, got.TenantID)
	}
	if got.OurDid != exp.OurDid {
		t.Errorf("Connection our did mismatch expected %s got %s", exp.OurDid, got.OurDid)
	}
	if got.TheirDid != exp.TheirDid {
		t.Errorf("Connection their did mismatch expected %s got %s", exp.TheirDid, got.TheirDid)
	}
	if got.TheirEndpoint != exp.TheirEndpoint {
		t.Errorf("Connection their endpoint mismatch expected %s got %s", exp.TheirEndpoint, got.TheirEndpoint)
	}
	if got.TheirLabel != exp.TheirLabel {
		t.Errorf("Connection their label mismatch expected %s got %s", exp.TheirLabel, got.TheirLabel)
	}
	if got.Invited != exp.Invited {
		t.Errorf("Connection invited mismatch expected %v got %v", exp.Invited, got.Invited)
	}
	if time.Since(got.Created) > time.Second {
		t.Errorf("Timestamp not in threshold %v", got.Created)
	}
	created := uint64(math.Round(float64(got.Created.UnixNano()) / float64(time.Millisecond.Nanoseconds())))
	if got.Cursor != created {
		t.Errorf("Cursor mismatch %v %v", got.Cursor, created)
	}
}

func TestAddConnection(t *testing.T) {
	testConnection := &model.Connection{
		TenantID:      testTenantID,
		OurDid:        "ourDid",
		TheirDid:      "theirDid",
		TheirEndpoint: "theirEndpoint",
		TheirLabel:    "theirLabel",
		Invited:       false,
	}

	// Add data
	c, err := pgDB.AddConnection(testConnection)
	if err != nil {
		t.Errorf("Failed to add connection %s", err.Error())
	} else {
		validateConnection(t, testConnection, c)
	}

	// Get data for id
	got, err := pgDB.GetConnection(c.ID, testTenantID)
	if err != nil {
		t.Errorf("Error fetching connection %s", err.Error())
	} else if !reflect.DeepEqual(&c, &got) {
		t.Errorf("Mismatch in fetched connection expected: %v  got: %v", c, got)
	}
	validateConnection(t, c, got)
}

func TestGetConnections(t *testing.T) {
	// add new agent with no pre-existing connections
	ctAgent := &model.Agent{AgentID: "connectionsTestAgentID", Label: "testAgent"}
	a, err := pgDB.AddAgent(ctAgent)
	if err != nil {
		panic(err)
	}

	size := 5
	all := fake.AddConnections(pgDB, a.ID, size)
	all = append(all, fake.AddConnections(pgDB, a.ID, size)...)
	all = append(all, fake.AddConnections(pgDB, a.ID, size)...)

	sort.Slice(all, func(i, j int) bool {
		return all[i].Created.Sub(all[j].Created) < 0
	})

	t.Run("get connections", func(t *testing.T) {
		tests := []struct {
			name   string
			args   *paginator.BatchInfo
			result *model.Connections
		}{
			{
				"first 5",
				&paginator.BatchInfo{Count: size, Tail: false},
				&model.Connections{HasNextPage: true, HasPreviousPage: false, Connections: all[:size]},
			},
			{
				"first next 5",
				&paginator.BatchInfo{Count: size, Tail: false, After: all[size-1].Cursor},
				&model.Connections{HasNextPage: true, HasPreviousPage: true, Connections: all[size : size*2]},
			},
			{
				"first last 5",
				&paginator.BatchInfo{Count: size, Tail: false, After: all[(size*2)-1].Cursor},
				&model.Connections{HasNextPage: false, HasPreviousPage: true, Connections: all[size*2:]},
			},
			{
				"last 5",
				&paginator.BatchInfo{Count: size, Tail: true},
				&model.Connections{HasNextPage: false, HasPreviousPage: true, Connections: all[size*2:]},
			},
			{
				"last next 5",
				&paginator.BatchInfo{Count: size, Tail: true, Before: all[size*2].Cursor},
				&model.Connections{HasNextPage: true, HasPreviousPage: true, Connections: all[size : size*2]},
			},
			{
				"last first 5",
				&paginator.BatchInfo{Count: size, Tail: true, Before: all[size].Cursor},
				&model.Connections{HasNextPage: true, HasPreviousPage: false, Connections: all[:size]},
			},
			{
				"all",
				&paginator.BatchInfo{Count: size * 3, Tail: false},
				&model.Connections{HasNextPage: false, HasPreviousPage: false, Connections: all},
			},
		}

		for _, testCase := range tests {
			tc := testCase
			t.Run(tc.name, func(t *testing.T) {
				c, err := pgDB.GetConnections(tc.args, a.ID)
				if err != nil {
					t.Errorf("Error fetching connections %s", err.Error())
				} else {
					if len(c.Connections) != tc.args.Count {
						t.Errorf("Mismatch in connection count: %v  got: %v", len(c.Connections), tc.args.Count)
					}
					if c.HasNextPage != tc.result.HasNextPage {
						t.Errorf("Batch next page mismatch %v got: %v", c.HasNextPage, tc.result.HasNextPage)
					}
					if c.HasPreviousPage != tc.result.HasPreviousPage {
						t.Errorf("Batch previous page mismatch %v got: %v", c.HasPreviousPage, tc.result.HasPreviousPage)
					}
					for index, connection := range c.Connections {
						validateConnection(t, tc.result.Connections[index], connection)
					}
				}
			})
		}
	})
}
