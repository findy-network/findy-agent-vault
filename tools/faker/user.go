package faker

import (
	"fmt"

	"github.com/bxcodec/faker/v3"
	"github.com/findy-network/findy-agent-vault/tools/data"
	"github.com/lainio/err2"
)

func fakeUser() (user data.InternalUser, err error) {
	defer err2.Return(&err)

	err2.Check(faker.FakeData(&user))
	fmt.Printf("var user = InternalUser")
	printObject(&user, user, false)

	return
}
