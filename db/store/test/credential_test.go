package test

import (
	"sort"
	"testing"
	"time"

	"github.com/findy-network/findy-agent-vault/db/fake"
	"github.com/findy-network/findy-agent-vault/db/model"
	graph "github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/paginator"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/lainio/err2/assert"
)

func validateAttributes(t *testing.T, exp, got []*graph.CredentialValue) {
	if len(got) == 0 {
		t.Errorf("No attributes found")
	}
	for index, a := range got {
		if a.ID == "" {
			t.Errorf("Credential attribute id invalid.")
		}
		if a.Name != exp[index].Name {
			t.Errorf("Credential attribute name mismatch: expected %s got %s.", exp[index].Name, a.Name)
		}
		if a.Value != exp[index].Value {
			t.Errorf("Credential attribute value mismatch: expected %s got %s.", exp[index].Value, a.Value)
		}
	}
}

func validateCredential(t *testing.T, exp, got *model.Credential) {
	if got == nil {
		t.Errorf("Expecting result, credential is nil")
		return
	}
	if got.ID == "" {
		t.Errorf("Credential id invalid.")
	}
	if got.TenantID != exp.TenantID {
		t.Errorf("Credential tenant id mismatch expected %s got %s", exp.TenantID, got.TenantID)
	}
	if got.ConnectionID != exp.ConnectionID {
		t.Errorf("Credential connection id mismatch expected %s got %s", exp.ConnectionID, got.ConnectionID)
	}
	if got.Role != exp.Role {
		t.Errorf("Credential Role mismatch expected %s got %s", exp.Role, got.Role)
	}
	if got.SchemaID != exp.SchemaID {
		t.Errorf("Credential SchemaID mismatch expected %s got %s", exp.SchemaID, got.SchemaID)
	}
	if got.CredDefID != exp.CredDefID {
		t.Errorf("Credential CredDefID mismatch expected %s got %s", exp.CredDefID, got.CredDefID)
	}
	if got.InitiatedByUs != exp.InitiatedByUs {
		t.Errorf("Credential InitiatedByUs mismatch expected %v got %v", exp.InitiatedByUs, got.InitiatedByUs)
	}
	validateTimestap(t, &exp.Approved, &got.Approved, "Approved")
	validateTimestap(t, &exp.Issued, &got.Issued, "Issued")
	validateTimestap(t, &exp.Failed, &got.Failed, "Failed")
	validateTimestap(t, &exp.Archived, &got.Archived, "Archived")
	validateCreatedTS(t, got.Cursor, &got.Created)
	validateAttributes(t, exp.Attributes, got.Attributes)
}

func validateCredentials(t *testing.T, expCount int, exp, got *model.Credentials) {
	if len(got.Credentials) != expCount {
		t.Errorf("Mismatch in credential count: %v  got: %v", len(got.Credentials), expCount)
	}
	if got.HasNextPage != exp.HasNextPage {
		t.Errorf("Batch next page mismatch %v got: %v", got.HasNextPage, exp.HasNextPage)
	}
	if got.HasPreviousPage != exp.HasPreviousPage {
		t.Errorf("Batch previous page mismatch %v got: %v", got.HasPreviousPage, exp.HasPreviousPage)
	}
	for index, credential := range got.Credentials {
		validateCredential(t, exp.Credentials[index], credential)
	}
}

type credTest struct {
	name   string
	args   *paginator.BatchInfo
	result *model.Credentials
}

func getCredTests(size int, all []*model.Credential) []*credTest {
	var credTests = []*credTest{
		{
			"first 5",
			&paginator.BatchInfo{Count: size, Tail: false},
			&model.Credentials{HasNextPage: true, HasPreviousPage: false, Credentials: all[:size]},
		},
		{
			"first next 5",
			&paginator.BatchInfo{Count: size, Tail: false, After: all[size-1].Cursor},
			&model.Credentials{HasNextPage: true, HasPreviousPage: true, Credentials: all[size : size*2]},
		},
		{
			"first last 5",
			&paginator.BatchInfo{Count: size, Tail: false, After: all[(size*2)-1].Cursor},
			&model.Credentials{HasNextPage: false, HasPreviousPage: true, Credentials: all[size*2:]},
		},
		{
			"last 5",
			&paginator.BatchInfo{Count: size, Tail: true},
			&model.Credentials{HasNextPage: false, HasPreviousPage: true, Credentials: all[size*2:]},
		},
		{
			"last next 5",
			&paginator.BatchInfo{Count: size, Tail: true, Before: all[size*2].Cursor},
			&model.Credentials{HasNextPage: true, HasPreviousPage: true, Credentials: all[size : size*2]},
		},
		{
			"last first 5",
			&paginator.BatchInfo{Count: size, Tail: true, Before: all[size].Cursor},
			&model.Credentials{HasNextPage: true, HasPreviousPage: false, Credentials: all[:size]},
		},
		{
			"all",
			&paginator.BatchInfo{Count: size * 3, Tail: false},
			&model.Credentials{HasNextPage: false, HasPreviousPage: false, Credentials: all},
		},
	}
	return credTests
}

func TestAddCredential(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("add credential "+s.name, func(t *testing.T) {
			testCredential = s.newTestCredential(testCredential)

			// Add data
			c, err := s.db.AddCredential(testCredential)
			if err != nil {
				t.Errorf("Failed to add credential %s", err.Error())
			} else {
				validateCredential(t, testCredential, c)
			}

			// Get data for id
			got, err := s.db.GetCredential(c.ID, s.testTenantID)
			if err != nil {
				t.Errorf("Error fetching credential %s", err.Error())
			} else {
				validateCredential(t, c, got)
			}
		})
	}
}

func TestUpdateCredential(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("update credential "+s.name, func(t *testing.T) {
			testCredential.TenantID = s.testTenantID
			testCredential.ConnectionID = s.testConnectionID

			// Add data
			c, err := s.db.AddCredential(testCredential)
			if err != nil {
				t.Errorf("Failed to add credential %s", err.Error())
			}

			// Update data
			now := time.Now().UTC()
			c.Approved = now
			c.Issued = now
			_, err = s.db.UpdateCredential(c)
			if err != nil {
				t.Errorf("Failed to update credential %s", err.Error())
			}

			// Get data for id
			got, err := s.db.GetCredential(c.ID, s.testTenantID)
			if err != nil {
				t.Errorf("Error fetching credential %s", err.Error())
			} else {
				validateCredential(t, c, got)
			}
		})
	}
}

func TestGetTenantCredentials(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("get credentials "+s.name, func(t *testing.T) {
			a, connections := AddAgentAndConnections(s.db, "TestGetTenantCredentials", 3)

			size := 5
			all := fake.AddCredentials(s.db, a.ID, connections[0].ID, size)
			all = append(all, fake.AddCredentials(s.db, a.ID, connections[1].ID, size)...)
			all = append(all, fake.AddCredentials(s.db, a.ID, connections[2].ID, size)...)

			sort.Slice(all, func(i, j int) bool {
				return all[i].Created.Sub(all[j].Created) < 0
			})

			t.Run("get credentials", func(t *testing.T) {
				tests := getCredTests(size, all)

				for _, testCase := range tests {
					tc := testCase
					t.Run(tc.name, func(t *testing.T) {
						c, err := s.db.GetCredentials(tc.args, a.ID, nil)
						if err != nil {
							t.Errorf("Error fetching credentials %s", err.Error())
						} else {
							validateCredentials(t, tc.args.Count, c, tc.result)
						}
					})
				}
			})
		})
	}
}

func TestGetConnectionCredentials(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("get connection credentials "+s.name, func(t *testing.T) {
			// add new agent with no pre-existing credentials
			a, connections := AddAgentAndConnections(s.db, "TestGetConnectionCredentials", 3)

			size := 5
			countPerConnection := size * 3
			fake.AddCredentials(s.db, a.ID, connections[0].ID, countPerConnection)
			fake.AddCredentials(s.db, a.ID, connections[1].ID, countPerConnection)
			all := fake.AddCredentials(s.db, a.ID, connections[2].ID, countPerConnection)

			sort.Slice(all, func(i, j int) bool {
				return all[i].Created.Sub(all[j].Created) < 0
			})

			t.Run("get credentials", func(t *testing.T) {
				tests := getCredTests(size, all)

				for _, testCase := range tests {
					tc := testCase
					t.Run(tc.name, func(t *testing.T) {
						c, err := s.db.GetCredentials(tc.args, a.ID, &connections[2].ID)
						if err != nil {
							t.Errorf("Error fetching connection credentials %s", err.Error())
						} else {
							validateCredentials(t, tc.args.Count, c, tc.result)
						}
					})
				}
			})
		})
	}
}

func TestGetCredentialCount(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("get credentials count "+s.name, func(t *testing.T) {
			// add new agent with no pre-existing credentials
			a, connections := AddAgentAndConnections(s.db, "TestGetCredentialCount", 3)
			size := 5
			fake.AddCredentials(s.db, a.ID, connections[0].ID, size)

			// Get count
			got, err := s.db.GetCredentialCount(a.ID, nil)
			if err != nil {
				t.Errorf("Error fetching count %s", err.Error())
			} else if got != size {
				t.Errorf("Mismatch in fetched credential count expected: %v  got: %v", size, got)
			}
		})
	}
}

func TestGetConnectionCredentialCount(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("get connection credentials count "+s.name, func(t *testing.T) {
			// add new agent with no pre-existing credentials
			a, connections := AddAgentAndConnections(s.db, "TestGetConnectionCredentialCount", 3)
			size := 5
			index := 0
			fake.AddCredentials(s.db, a.ID, connections[index].ID, (index+1)*size)
			index++
			fake.AddCredentials(s.db, a.ID, connections[index].ID, (index+1)*size)
			index++
			fake.AddCredentials(s.db, a.ID, connections[index].ID, index*size)

			// Get count
			expected := index * size
			got, err := s.db.GetCredentialCount(a.ID, &connections[index].ID)
			if err != nil {
				t.Errorf("Error fetching count %s", err.Error())
			} else if got != expected {
				t.Errorf("Mismatch in fetched credential count expected: %v  got: %v", expected, got)
			}
		})
	}
}

func TestGetConnectionForCredential(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("get connection for credential"+s.name, func(t *testing.T) {
			a, connections := AddAgentAndConnections(s.db, "TestGetConnectionForCredential", 3)
			connection := connections[0]
			credentials := fake.AddCredentials(s.db, a.ID, connection.ID, 1)
			credential := credentials[0]

			// Get data for id
			got, err := s.db.GetConnectionForCredential(credential.ID, a.ID)
			if err != nil {
				t.Errorf("Error fetching connection %s", err.Error())
			} else {
				validateConnection(t, connection, got)
			}
		})
	}
}

func TestArchiveCredential(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("archive credential "+s.name, func(t *testing.T) {
			testCredential = s.newTestCredential(testCredential)

			// Add data
			c, err := s.db.AddCredential(testCredential)
			assert.That(err == nil)

			now := utils.CurrentTime()
			err = s.db.ArchiveCredential(c.ID, c.TenantID)
			if err != nil {
				t.Errorf("Failed to archive credential %s", err.Error())
			}

			// Get data for id
			got, err := s.db.GetCredential(c.ID, c.TenantID)
			assert.That(err == nil)

			c.Archived = now
			validateCredential(t, c, got)
		})
	}
}

func TestSearchCredentials(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("search credentials "+s.name, func(t *testing.T) {
			now := utils.CurrentTime()
			testCredential = s.newTestCredential(testCredential)

			// Add data
			testCredential.CredDefID = "searchCredentials"
			c, err := s.db.AddCredential(testCredential)
			assert.That(err == nil)

			c.Issued = now
			c, err = s.db.UpdateCredential(c)
			assert.That(err == nil)

			proofRequest := &graph.Proof{Attributes: testProof.Attributes}
			proofRequest.Attributes[0].Name = testCredential.Attributes[0].Name
			proofRequest.Attributes[0].CredDefID = testCredential.CredDefID
			res, err := s.db.SearchCredentials(s.testTenantID, proofRequest.Attributes)

			if err != nil {
				t.Errorf("Encountered error when searching for creds %s", err)
			}
			if len(res) != len(proofRequest.Attributes) {
				t.Errorf("Unexpected credential count %d", len(res))
			}
			if len(res[0].Credentials) != 1 {
				t.Errorf("Expected credential was not found")
			}
			if res[0].Credentials[0].CredentialID != c.ID {
				t.Errorf("Found Unexpected credential id, got: %s, exp: %s", res[0].Credentials[0].CredentialID, c.ID)
			}
		})
	}
}

func TestSearchCredentialsWithApostrophe(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("search credentials with apostrophe "+s.name, func(t *testing.T) {
			now := utils.CurrentTime()
			testCredential = s.newTestCredential(testCredential)

			// Add data
			testCredential.CredDefID = "searchCredentials's"
			c, err := s.db.AddCredential(testCredential)
			assert.That(err == nil)

			c.Issued = now
			c, err = s.db.UpdateCredential(c)
			assert.That(err == nil)

			proofRequest := &graph.Proof{Attributes: testProof.Attributes}
			proofRequest.Attributes[0].CredDefID = testCredential.CredDefID
			res, err := s.db.SearchCredentials(s.testTenantID, proofRequest.Attributes)

			if err != nil {
				t.Errorf("Encountered error when searching for creds %s", err)
			}
			if len(res) != len(proofRequest.Attributes) {
				t.Errorf("Unexpected credential count %d", len(res))
			}
			if len(res[0].Credentials) != 1 {
				t.Errorf("Expected credential was not found")
			}
			if res[0].Credentials[0].CredentialID != c.ID {
				t.Errorf("Found Unexpected credential id, got: %s, exp: %s", res[0].Credentials[0].CredentialID, c.ID)
			}
		})
	}
}
