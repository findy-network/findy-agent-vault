package test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/findy-network/findy-agent-vault/db/fake"
	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/paginator"
)

func validateEvent(t *testing.T, exp, got *model.Event) {
	if got == nil {
		t.Errorf("Expecting result, event  is nil")
		return
	}
	if got.ID == "" {
		t.Errorf("Event id invalid.")
	}
	if got.TenantID != exp.TenantID {
		t.Errorf("Event tenant id mismatch expected %s got %s", exp.TenantID, got.TenantID)
	}
	validateStrPtr(t, exp.ConnectionID, got.ConnectionID, "ConnectionID")
	validateStrPtr(t, exp.JobID, got.JobID, "JobID")
	if got.Description != exp.Description {
		t.Errorf("Event Description mismatch expected %s got %s", exp.Description, got.Description)
	}
	if got.Read != exp.Read {
		t.Errorf("Event Read mismatch expected %v got %v", exp.Read, got.Read)
	}
	validateCreatedTS(t, got.Cursor, &got.Created)
}

func validateEvents(t *testing.T, expCount int, exp, got *model.Events) {
	if len(got.Events) != expCount {
		t.Errorf("Mismatch in event  count: %v  got: %v", len(got.Events), expCount)
	}
	if got.HasNextPage != exp.HasNextPage {
		t.Errorf("Batch next page mismatch %v got: %v", got.HasNextPage, exp.HasNextPage)
	}
	if got.HasPreviousPage != exp.HasPreviousPage {
		t.Errorf("Batch previous page mismatch %v got: %v", got.HasPreviousPage, exp.HasPreviousPage)
	}
	for index, event := range got.Events {
		validateEvent(t, exp.Events[index], event)
	}
}

type eventTest struct {
	name   string
	args   *paginator.BatchInfo
	result *model.Events
}

func getEventTests(size int, all []*model.Event) []*eventTest {
	var eventTests = []*eventTest{
		{
			"first 5",
			&paginator.BatchInfo{Count: size, Tail: false},
			&model.Events{HasNextPage: true, HasPreviousPage: false, Events: all[:size]},
		},
		{
			"first next 5",
			&paginator.BatchInfo{Count: size, Tail: false, After: all[size-1].Cursor},
			&model.Events{HasNextPage: true, HasPreviousPage: true, Events: all[size : size*2]},
		},
		{
			"first last 5",
			&paginator.BatchInfo{Count: size, Tail: false, After: all[(size*2)-1].Cursor},
			&model.Events{HasNextPage: false, HasPreviousPage: true, Events: all[size*2:]},
		},
		{
			"last 5",
			&paginator.BatchInfo{Count: size, Tail: true},
			&model.Events{HasNextPage: false, HasPreviousPage: true, Events: all[size*2:]},
		},
		{
			"last next 5",
			&paginator.BatchInfo{Count: size, Tail: true, Before: all[size*2].Cursor},
			&model.Events{HasNextPage: true, HasPreviousPage: true, Events: all[size : size*2]},
		},
		{
			"last first 5",
			&paginator.BatchInfo{Count: size, Tail: true, Before: all[size].Cursor},
			&model.Events{HasNextPage: true, HasPreviousPage: false, Events: all[:size]},
		},
		{
			"all",
			&paginator.BatchInfo{Count: size * 3, Tail: false},
			&model.Events{HasNextPage: false, HasPreviousPage: false, Events: all},
		},
	}
	return eventTests
}

func TestAddEvent(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("add event  "+s.name, func(t *testing.T) {
			testEvent = model.NewEvent(s.testTenantID, testEvent)
			testEvent.TenantID = s.testTenantID
			testEvent.ConnectionID = &s.testConnectionID

			// Add data
			c, err := s.db.AddEvent(testEvent)
			if err != nil {
				t.Errorf("Failed to add event  %s", err.Error())
			} else {
				validateEvent(t, testEvent, c)
			}

			// Get data for id
			got, err := s.db.GetEvent(c.ID, s.testTenantID)
			if err != nil {
				t.Errorf("Error fetching event  %s", err.Error())
			} else if !reflect.DeepEqual(&c, &got) {
				t.Errorf("Mismatch in fetched event  expected: %v  got: %v", c, got)
			}
			validateEvent(t, c, got)
		})
	}
}

func TestMarkEventRead(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("update event  "+s.name, func(t *testing.T) {
			testEvent.TenantID = s.testTenantID
			testEvent.ConnectionID = &s.testConnectionID

			// Add data
			e, err := s.db.AddEvent(testEvent)
			if err != nil {
				t.Errorf("Failed to add event  %s", err.Error())
			}

			// Update data
			e.Read = true
			got, err := s.db.MarkEventRead(e.ID, testEvent.TenantID)
			if err != nil {
				t.Errorf("Failed to mark event read  %s", err.Error())
			}
			if !reflect.DeepEqual(&e, &got) {
				t.Errorf("Mismatch in fetched event  expected: %v  got: %v", e, got)
			}
			validateEvent(t, e, got)
		})
	}
}

func TestGetTenantEvents(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("get event s "+s.name, func(t *testing.T) {
			// add new agent with no pre-existing events
			a, connections := AddAgentAndConnections(s.db, "TestGetTenantEvents", 3)

			size := 5
			all := fake.AddEvents(s.db, a.ID, connections[0].ID, nil, size)
			all = append(all, fake.AddEvents(s.db, a.ID, connections[1].ID, nil, size)...)
			all = append(all, fake.AddEvents(s.db, a.ID, connections[2].ID, nil, size)...)

			sort.Slice(all, func(i, j int) bool {
				return all[i].Created.Sub(all[j].Created) < 0
			})

			t.Run("get event s", func(t *testing.T) {
				tests := getEventTests(size, all)

				for _, testCase := range tests {
					tc := testCase
					t.Run(tc.name, func(t *testing.T) {
						c, err := s.db.GetEvents(tc.args, a.ID, nil)
						if err != nil {
							t.Errorf("Error fetching event s %s", err.Error())
						} else {
							validateEvents(t, tc.args.Count, c, tc.result)
						}
					})
				}
			})
		})
	}
}

func TestGetConnectionEvents(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("get connection event s "+s.name, func(t *testing.T) {
			// add new agent with no pre-existing event s
			a, connections := AddAgentAndConnections(s.db, "TestGetConnectionEvents", 3)

			size := 5
			countPerConnection := size * 3
			fake.AddEvents(s.db, a.ID, connections[0].ID, nil, countPerConnection)
			fake.AddEvents(s.db, a.ID, connections[1].ID, nil, countPerConnection)
			all := fake.AddEvents(s.db, a.ID, connections[2].ID, nil, countPerConnection)

			sort.Slice(all, func(i, j int) bool {
				return all[i].Created.Sub(all[j].Created) < 0
			})

			t.Run("get event s", func(t *testing.T) {
				tests := getEventTests(size, all)

				for _, testCase := range tests {
					tc := testCase
					t.Run(tc.name, func(t *testing.T) {
						c, err := s.db.GetEvents(tc.args, a.ID, &connections[2].ID)
						if err != nil {
							t.Errorf("Error fetching connection event s %s", err.Error())
						} else {
							validateEvents(t, tc.args.Count, c, tc.result)
						}
					})
				}
			})
		})
	}
}

func TestGetEventCount(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("get event s count "+s.name, func(t *testing.T) {
			// add new agent with no pre-existing event s
			a, connections := AddAgentAndConnections(s.db, "TestGetEventCount", 3)
			size := 5
			fake.AddEvents(s.db, a.ID, connections[0].ID, nil, size)

			// Get count
			got, err := s.db.GetEventCount(a.ID, nil)
			if err != nil {
				t.Errorf("Error fetching count %s", err.Error())
			} else if got != size {
				t.Errorf("Mismatch in fetched event  count expected: %v  got: %v", size, got)
			}
		})
	}
}

func TestGetConnectionEventCount(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("get connection event s count "+s.name, func(t *testing.T) {
			// add new agent with no pre-existing event s
			a, connections := AddAgentAndConnections(s.db, "TestGetConnectionEventCount", 3)
			size := 5
			index := 0
			fake.AddEvents(s.db, a.ID, connections[index].ID, nil, (index+1)*size)
			index++
			fake.AddEvents(s.db, a.ID, connections[index].ID, nil, (index+1)*size)
			index++
			fake.AddEvents(s.db, a.ID, connections[index].ID, nil, index*size)

			// Get count
			expected := index * size
			got, err := s.db.GetEventCount(a.ID, &connections[index].ID)
			if err != nil {
				t.Errorf("Error fetching count %s", err.Error())
			} else if got != expected {
				t.Errorf("Mismatch in fetched event  count expected: %v  got: %v", expected, got)
			}
		})
	}
}

func TestGetConnectionForEvent(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("get connection for event"+s.name, func(t *testing.T) {
			a, connections := AddAgentAndConnections(s.db, "TestGetConnectionForEvent", 3)
			connection := connections[0]
			events := fake.AddEvents(s.db, a.ID, connection.ID, nil, 1)
			event := events[0]

			// Get data for id
			got, err := s.db.GetConnectionForEvent(event.ID, a.ID)
			if err != nil {
				t.Errorf("Error fetching connection %s", err.Error())
			} else {
				validateConnection(t, connection, got)
			}
		})
	}
}

func TestGetJobForEvent(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("get job for event "+s.name, func(t *testing.T) {
			a, connections := AddAgentAndConnections(s.db, "TestGetJobForEvent", 3)
			connection := connections[0]
			jobs := fake.AddJobs(s.db, a.ID, connection.ID, 1)
			job := jobs[0]
			events := fake.AddEvents(s.db, a.ID, connection.ID, &job.ID, 1)
			event := events[0]

			// Get data for id
			got, err := s.db.GetJobForEvent(event.ID, a.ID)
			if err != nil {
				t.Errorf("Error fetching connection %s", err.Error())
			} else {
				validateJob(t, job, got)
			}
		})
	}
}
