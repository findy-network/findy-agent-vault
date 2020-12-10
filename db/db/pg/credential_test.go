package pg

import (
	"math"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/findy-network/findy-agent-vault/db/fake"
	"github.com/findy-network/findy-agent-vault/db/model"
	graph "github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/paginator"
)

func validateTimestap(t *testing.T, exp, got *time.Time, name string) {
	fail := false
	if got != exp {
		fail = true
		if got != nil && exp != nil && got.Sub(*exp) == 0 {
			fail = false
		}
	}
	if fail {
		t.Errorf("Credential %s mismatch expected %s got %s", name, exp, got)
	}
}

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
	if time.Since(got.Created) > time.Second {
		t.Errorf("Timestamp not in threshold %v", got.Created)
	}
	validateTimestap(t, exp.Approved, got.Approved, "Approved")
	validateTimestap(t, exp.Issued, got.Issued, "Issued")
	validateTimestap(t, exp.Failed, got.Failed, "Failed")
	created := uint64(math.Round(float64(got.Created.UnixNano()) / float64(time.Millisecond.Nanoseconds())))
	if got.Cursor != created {
		t.Errorf("Cursor mismatch %v %v", got.Cursor, created)
	}
	validateAttributes(t, exp.Attributes, got.Attributes)
}

func TestAddCredential(t *testing.T) {
	testCredential.TenantID = testTenantID
	testCredential.ConnectionID = testConnectionID

	// Add data
	c, err := pgDB.AddCredential(testCredential)
	if err != nil {
		t.Errorf("Failed to add credential %s", err.Error())
	} else {
		validateCredential(t, testCredential, c)
	}

	// Get data for id
	got, err := pgDB.GetCredential(c.ID, testTenantID)
	if err != nil {
		t.Errorf("Error fetching credential %s", err.Error())
	} else if !reflect.DeepEqual(&c, &got) {
		t.Errorf("Mismatch in fetched credential expected: %v  got: %v", c, got)
	}
	validateCredential(t, c, got)
}

func TestUpdateCredential(t *testing.T) {
	testCredential.TenantID = testTenantID
	testCredential.ConnectionID = testConnectionID

	// Add data
	c, err := pgDB.AddCredential(testCredential)
	if err != nil {
		t.Errorf("Failed to add credential %s", err.Error())
	}

	// Update data
	now := time.Now().UTC()
	c.Approved = &now
	c.Issued = &now
	_, err = pgDB.UpdateCredential(c)
	if err != nil {
		t.Errorf("Failed to update credential %s", err.Error())
	}

	// Get data for id
	got, err := pgDB.GetCredential(c.ID, testTenantID)
	if err != nil {
		t.Errorf("Error fetching credential %s", err.Error())
	} else if !reflect.DeepEqual(&c, &got) {
		t.Errorf("Mismatch in fetched credential expected: %v  got: %v", c, got)
	}
	validateCredential(t, c, got)
}

func addAgentAndConnections(agentID string) (*model.Agent, []*model.Connection) {
	// add new agent with no pre-existing credentials
	ctAgent := &model.Agent{AgentID: agentID, Label: "testAgent"}
	a, err := pgDB.AddAgent(ctAgent)
	if err != nil {
		panic(err)
	}
	// add new connections
	connCount := 3
	connections := make([]*model.Connection, connCount)
	for i := 0; i < connCount; i++ {
		c, err := pgDB.AddConnection(testConnection)
		if err != nil {
			panic(err)
		}
		connections[i] = c
	}
	return a, connections
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

func TestGetTenantCredentials(t *testing.T) {
	// add new agent with no pre-existing credentials
	a, connections := addAgentAndConnections("TestGetTenantCredentials")

	size := 5
	all := fake.AddCredentials(pgDB, a.ID, connections[0].ID, size)
	all = append(all, fake.AddCredentials(pgDB, a.ID, connections[1].ID, size)...)
	all = append(all, fake.AddCredentials(pgDB, a.ID, connections[2].ID, size)...)

	sort.Slice(all, func(i, j int) bool {
		return all[i].Created.Sub(all[j].Created) < 0
	})

	t.Run("get credentials", func(t *testing.T) {
		tests := getCredTests(size, all)

		for _, testCase := range tests {
			tc := testCase
			t.Run(tc.name, func(t *testing.T) {
				c, err := pgDB.GetCredentials(tc.args, a.ID)
				if err != nil {
					t.Errorf("Error fetching credentials %s", err.Error())
				} else {
					if len(c.Credentials) != tc.args.Count {
						t.Errorf("Mismatch in credential count: %v  got: %v", len(c.Credentials), tc.args.Count)
					}
					if c.HasNextPage != tc.result.HasNextPage {
						t.Errorf("Batch next page mismatch %v got: %v", c.HasNextPage, tc.result.HasNextPage)
					}
					if c.HasPreviousPage != tc.result.HasPreviousPage {
						t.Errorf("Batch previous page mismatch %v got: %v", c.HasPreviousPage, tc.result.HasPreviousPage)
					}
					for index, credential := range c.Credentials {
						validateCredential(t, tc.result.Credentials[index], credential)
					}
				}
			})
		}
	})
}

func TestGetConnectionCredentials(t *testing.T) {
	// add new agent with no pre-existing credentials
	a, connections := addAgentAndConnections("TestGetConnectionCredentials")

	size := 5
	countPerConnection := size * 3
	fake.AddCredentials(pgDB, a.ID, connections[0].ID, countPerConnection)
	fake.AddCredentials(pgDB, a.ID, connections[1].ID, countPerConnection)
	all := fake.AddCredentials(pgDB, a.ID, connections[2].ID, countPerConnection)

	sort.Slice(all, func(i, j int) bool {
		return all[i].Created.Sub(all[j].Created) < 0
	})

	t.Run("get credentials", func(t *testing.T) {
		tests := getCredTests(size, all)

		for _, testCase := range tests {
			tc := testCase
			t.Run(tc.name, func(t *testing.T) {
				c, err := pgDB.GetConnectionCredentials(tc.args, a.ID, connections[2].ID)
				if err != nil {
					t.Errorf("Error fetching credentials %s", err.Error())
				} else {
					if len(c.Credentials) != tc.args.Count {
						t.Errorf("Mismatch in credential count: %v  got: %v", len(c.Credentials), tc.args.Count)
					}
					if c.HasNextPage != tc.result.HasNextPage {
						t.Errorf("Batch next page mismatch %v got: %v", c.HasNextPage, tc.result.HasNextPage)
					}
					if c.HasPreviousPage != tc.result.HasPreviousPage {
						t.Errorf("Batch previous page mismatch %v got: %v", c.HasPreviousPage, tc.result.HasPreviousPage)
					}
					for index, credential := range c.Credentials {
						validateCredential(t, tc.result.Credentials[index], credential)
					}
				}
			})
		}
	})
}

func TestGetCredentialCount(t *testing.T) {
	// add new agent with no pre-existing credentials
	a, connections := addAgentAndConnections("TestGetCredentialCount")
	size := 5
	fake.AddCredentials(pgDB, a.ID, connections[0].ID, size)

	// Get count
	got, err := pgDB.GetCredentialCount(a.ID)
	if err != nil {
		t.Errorf("Error fetching count %s", err.Error())
	} else if got != size {
		t.Errorf("Mismatch in fetched credential count expected: %v  got: %v", size, got)
	}
}

func TestGetConnectionCredentialCount(t *testing.T) {
	// add new agent with no pre-existing credentials
	a, connections := addAgentAndConnections("TestGetCredentialCount")
	size := 5
	index := 0
	fake.AddCredentials(pgDB, a.ID, connections[index].ID, (index+1)*size)
	index++
	fake.AddCredentials(pgDB, a.ID, connections[index].ID, (index+1)*size)
	index++
	fake.AddCredentials(pgDB, a.ID, connections[index].ID, index*size)

	// Get count
	expected := index * size
	got, err := pgDB.GetConnectionCredentialCount(a.ID, connections[index].ID)
	if err != nil {
		t.Errorf("Error fetching count %s", err.Error())
	} else if got != expected {
		t.Errorf("Mismatch in fetched credential count expected: %v  got: %v", expected, got)
	}

}