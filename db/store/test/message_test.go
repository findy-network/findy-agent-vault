package test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/findy-network/findy-agent-vault/db/fake"
	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/paginator"
	"github.com/findy-network/findy-agent-vault/utils"
	"github.com/lainio/err2/assert"
)

func validateBoolPtr(t *testing.T, exp, got *bool, name string) {
	fail := false
	if got != exp {
		fail = true
		if got != nil && exp != nil && *got == *exp {
			fail = false
		}
	}
	if fail {
		t.Errorf("Message %s mismatch expected %v got %v", name, exp, got)
	}
}

func validateMessage(t *testing.T, exp, got *model.Message) {
	if got == nil {
		t.Errorf("Expecting result, message  is nil")
		return
	}
	if got.ID == "" {
		t.Errorf("Message id invalid.")
	}
	if got.TenantID != exp.TenantID {
		t.Errorf("Message tenant id mismatch expected %s got %s", exp.TenantID, got.TenantID)
	}
	if got.ConnectionID != exp.ConnectionID {
		t.Errorf("Message connection id mismatch expected %s got %s", exp.ConnectionID, got.ConnectionID)
	}
	if got.Message != exp.Message {
		t.Errorf("Message Message mismatch expected %s got %s", exp.Message, got.Message)
	}
	if got.SentByMe != exp.SentByMe {
		t.Errorf("Message SentByMe mismatch expected %v got %v", exp.SentByMe, got.SentByMe)
	}
	validateBoolPtr(t, exp.Delivered, got.Delivered, "Delivered")
	validateCreatedTS(t, got.Cursor, &got.Created)
	validateTimestap(t, exp.Archived, got.Archived, "Archived")
}

func validateMessages(t *testing.T, expCount int, exp, got *model.Messages) {
	if len(got.Messages) != expCount {
		t.Errorf("Mismatch in message  count: %v  got: %v", len(got.Messages), expCount)
	}
	if got.HasNextPage != exp.HasNextPage {
		t.Errorf("Batch next page mismatch %v got: %v", got.HasNextPage, exp.HasNextPage)
	}
	if got.HasPreviousPage != exp.HasPreviousPage {
		t.Errorf("Batch previous page mismatch %v got: %v", got.HasPreviousPage, exp.HasPreviousPage)
	}
	for index, message := range got.Messages {
		validateMessage(t, exp.Messages[index], message)
	}
}

type messageTest struct {
	name   string
	args   *paginator.BatchInfo
	result *model.Messages
}

func getMessageTests(size int, all []*model.Message) []*messageTest {
	var messageTests = []*messageTest{
		{
			"first 5",
			&paginator.BatchInfo{Count: size, Tail: false},
			&model.Messages{HasNextPage: true, HasPreviousPage: false, Messages: all[:size]},
		},
		{
			"first next 5",
			&paginator.BatchInfo{Count: size, Tail: false, After: all[size-1].Cursor},
			&model.Messages{HasNextPage: true, HasPreviousPage: true, Messages: all[size : size*2]},
		},
		{
			"first last 5",
			&paginator.BatchInfo{Count: size, Tail: false, After: all[(size*2)-1].Cursor},
			&model.Messages{HasNextPage: false, HasPreviousPage: true, Messages: all[size*2:]},
		},
		{
			"last 5",
			&paginator.BatchInfo{Count: size, Tail: true},
			&model.Messages{HasNextPage: false, HasPreviousPage: true, Messages: all[size*2:]},
		},
		{
			"last next 5",
			&paginator.BatchInfo{Count: size, Tail: true, Before: all[size*2].Cursor},
			&model.Messages{HasNextPage: true, HasPreviousPage: true, Messages: all[size : size*2]},
		},
		{
			"last first 5",
			&paginator.BatchInfo{Count: size, Tail: true, Before: all[size].Cursor},
			&model.Messages{HasNextPage: true, HasPreviousPage: false, Messages: all[:size]},
		},
		{
			"all",
			&paginator.BatchInfo{Count: size * 3, Tail: false},
			&model.Messages{HasNextPage: false, HasPreviousPage: false, Messages: all},
		},
	}
	return messageTests
}

func TestAddMessage(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("add message  "+s.name, func(t *testing.T) {
			testMessage = s.newTestMessage(testMessage)

			// Add data
			m, err := s.db.AddMessage(testMessage)
			if err != nil {
				t.Errorf("Failed to add message  %s", err.Error())
			} else {
				validateMessage(t, testMessage, m)
			}

			// Get data for id
			got, err := s.db.GetMessage(m.ID, s.testTenantID)
			if err != nil {
				t.Errorf("Error fetching message  %s", err.Error())
			} else if !reflect.DeepEqual(&m, &got) {
				t.Errorf("Mismatch in fetched message  expected: %v  got: %v", m, got)
			}
			validateMessage(t, m, got)
		})
	}
}

func TestUpdateMessage(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("update message  "+s.name, func(t *testing.T) {
			testMessage.TenantID = s.testTenantID
			testMessage.ConnectionID = s.testConnectionID

			// Add data
			m, err := s.db.AddMessage(testMessage)
			if err != nil {
				t.Errorf("Failed to add message  %s", err.Error())
			}

			// Update data
			delivered := true
			m.Delivered = &delivered
			got, err := s.db.UpdateMessage(m)
			if err != nil {
				t.Errorf("Failed to update message  %s", err.Error())
			}
			if !reflect.DeepEqual(&m, &got) {
				t.Errorf("Mismatch in fetched message  expected: %v  got: %v", m, got)
			}
			validateMessage(t, m, got)
		})
	}
}

func TestGetTenantMessages(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("get message s "+s.name, func(t *testing.T) {
			// add new agent with no pre-existing messages
			a, connections := AddAgentAndConnections(s.db, "TestGetTenantMessages", 3)

			size := 5
			all := fake.AddMessages(s.db, a.ID, connections[0].ID, size)
			all = append(all, fake.AddMessages(s.db, a.ID, connections[1].ID, size)...)
			all = append(all, fake.AddMessages(s.db, a.ID, connections[2].ID, size)...)

			sort.Slice(all, func(i, j int) bool {
				return all[i].Created.Sub(all[j].Created) < 0
			})

			t.Run("get message s", func(t *testing.T) {
				tests := getMessageTests(size, all)

				for _, testCase := range tests {
					tc := testCase
					t.Run(tc.name, func(t *testing.T) {
						m, err := s.db.GetMessages(tc.args, a.ID, nil)
						if err != nil {
							t.Errorf("Error fetching message s %s", err.Error())
						} else {
							validateMessages(t, tc.args.Count, m, tc.result)
						}
					})
				}
			})
		})
	}
}

func TestGetConnectionMessages(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("get connection message s "+s.name, func(t *testing.T) {
			// add new agent with no pre-existing message s
			a, connections := AddAgentAndConnections(s.db, "TestGetConnectionMessages", 3)

			size := 5
			countPerConnection := size * 3
			fake.AddMessages(s.db, a.ID, connections[0].ID, countPerConnection)
			fake.AddMessages(s.db, a.ID, connections[1].ID, countPerConnection)
			all := fake.AddMessages(s.db, a.ID, connections[2].ID, countPerConnection)

			sort.Slice(all, func(i, j int) bool {
				return all[i].Created.Sub(all[j].Created) < 0
			})

			t.Run("get message s", func(t *testing.T) {
				tests := getMessageTests(size, all)

				for _, testCase := range tests {
					tc := testCase
					t.Run(tc.name, func(t *testing.T) {
						m, err := s.db.GetMessages(tc.args, a.ID, &connections[2].ID)
						if err != nil {
							t.Errorf("Error fetching connection message s %s", err.Error())
						} else {
							validateMessages(t, tc.args.Count, m, tc.result)
						}
					})
				}
			})
		})
	}
}

func TestGetMessageCount(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("get message s count "+s.name, func(t *testing.T) {
			// add new agent with no pre-existing message s
			a, connections := AddAgentAndConnections(s.db, "TestGetMessageCount", 3)
			size := 5
			fake.AddMessages(s.db, a.ID, connections[0].ID, size)

			// Get count
			got, err := s.db.GetMessageCount(a.ID, nil)
			if err != nil {
				t.Errorf("Error fetching count %s", err.Error())
			} else if got != size {
				t.Errorf("Mismatch in fetched message  count expected: %v  got: %v", size, got)
			}
		})
	}
}

func TestGetConnectionMessageCount(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("get connection message s count "+s.name, func(t *testing.T) {
			// add new agent with no pre-existing message s
			a, connections := AddAgentAndConnections(s.db, "TestGetConnectionMessageCount", 3)
			size := 5
			index := 0
			fake.AddMessages(s.db, a.ID, connections[index].ID, (index+1)*size)
			index++
			fake.AddMessages(s.db, a.ID, connections[index].ID, (index+1)*size)
			index++
			fake.AddMessages(s.db, a.ID, connections[index].ID, index*size)

			// Get count
			expected := index * size
			got, err := s.db.GetMessageCount(a.ID, &connections[index].ID)
			if err != nil {
				t.Errorf("Error fetching count %s", err.Error())
			} else if got != expected {
				t.Errorf("Mismatch in fetched message  count expected: %v  got: %v", expected, got)
			}
		})
	}
}

func TestGetConnectionForMessage(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("get connection for message "+s.name, func(t *testing.T) {
			a, connections := AddAgentAndConnections(s.db, "TestGetConnectionForMessage", 3)
			connection := connections[0]
			messages := fake.AddMessages(s.db, a.ID, connection.ID, 1)
			message := messages[0]

			// Get data for id
			got, err := s.db.GetConnectionForMessage(message.ID, a.ID)
			if err != nil {
				t.Errorf("Error fetching connection %s", err.Error())
			} else {
				validateConnection(t, connection, got)
			}
		})
	}
}

func TestArchiveMessage(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("archive message "+s.name, func(t *testing.T) {
			testMessage = s.newTestMessage(testMessage)

			// Add data
			m, err := s.db.AddMessage(testMessage)
			assert.D.True(err == nil)

			now := utils.CurrentTime()
			err = s.db.ArchiveMessage(m.ID, m.TenantID)
			if err != nil {
				t.Errorf("Failed to archive message %s", err.Error())
			}

			// Get data for id
			got, err := s.db.GetMessage(m.ID, m.TenantID)
			assert.D.True(err == nil)

			m.Archived = &now
			validateMessage(t, m, got)
		})
	}
}
