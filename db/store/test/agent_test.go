package test

import (
	"testing"
	"time"

	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/paginator"
)

func TestGetListenerAgents(t *testing.T) {
	token := "jwt-token"
	testAgent1 := model.NewAgent(nil)
	testAgent1.AgentID = "agentID1"
	testAgent1.Label = "agentLabel1"
	testAgent1.RawJWT = &token

	testAgent2 := model.NewAgent(nil)
	testAgent2.AgentID = "agentID2"
	testAgent2.Label = "agentLabel2"
	testAgent2.RawJWT = nil

	for index := range DBs {
		s := DBs[index]
		t.Run("get listener agents "+s.name, func(t *testing.T) {
			a1, err := s.db.AddAgent(testAgent1)
			if err != nil {
				t.Errorf("Failed to add agent %s", err.Error())
			}
			a2, err := s.db.AddAgent(testAgent2)
			if err != nil {
				t.Errorf("Failed to add agent %s", err.Error())
			}

			pagination := paginator.BatchInfo{Count: 5}
			agents, err := s.db.GetListenerAgents(&pagination)

			if err != nil {
				t.Errorf("Failed to get listener agents %s", err.Error())
			} else {
				if len(agents.Agents) < 1 {
					t.Errorf("Agents count mismatch got %v, expected at least 1", len(agents.Agents))
				}
				foundWithToken := false
				foundWithoutToken := false
				for _, a := range agents.Agents {
					if a.ID == a1.ID {
						foundWithToken = true
					}
					if a.ID == a2.ID {
						foundWithoutToken = true
					}
				}
				if !foundWithToken {
					t.Errorf("Did not receive listener agent")
				}
				if foundWithoutToken {
					t.Errorf("Received listener agent without token")
				}
			}
		})
	}
}

func TestAddAgent(t *testing.T) {
	for index := range DBs {
		s := DBs[index]
		t.Run("add agent "+s.name, func(t *testing.T) {
			testAgent := model.NewAgent(nil)
			testAgent.AgentID = "agentID"
			testAgent.Label = "agentLabel"

			var validateAgent = func(a *model.Agent) {
				if a == nil {
					t.Errorf("Expecting result, agent is nil")
					return
				}
				if a.AgentID != testAgent.AgentID {
					t.Errorf("Agent id mismatch expected %s got %s", testAgent.AgentID, a.AgentID)
				}
				if a.Label != testAgent.Label {
					t.Errorf("Agent label mismatch expected %s got %s", testAgent.Label, a.Label)
				}
				if a.RawJWT != testAgent.RawJWT {
					t.Errorf("Agent RawJWT mismatch expected %v got %v", testAgent.RawJWT, a.RawJWT)
				}
				if a.ID == "" {
					t.Errorf("Invalid agent id %s", a.ID)
				}
				if time.Since(a.Created) > time.Second {
					t.Errorf("Timestamp not in threshold %v", a.Created)
				}
			}

			// Add data
			a, err := s.db.AddAgent(testAgent)
			if err != nil {
				t.Errorf("Failed to add agent %s", err.Error())
			} else {
				validateAgent(a)
			}
			time.Sleep(time.Nanosecond)

			// Only update timestamp if already exists
			var updatedAgent *model.Agent
			if updatedAgent, err = s.db.AddAgent(a); err != nil {
				t.Errorf("Failed to update agent %s", err.Error())
			} else if err == nil {
				if updatedAgent.LastAccessed.Sub(a.LastAccessed) == 0 {
					t.Errorf("Timestamp not updated %v from %v", updatedAgent.LastAccessed, a.LastAccessed)
				}
			}

			// Get data for id
			a1, err := s.db.GetAgent(&a.ID, nil)
			if err != nil {
				t.Errorf("Error fetching agent %s", err.Error())
			} else {
				validateAgent(a1)
			}

			// Get data for agent id
			a2, err := s.db.GetAgent(nil, &a.AgentID)
			if err != nil {
				t.Errorf("Error fetching agent %s", err.Error())
			} else {
				validateAgent(a2)
			}
		})
	}
}
