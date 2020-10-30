// +build findy

package agency

import (
	"testing"
)

func TestInvite(t *testing.T) {
	invitation, err := Instance.Invite()
	if err != nil || len(invitation) == 0 {
		t.Errorf("Invitation failed = err (%v), invitation (%s)", err, invitation)
	}
	t.Log(invitation)
}
