package test

import (
	"math"
	"sort"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/findy-network/findy-agent-vault/db/fake"
	"github.com/findy-network/findy-agent-vault/db/model"
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
	validateStrPtr(t, exp.ProtocolID, got.ProtocolID, "ProtocolID")
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

	if time.Since(got.Created) > time.Second {
		t.Errorf("Timestamp not in threshold %v", got.Created)
	}
	created := uint64(math.Round(float64(got.Created.UnixNano()) / float64(time.Millisecond.Nanoseconds())))
	if got.Cursor != created {
		t.Errorf("Cursor mismatch %v %v", got.Cursor, created)
	}
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
			testJob = model.NewJob(testJob)
			testJob.TenantID = s.testTenantID
			testJob.ConnectionID = &s.testConnectionID

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

func TestUpdateJob(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("update job  "+s.name, func(t *testing.T) {
			testJob.TenantID = s.testTenantID
			testJob.ConnectionID = &s.testConnectionID

			// Add data
			j, err := s.db.AddJob(testJob)
			if err != nil {
				t.Errorf("Failed to add job  %s", err.Error())
			}

			// Update data
			pID := faker.UUIDHyphenated()
			j.ProtocolID = &pID
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
						c, err := s.db.GetJobs(tc.args, a.ID, nil)
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
						c, err := s.db.GetJobs(tc.args, a.ID, &connections[2].ID)
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
			got, err := s.db.GetJobCount(a.ID, nil)
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
			got, err := s.db.GetJobCount(a.ID, &connections[index].ID)
			if err != nil {
				t.Errorf("Error fetching count %s", err.Error())
			} else if got != expected {
				t.Errorf("Mismatch in fetched job  count expected: %v  got: %v", expected, got)
			}
		})
	}
}
