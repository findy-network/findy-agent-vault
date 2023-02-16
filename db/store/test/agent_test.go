package test

import (
	"testing"
	"time"

	"github.com/findy-network/findy-agent-vault/db/model"
	"github.com/findy-network/findy-agent-vault/paginator"
	"github.com/lainio/err2/assert"
)

func TestGetListenerAgents(t *testing.T) {
	token := "jwt-token"
	testAgent1 := &model.Agent{}
	testAgent1.AgentID = "agentID1"
	testAgent1.Label = "agentLabel1"
	testAgent1.RawJWT = token

	testAgent2 := &model.Agent{}
	testAgent2.AgentID = "agentID2"
	testAgent2.Label = "agentLabel2"
	testAgent2.RawJWT = ""

	for index := range DBs {
		s := DBs[index]
		t.Run("get listener agents "+s.name, func(t *testing.T) {
			assert.PushTester(t)
			defer assert.PopTester()

			agent1, err := s.db.AddAgent(testAgent1)
			assert.NoError(err, "Failed to add agent %v", err)

			agent2, err := s.db.AddAgent(testAgent2)
			assert.NoError(err, "Failed to add agent %v", err)

			pagination := paginator.BatchInfo{Count: 5}
			agents, err := s.db.GetListenerAgents(&pagination)

			assert.NoError(err, "Failed to get listener agents %v", err)

			assert.That(len(agents.Agents) >= 1, "Agents count mismatch got %v, expected at least 1", len(agents.Agents))

			foundWithToken := false
			foundWithoutToken := false
			for _, a := range agents.Agents {
				if a.ID == agent1.ID {
					foundWithToken = true
				}
				if a.ID == agent2.ID {
					foundWithoutToken = true
				}
			}
			assert.That(foundWithToken, "Did not receive listener agent")
			assert.ThatNot(foundWithoutToken, "Received listener agent without token")
		})
	}
}

func agentValidator(t *testing.T, testAgent *model.Agent) func(a *model.Agent) {
	return func(agent *model.Agent) {
		if agent == nil {
			t.Errorf("Expecting result, agent is nil")
			return
		}
		if agent.AgentID != testAgent.AgentID {
			t.Errorf("Agent id mismatch expected %s got %s", testAgent.AgentID, agent.AgentID)
		}
		if agent.Label != testAgent.Label {
			t.Errorf("Agent label mismatch expected %s got %s", testAgent.Label, agent.Label)
		}
		if agent.RawJWT != testAgent.RawJWT {
			t.Errorf("Agent RawJWT mismatch expected %v got %v", testAgent.RawJWT, agent.RawJWT)
		}
		if agent.ID == "" {
			t.Errorf("Invalid agent id %s", agent.ID)
		}
		if time.Since(agent.Created) > time.Second {
			t.Errorf("Timestamp not in threshold %v", agent.Created)
		}
		if agent.Cursor == 0 {
			t.Errorf("Cursor invalid %v", agent.Cursor)
		}
	}
}

func TestAddAgent(t *testing.T) {
	for index := range DBs {
		store := DBs[index]
		t.Run("add agent "+store.name, func(t *testing.T) {
			testAgent := &model.Agent{}
			testAgent.AgentID = "agentID"
			testAgent.Label = "agentLabel"
			testJwt := "jwt"
			testAgent.RawJWT = testJwt

			validateAgent := agentValidator(t, testAgent)

			// Add data
			agent, err := store.db.AddAgent(testAgent)
			if err != nil {
				t.Errorf("Failed to add agent %s", err.Error())
			} else {
				validateAgent(agent)
			}
			time.Sleep(time.Nanosecond)

			// Get data for id
			a1, err := store.db.GetAgent(&agent.ID, nil)
			if err != nil {
				t.Errorf("Error fetching agent %s", err.Error())
			} else {
				validateAgent(a1)
			}

			// Get data for agent id
			a2, err := store.db.GetAgent(nil, &agent.AgentID)
			if err != nil {
				t.Errorf("Error fetching agent %s", err.Error())
			} else {
				validateAgent(a2)
			}

			var updatedAgent *model.Agent
			newJwt := "new jwt"
			agent.RawJWT = newJwt
			if updatedAgent, err = store.db.AddAgent(agent); err != nil {
				t.Errorf("Failed to update agent %s", err.Error())
			} else if err == nil {
				if updatedAgent.LastAccessed.Sub(agent.LastAccessed) == 0 {
					t.Errorf("Timestamp not updated %v from %v", updatedAgent.LastAccessed, agent.LastAccessed)
				}
				if newJwt != updatedAgent.RawJWT {
					t.Errorf("Token not updated %v expected %v", updatedAgent.RawJWT, newJwt)
				}
			}
		})
	}
}
