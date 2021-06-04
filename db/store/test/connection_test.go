package test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/findy-network/findy-agent-vault/db/fake"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/lainio/err2/assert"

	"github.com/findy-network/findy-agent-vault/paginator"

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
	validateCreatedTS(t, got.Cursor, &got.Created)
	validateTimestap(t, &exp.Approved, &got.Approved, "Approved")
	validateTimestap(t, &exp.Archived, &got.Archived, "Archived")
}

func TestAddConnection(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("add connection "+s.name, func(t *testing.T) {
			s.updateTestConnection()
			// Add data
			c, err := s.db.AddConnection(s.testConnection)
			if err != nil {
				t.Errorf("Failed to add connection %s", err.Error())
			} else {
				validateConnection(t, s.testConnection, c)
			}

			// Get data for id
			got, err := s.db.GetConnection(c.ID, s.testTenantID)
			if err != nil {
				t.Errorf("Error fetching connection %s", err.Error())
			} else if !reflect.DeepEqual(&c, &got) {
				t.Errorf("Mismatch in fetched connection expected: %+v  got: %+v", c, got)
			}
		})
	}
}

func TestAddConnectionSameIDDifferentTenant(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("add connection with same id "+s.name, func(t *testing.T) {
			s.updateTestConnection()
			// Add data
			connection1, err := s.db.AddConnection(s.testConnection)
			if err != nil {
				t.Errorf("Failed to add connection %s", err.Error())
			} else {
				validateConnection(t, s.testConnection, connection1)
			}

			// Add connection with same id
			agent2 := fake.AddAgent(s.db)
			s.testConnection.TenantID = agent2.TenantID
			connection2, err := s.db.AddConnection(s.testConnection)
			if err != nil {
				t.Errorf("Failed to add connection with same id %s", err.Error())
			} else {
				validateConnection(t, s.testConnection, connection2)
			}
		})
	}
}

func TestGetConnections(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("get connections "+s.name, func(t *testing.T) {
			size := 5
			a, all := AddAgentAndConnections(s.db, "TestGetConnections", size*3)

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
						c, err := s.db.GetConnections(tc.args, a.ID)
						if err != nil {
							t.Errorf("Error fetching connections %s", err.Error())
						} else {
							if len(c.Connections) != tc.args.Count {
								t.Errorf("Mismatch in connection count: %v  expected: %v", len(c.Connections), tc.args.Count)
							}
							if c.HasNextPage != tc.result.HasNextPage {
								t.Errorf("Batch next page mismatch %v expected: %v", c.HasNextPage, tc.result.HasNextPage)
							}
							if c.HasPreviousPage != tc.result.HasPreviousPage {
								t.Errorf("Batch previous page mismatch %v expected: %v", c.HasPreviousPage, tc.result.HasPreviousPage)
							}
							for index, connection := range c.Connections {
								validateConnection(t, tc.result.Connections[index], connection)
							}
						}
					})
				}
			})
		})
	}
}

func TestGetConnectionCount(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("get connection count "+s.name, func(t *testing.T) {
			size := 5
			a, _ := AddAgentAndConnections(s.db, "TestGetConnectionCount", size)

			// Get count
			got, err := s.db.GetConnectionCount(a.ID)
			if err != nil {
				t.Errorf("Error fetching connection %s", err.Error())
			} else if got != size {
				t.Errorf("Mismatch in fetched connection count expected: %v  got: %v", size, got)
			}
		})
	}
}

func TestArchiveConnection(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("archive connection "+s.name, func(t *testing.T) {
			s.updateTestConnection()
			// Add data
			c, err := s.db.AddConnection(s.testConnection)
			assert.D.True(err == nil)

			// Archive object
			now := utils.CurrentTime()
			err = s.db.ArchiveConnection(c.ID, c.TenantID)
			if err != nil {
				t.Errorf("Failed to archive connection %s", err.Error())
			}

			// Get data for id
			got, err := s.db.GetConnection(c.ID, c.TenantID)
			assert.D.True(err == nil)

			c.Archived = now
			validateConnection(t, c, got)
		})
	}
}
