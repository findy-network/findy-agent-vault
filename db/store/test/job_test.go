package test

import (
	"sort"
	"testing"
	"time"

	"github.com/findy-network/findy-agent-vault/db/fake"
	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/db/store"
	graph "github.com/findy-network/findy-agent-vault/graph/model"
	"github.com/findy-network/findy-agent-vault/paginator"
)

func validateJob(t *testing.T, exp, got *model.Job) {
	if got == nil {
		t.Errorf("Expecting result, job  is nil")
		return
	}
	if got.ID == "" {
		t.Errorf("Job id invalid.")
	}
	if got.TenantID != exp.TenantID {
		t.Errorf("Job tenant id mismatch expected %s got %s", exp.TenantID, got.TenantID)
	}
	if got.ProtocolType != exp.ProtocolType {
		t.Errorf("Job protocol type mismatch expected %s got %s", exp.ProtocolType, got.ProtocolType)
	}
	validateStrPtr(t, exp.ProtocolConnectionID, got.ProtocolConnectionID, "ProtocolConnectionID")
	validateStrPtr(t, exp.ProtocolCredentialID, got.ProtocolCredentialID, "ProtocolCredentialID")
	validateStrPtr(t, exp.ProtocolProofID, got.ProtocolProofID, "ProtocolProofID")
	validateStrPtr(t, exp.ProtocolMessageID, got.ProtocolMessageID, "ProtocolMessageID")
	validateStrPtr(t, exp.ConnectionID, got.ConnectionID, "ConnectionID")
	if got.Status != exp.Status {
		t.Errorf("Job status mismatch expected %s got %s", exp.Status, got.Status)
	}
	if got.Result != exp.Result {
		t.Errorf("Job result mismatch expected %s got %s", exp.Result, got.Result)
	}
	if got.InitiatedByUs != exp.InitiatedByUs {
		t.Errorf("Job initiatedByUs mismatch expected %v got %v", exp.InitiatedByUs, got.InitiatedByUs)
	}

	validateCreatedTS(t, got.Cursor, &got.Created)
}

func validateJobs(t *testing.T, expCount int, exp, got *model.Jobs) {
	if len(got.Jobs) != expCount {
		t.Errorf("Mismatch in job  count: %v  got: %v", len(got.Jobs), expCount)
	}
	if got.HasNextPage != exp.HasNextPage {
		t.Errorf("Batch next page mismatch %v got: %v", got.HasNextPage, exp.HasNextPage)
	}
	if got.HasPreviousPage != exp.HasPreviousPage {
		t.Errorf("Batch previous page mismatch %v got: %v", got.HasPreviousPage, exp.HasPreviousPage)
	}
	for index, job := range got.Jobs {
		validateJob(t, exp.Jobs[index], job)
	}
}

type jobTest struct {
	name   string
	args   *paginator.BatchInfo
	result *model.Jobs
}

func getJobTests(size int, all []*model.Job) []*jobTest {
	var jobTests = []*jobTest{
		{
			"first 5",
			&paginator.BatchInfo{Count: size, Tail: false},
			&model.Jobs{HasNextPage: true, HasPreviousPage: false, Jobs: all[:size]},
		},
		{
			"first next 5",
			&paginator.BatchInfo{Count: size, Tail: false, After: all[size-1].Cursor},
			&model.Jobs{HasNextPage: true, HasPreviousPage: true, Jobs: all[size : size*2]},
		},
		{
			"first last 5",
			&paginator.BatchInfo{Count: size, Tail: false, After: all[(size*2)-1].Cursor},
			&model.Jobs{HasNextPage: false, HasPreviousPage: true, Jobs: all[size*2:]},
		},
		{
			"last 5",
			&paginator.BatchInfo{Count: size, Tail: true},
			&model.Jobs{HasNextPage: false, HasPreviousPage: true, Jobs: all[size*2:]},
		},
		{
			"last next 5",
			&paginator.BatchInfo{Count: size, Tail: true, Before: all[size*2].Cursor},
			&model.Jobs{HasNextPage: true, HasPreviousPage: true, Jobs: all[size : size*2]},
		},
		{
			"last first 5",
			&paginator.BatchInfo{Count: size, Tail: true, Before: all[size].Cursor},
			&model.Jobs{HasNextPage: true, HasPreviousPage: false, Jobs: all[:size]},
		},
		{
			"all",
			&paginator.BatchInfo{Count: size * 3, Tail: false},
			&model.Jobs{HasNextPage: false, HasPreviousPage: false, Jobs: all},
		},
	}
	return jobTests
}

func TestAddJob(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("add job  "+s.name, func(t *testing.T) {
			testJob = s.newTestJob(testJob)

			// Add data
			j, err := s.db.AddJob(testJob)
			if err != nil {
				t.Errorf("Failed to add job  %s", err.Error())
			} else {
				validateJob(t, testJob, j)
			}

			// Get data for id
			got, err := s.db.GetJob(j.ID, s.testTenantID)
			if err != nil {
				t.Errorf("Error fetching job  %s", err.Error())
			} else {
				validateJob(t, j, got)
			}
		})
	}
}

func TestGetNonexistentJob(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("get non-existent job  "+s.name, func(t *testing.T) {
			_, err := s.db.GetJob("b49b092e-2812-4adc-b3da-79a4ebbf864c", s.testTenantID)
			if err == nil || store.ErrorCode(err) != store.ErrCodeNotFound {
				t.Errorf("Error fetching non-existent job  %s", err)
			}
		})
	}
}

func TestAddJobSameIDDifferentTenant(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("add job same id "+s.name, func(t *testing.T) {
			testJob = s.newTestJob(testJob)
			testJob.ConnectionID = nil

			// Add data
			job1, err := s.db.AddJob(testJob)
			if err != nil {
				t.Errorf("Failed to add job  %s", err.Error())
			} else {
				validateJob(t, testJob, job1)
			}

			// Add data
			a2 := fake.AddAgent(s.db)
			testJob.TenantID = a2.ID
			job2, err := s.db.AddJob(testJob)
			if err != nil {
				t.Errorf("Failed to add job with same id %s", err.Error())
			} else {
				validateJob(t, testJob, job2)
			}
		})
	}
}

func TestUpdateJob(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("update job  "+s.name, func(t *testing.T) {
			testJob = s.newTestJob(testJob)

			// Add data
			j, err := s.db.AddJob(testJob)
			if err != nil {
				t.Errorf("Failed to add job  %s", err.Error())
			}

			// Update data
			j.ProtocolConnectionID = &s.testConnectionID
			j.Status = graph.JobStatusComplete
			j.Result = graph.JobResultSuccess
			j.Updated = time.Now().UTC()
			got, err := s.db.UpdateJob(j)
			if err != nil {
				t.Errorf("Failed to update job %s", err.Error())
			} else {
				if got.Updated.Sub(j.Updated) == 0 || time.Since(got.Updated) > time.Second {
					t.Errorf("Updated timestamp not updated, got: %v was: %v", got.Created, j.Updated)
				}
				validateJob(t, j, got)
			}
		})
	}
}

func TestGetTenantJobs(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("get job s "+s.name, func(t *testing.T) {
			// add new agent with no pre-existing jobs
			a, connections := AddAgentAndConnections(s.db, "TestGetTenantJobs", 3)

			size := 5
			all := fake.AddJobs(s.db, a.ID, connections[0].ID, size)
			all = append(all, fake.AddJobs(s.db, a.ID, connections[1].ID, size)...)
			all = append(all, fake.AddJobs(s.db, a.ID, connections[2].ID, size)...)

			sort.Slice(all, func(i, j int) bool {
				return all[i].Created.Sub(all[j].Created) < 0
			})

			t.Run("get job s", func(t *testing.T) {
				tests := getJobTests(size, all)

				for _, testCase := range tests {
					tc := testCase
					t.Run(tc.name, func(t *testing.T) {
						completed := true
						c, err := s.db.GetJobs(tc.args, a.ID, nil, &completed)
						if err != nil {
							t.Errorf("Error fetching job s %s", err.Error())
						} else {
							validateJobs(t, tc.args.Count, c, tc.result)
						}
					})
				}
			})
		})
	}
}

func TestGetConnectionJobs(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("get connection job s "+s.name, func(t *testing.T) {
			// add new agent with no pre-existing job s
			a, connections := AddAgentAndConnections(s.db, "TestGetConnectionJobs", 3)

			size := 5
			countPerConnection := size * 3
			fake.AddJobs(s.db, a.ID, connections[0].ID, countPerConnection)
			fake.AddJobs(s.db, a.ID, connections[1].ID, countPerConnection)
			all := fake.AddJobs(s.db, a.ID, connections[2].ID, countPerConnection)

			sort.Slice(all, func(i, j int) bool {
				return all[i].Created.Sub(all[j].Created) < 0
			})

			t.Run("get job s", func(t *testing.T) {
				tests := getJobTests(size, all)

				for _, testCase := range tests {
					tc := testCase
					t.Run(tc.name, func(t *testing.T) {
						completed := true
						c, err := s.db.GetJobs(tc.args, a.ID, &connections[2].ID, &completed)
						if err != nil {
							t.Errorf("Error fetching connection job s %s", err.Error())
						} else {
							validateJobs(t, tc.args.Count, c, tc.result)
						}
					})
				}
			})
		})
	}
}

func TestGetJobCount(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("get job s count "+s.name, func(t *testing.T) {
			// add new agent with no pre-existing job s
			a, connections := AddAgentAndConnections(s.db, "TestGetJobCount", 3)
			size := 5
			fake.AddJobs(s.db, a.ID, connections[0].ID, size)

			// Get count
			completed := true
			got, err := s.db.GetJobCount(a.ID, nil, &completed)
			if err != nil {
				t.Errorf("Error fetching count %s", err.Error())
			} else if got != size {
				t.Errorf("Mismatch in fetched job  count expected: %v  got: %v", size, got)
			}
		})
	}
}

func TestGetConnectionJobCount(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("get connection job s count "+s.name, func(t *testing.T) {
			// add new agent with no pre-existing job s
			a, connections := AddAgentAndConnections(s.db, "TestGetConnectionJobCount", 3)
			size := 5
			index := 0
			fake.AddJobs(s.db, a.ID, connections[index].ID, (index+1)*size)
			index++
			fake.AddJobs(s.db, a.ID, connections[index].ID, (index+1)*size)
			index++
			fake.AddJobs(s.db, a.ID, connections[index].ID, index*size)

			// Get count
			expected := index * size
			completed := true
			got, err := s.db.GetJobCount(a.ID, &connections[index].ID, &completed)
			if err != nil {
				t.Errorf("Error fetching count %s", err.Error())
			} else if got != expected {
				t.Errorf("Mismatch in fetched job  count expected: %v  got: %v", expected, got)
			}
		})
	}
}

func TestGetConnectionForJob(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("get connection for job "+s.name, func(t *testing.T) {
			a, connections := AddAgentAndConnections(s.db, "TestGetConnectionForJob", 3)
			connection := connections[0]
			jobs := fake.AddJobs(s.db, a.ID, connection.ID, 1)
			job := jobs[0]

			// Get data for id
			got, err := s.db.GetConnectionForJob(job.ID, a.ID)
			if err != nil {
				t.Errorf("Error fetching connection %s", err.Error())
			} else {
				validateConnection(t, connection, got)
			}
		})
	}
}

func TestGetConnectionOutputForJob(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("get connection output for job "+s.name, func(t *testing.T) {
			a, connections := AddAgentAndConnections(s.db, "TestGetConnectionOutputForJob", 3)
			connection := connections[0]
			jobs := fake.AddConnectionJobs(s.db, a.ID, connection.ID, connection.ID, 1)
			job := jobs[0]

			// Get data for id
			got, err := s.db.GetJobOutput(job.ID, a.ID, graph.ProtocolTypeConnection)
			if err != nil {
				t.Errorf("Error fetching output %s", err.Error())
			} else {
				validateConnection(t, connection, got.Connection)
			}
		})
	}
}

func TestGetCredentialOutputForJob(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("get credential output for job "+s.name, func(t *testing.T) {
			a, connections := AddAgentAndConnections(s.db, "TestGetCredentialOutputForJob", 3)
			connection := connections[0]
			credentials := fake.AddCredentials(s.db, a.ID, connection.ID, 1)
			credential := credentials[0]
			jobs := fake.AddCredentialJobs(s.db, a.ID, connection.ID, credential.ID, 1)
			job := jobs[0]

			// Get data for id
			got, err := s.db.GetJobOutput(job.ID, a.ID, graph.ProtocolTypeCredential)
			if err != nil {
				t.Errorf("Error fetching output %s", err.Error())
			} else {
				validateCredential(t, credential, got.Credential)
			}
		})
	}
}

func TestGetProofOutputForJob(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("get proof output for job "+s.name, func(t *testing.T) {
			a, connections := AddAgentAndConnections(s.db, "TestGetProofOutputForJob", 3)
			connection := connections[0]
			proofs := fake.AddProofs(s.db, a.ID, connection.ID, 1, true)
			proof := proofs[0]
			jobs := fake.AddProofJobs(s.db, a.ID, connection.ID, proof.ID, 1, graph.JobStatusComplete)
			job := jobs[0]

			// Get data for id
			got, err := s.db.GetJobOutput(job.ID, a.ID, graph.ProtocolTypeProof)
			if err != nil {
				t.Errorf("Error fetching output %s", err.Error())
			} else {
				validateProof(t, proof, got.Proof)
			}
		})
	}
}

func TestGetMessageOutputForJob(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("get message output for job "+s.name, func(t *testing.T) {
			a, connections := AddAgentAndConnections(s.db, "TestGetMessageOutputForJob", 3)
			connection := connections[0]
			messages := fake.AddMessages(s.db, a.ID, connection.ID, 1)
			message := messages[0]
			jobs := fake.AddMessageJobs(s.db, a.ID, connection.ID, message.ID, 1)
			job := jobs[0]

			// Get data for id
			got, err := s.db.GetJobOutput(job.ID, a.ID, graph.ProtocolTypeBasicMessage)
			if err != nil {
				t.Errorf("Error fetching output %s", err.Error())
			} else {
				validateMessage(t, message, got.Message)
			}
		})
	}
}

func TestGetOpenProofJobs(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("get open proof jobs "+s.name, func(t *testing.T) {
			a, connections := AddAgentAndConnections(s.db, "TestGetOpenProofJobs", 3)
			connection := connections[0]
			proofs := fake.AddProofs(s.db, a.ID, connection.ID, 1, false)
			proof := proofs[0]
			jobs := fake.AddProofJobs(s.db, a.ID, connection.ID, proof.ID, 1, graph.JobStatusBlocked)
			job := jobs[0]

			// Get data for id
			got, err := s.db.GetOpenProofJobs(job.TenantID, proof.Attributes)
			if err != nil {
				t.Errorf("Error getting jobs %s", err.Error())
			} else if got[0].ID != job.ID {
				t.Errorf("Open proof job was not found %s != %s", got[0].ID, job.ID)
			}
		})
	}
}
