package pg

import (
	"testing"
	"time"

	"github.com/findy-network/findy-agent-vault/db/model"
)

func TestAddAgent(t *testing.T) {
	testAgent := &model.Agent{AgentID: "agentID", Label: "agentLabel"}

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
	a, err := pgDB.AddAgent(testAgent)
	if err != nil {
		t.Errorf("Failed to add agent %s", err.Error())
	} else {
		validateAgent(a)
	}

	// Error with duplicate id
	if _, err := pgDB.AddAgent(testAgent); err == nil {
		t.Errorf("Expecting duplicate key error")
	} else {
		if pgErr, ok := err.(*PgError); ok {
			if pgErr.code != PgErrorUniqueViolation {
				t.Errorf("Expecting duplicate key error %s", pgErr.code)
			}
		} else {
			t.Errorf("Expecting pg error %v", err)
		}
	}

	// Get data for id
	a1, err := pgDB.GetAgent(&a.ID, nil)
	if err != nil {
		t.Errorf("Error fetching agent %s", err.Error())
	} else {
		validateAgent(a1)
	}

	// Get data for agent id
	a2, err := pgDB.GetAgent(nil, &a.AgentID)
	if err != nil {
		t.Errorf("Error fetching agent %s", err.Error())
	} else {
		validateAgent(a2)
	}

}
