package test

import (
	"testing"
	"time"

	"github.com/findy-network/findy-agent-vault/db/model"
)

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
