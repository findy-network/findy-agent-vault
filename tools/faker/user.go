package faker

import (
	"fmt"

	"github.com/findy-network/findy-agent-vault/tools/data/model"

	"github.com/bxcodec/faker/v3"
	"github.com/lainio/err2"
)

func fakeUser(skipPrint bool) (user model.InternalUser, err error) {
	defer err2.Return(&err)

	err2.Check(faker.FakeData(&user))
	if !skipPrint {
		fmt.Printf("var user = InternalUser")
		printObject(&user, user, false)
	}

	return
}
