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

func validateProofAttributes(t *testing.T, exp, got []*graph.ProofAttribute) {
	if len(got) == 0 {
		t.Errorf("No attributes found")
	}
	for index, a := range got {
		if a.ID == "" {
			t.Errorf("Proof attribute id invalid.")
		}
		if a.Name != exp[index].Name {
			t.Errorf("Proof attribute name mismatch: expected %s got %s.", exp[index].Name, a.Name)
		}
		if a.CredDefID != exp[index].CredDefID {
			t.Errorf("Proof attribute cred def mismatch: expected %s got %s.", exp[index].CredDefID, a.CredDefID)
		}
	}
}

func validateProofValues(t *testing.T, exp, got []*graph.ProofValue) {
	if len(got) != len(exp) {
		t.Errorf("No expected values found")
	}
	for index, a := range got {
		if a.ID == "" {
			t.Errorf("Proof value id invalid.")
		}
		if a.AttributeID != exp[index].AttributeID {
			t.Errorf("Proof value attribute id mismatch: expected %s got %s.", exp[index].AttributeID, a.AttributeID)
		}
		if a.Value != exp[index].Value {
			t.Errorf("Proof value mismatch: expected %s got %s.", exp[index].Value, a.Value)
		}
	}
}

func validateProof(t *testing.T, exp, got *model.Proof) {
	if got == nil {
		t.Errorf("Expecting result, proof is nil")
		return
	}
	if got.ID == "" {
		t.Errorf("Proof id invalid.")
	}
	if got.TenantID != exp.TenantID {
		t.Errorf("Proof tenant id mismatch expected %s got %s", exp.TenantID, got.TenantID)
	}
	if got.ConnectionID != exp.ConnectionID {
		t.Errorf("Proof connection id mismatch expected %s got %s", exp.ConnectionID, got.ConnectionID)
	}
	if got.Role != exp.Role {
		t.Errorf("Proof Role mismatch expected %s got %s", exp.Role, got.Role)
	}
	if got.InitiatedByUs != exp.InitiatedByUs {
		t.Errorf("Proof InitiatedByUs mismatch expected %v got %v", exp.InitiatedByUs, got.InitiatedByUs)
	}
	if got.Result != exp.Result {
		t.Errorf("Proof Result mismatch expected %v got %v", exp.Result, got.Result)
	}
	validateTimestap(t, &exp.Approved, &got.Approved, "Approved")
	validateTimestap(t, &exp.Verified, &got.Verified, "Verified")
	validateTimestap(t, &exp.Failed, &got.Failed, "Failed")
	validateTimestap(t, &exp.Archived, &got.Archived, "Archived")
	validateCreatedTS(t, got.Cursor, &got.Created)

	validateProofAttributes(t, exp.Attributes, got.Attributes)
	validateProofValues(t, exp.Values, got.Values)
}

func validateProofs(t *testing.T, expCount int, exp, got *model.Proofs) {
	if len(got.Proofs) != expCount {
		t.Errorf("Mismatch in proof count: %v  got: %v", len(got.Proofs), expCount)
	}
	if got.HasNextPage != exp.HasNextPage {
		t.Errorf("Batch next page mismatch %v got: %v", got.HasNextPage, exp.HasNextPage)
	}
	if got.HasPreviousPage != exp.HasPreviousPage {
		t.Errorf("Batch previous page mismatch %v got: %v", got.HasPreviousPage, exp.HasPreviousPage)
	}
	for index, proof := range got.Proofs {
		validateProof(t, exp.Proofs[index], proof)
	}
}

type proofTest struct {
	name   string
	args   *paginator.BatchInfo
	result *model.Proofs
}

func getProofTests(size int, all []*model.Proof) []*proofTest {
	var proofTests = []*proofTest{
		{
			"first 5",
			&paginator.BatchInfo{Count: size, Tail: false},
			&model.Proofs{HasNextPage: true, HasPreviousPage: false, Proofs: all[:size]},
		},
		{
			"first next 5",
			&paginator.BatchInfo{Count: size, Tail: false, After: all[size-1].Cursor},
			&model.Proofs{HasNextPage: true, HasPreviousPage: true, Proofs: all[size : size*2]},
		},
		{
			"first last 5",
			&paginator.BatchInfo{Count: size, Tail: false, After: all[(size*2)-1].Cursor},
			&model.Proofs{HasNextPage: false, HasPreviousPage: true, Proofs: all[size*2:]},
		},
		{
			"last 5",
			&paginator.BatchInfo{Count: size, Tail: true},
			&model.Proofs{HasNextPage: false, HasPreviousPage: true, Proofs: all[size*2:]},
		},
		{
			"last next 5",
			&paginator.BatchInfo{Count: size, Tail: true, Before: all[size*2].Cursor},
			&model.Proofs{HasNextPage: true, HasPreviousPage: true, Proofs: all[size : size*2]},
		},
		{
			"last first 5",
			&paginator.BatchInfo{Count: size, Tail: true, Before: all[size].Cursor},
			&model.Proofs{HasNextPage: true, HasPreviousPage: false, Proofs: all[:size]},
		},
		{
			"all",
			&paginator.BatchInfo{Count: size * 3, Tail: false},
			&model.Proofs{HasNextPage: false, HasPreviousPage: false, Proofs: all},
		},
	}
	return proofTests
}

func TestAddProof(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("add proof "+s.name, func(t *testing.T) {
			testProof = s.newTestProof(testProof)

			// Add data
			p, err := s.db.AddProof(testProof)
			if err != nil {
				t.Errorf("Failed to add proof %s", err.Error())
			} else {
				validateProof(t, testProof, p)
			}

			// Get data for id
			got, err := s.db.GetProof(p.ID, s.testTenantID)
			if err != nil {
				t.Errorf("Error fetching proof %s", err.Error())
			} else {
				validateProof(t, p, got)
			}
		})
	}
}

func TestUpdateProof(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("update proof "+s.name, func(t *testing.T) {
			testProof.TenantID = s.testTenantID
			testProof.ConnectionID = s.testConnectionID

			// Add data
			p, err := s.db.AddProof(testProof)
			if err != nil {
				t.Errorf("Failed to add proof %s", err.Error())
			}

			// Update data
			now := time.Now().UTC()
			p.Approved = now
			p.Verified = now
			p.Values = make([]*graph.ProofValue, 0)
			for _, attr := range p.Attributes {
				p.Values = append(p.Values, &graph.ProofValue{ID: attr.ID, AttributeID: attr.ID, Value: "value"})
			}
			_, err = s.db.UpdateProof(p)
			if err != nil {
				t.Errorf("Failed to update proof %s", err.Error())
			}

			// Get data for id
			got, err := s.db.GetProof(p.ID, s.testTenantID)
			if err != nil {
				t.Errorf("Error fetching proof %s", err.Error())
			} else {
				validateProof(t, p, got)
			}
		})
	}
}

func TestGetTenantProofs(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("get proofs "+s.name, func(t *testing.T) {
			a, connections := AddAgentAndConnections(s.db, "TestGetTenantProofs", 3)

			size := 5
			all := fake.AddProofs(s.db, a.ID, connections[0].ID, size, true)
			all = append(all, fake.AddProofs(s.db, a.ID, connections[1].ID, size, true)...)
			all = append(all, fake.AddProofs(s.db, a.ID, connections[2].ID, size, true)...)

			sort.Slice(all, func(i, j int) bool {
				return all[i].Created.Sub(all[j].Created) < 0
			})

			t.Run("get proofs", func(t *testing.T) {
				tests := getProofTests(size, all)

				for _, testCase := range tests {
					tc := testCase
					t.Run(tc.name, func(t *testing.T) {
						p, err := s.db.GetProofs(tc.args, a.ID, nil)
						if err != nil {
							t.Errorf("Error fetching proofs %s", err.Error())
						} else {
							validateProofs(t, tc.args.Count, p, tc.result)
						}
					})
				}
			})
		})
	}
}

func TestGetConnectionProofs(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("get connection proofs "+s.name, func(t *testing.T) {
			// add new agent with no pre-existing proofs
			a, connections := AddAgentAndConnections(s.db, "TestGetConnectionProofs", 3)

			size := 5
			countPerConnection := size * 3
			fake.AddProofs(s.db, a.ID, connections[0].ID, countPerConnection, true)
			fake.AddProofs(s.db, a.ID, connections[1].ID, countPerConnection, true)
			all := fake.AddProofs(s.db, a.ID, connections[2].ID, countPerConnection, true)

			sort.Slice(all, func(i, j int) bool {
				return all[i].Created.Sub(all[j].Created) < 0
			})

			t.Run("get proofs", func(t *testing.T) {
				tests := getProofTests(size, all)

				for _, testCase := range tests {
					tc := testCase
					t.Run(tc.name, func(t *testing.T) {
						p, err := s.db.GetProofs(tc.args, a.ID, &connections[2].ID)
						if err != nil {
							t.Errorf("Error fetching connection proofs %s", err.Error())
						} else {
							validateProofs(t, tc.args.Count, p, tc.result)
						}
					})
				}
			})
		})
	}
}

func TestGetProofCount(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("get proofs count "+s.name, func(t *testing.T) {
			// add new agent with no pre-existing proofs
			a, connections := AddAgentAndConnections(s.db, "TestGetProofCount", 3)
			size := 5
			fake.AddProofs(s.db, a.ID, connections[0].ID, size, true)

			// Get count
			got, err := s.db.GetProofCount(a.ID, nil)
			if err != nil {
				t.Errorf("Error fetching count %s", err.Error())
			} else if got != size {
				t.Errorf("Mismatch in fetched proof count expected: %v  got: %v", size, got)
			}
		})
	}
}

func TestGetConnectionProofCount(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("get connection proofs count "+s.name, func(t *testing.T) {
			// add new agent with no pre-existing proofs
			a, connections := AddAgentAndConnections(s.db, "TestGetConnectionProofCount", 3)
			size := 5
			index := 0
			fake.AddProofs(s.db, a.ID, connections[index].ID, (index+1)*size, true)
			index++
			fake.AddProofs(s.db, a.ID, connections[index].ID, (index+1)*size, true)
			index++
			fake.AddProofs(s.db, a.ID, connections[index].ID, index*size, true)

			// Get count
			expected := index * size
			got, err := s.db.GetProofCount(a.ID, &connections[index].ID)
			if err != nil {
				t.Errorf("Error fetching count %s", err.Error())
			} else if got != expected {
				t.Errorf("Mismatch in fetched proof count expected: %v  got: %v", expected, got)
			}
		})
	}
}

func TestGetConnectionForProof(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("get connection for proof "+s.name, func(t *testing.T) {
			a, connections := AddAgentAndConnections(s.db, "TestGetConnectionForProof", 3)
			connection := connections[0]
			proofs := fake.AddProofs(s.db, a.ID, connection.ID, 1, true)
			proof := proofs[0]

			// Get data for id
			got, err := s.db.GetConnectionForProof(proof.ID, a.ID)
			if err != nil {
				t.Errorf("Error fetching connection %s", err.Error())
			} else {
				validateConnection(t, connection, got)
			}
		})
	}
}

func TestArchiveProof(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("archive proof "+s.name, func(t *testing.T) {
			testProof = s.newTestProof(testProof)

			// Add data
			p, err := s.db.AddProof(testProof)
			assert.D.True(err == nil)

			now := utils.CurrentTime()
			err = s.db.ArchiveProof(p.ID, p.TenantID)
			if err != nil {
				t.Errorf("Failed to archive proof %s", err.Error())
			}

			// Get data for id
			got, err := s.db.GetProof(p.ID, p.TenantID)
			assert.D.True(err == nil)

			p.Archived = now
			validateProof(t, p, got)
		})
	}
}
