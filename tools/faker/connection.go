package faker

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/findy-network/findy-agent-api/tools/utils"

	"github.com/bxcodec/faker/v3"
	"github.com/findy-network/findy-agent-api/tools/data"
	"github.com/lainio/err2"
)

type fakeLastName struct {
	Name string `faker:"last_name"`
}

func fakeConnections(count int) (conns []data.InternalPairwise, err error) {
	defer err2.Return(&err)
	conns = make([]data.InternalPairwise, count)
	err = faker.AddProvider("organisationLabel", func(v reflect.Value) (interface{}, error) {
		orgs := []string{"Bank", "Ltd", "Agency", "Company", "United"}
		index := utils.Random(len(orgs))
		f := fakeLastName{}
		_ = faker.FakeData(&f)
		return f.Name + " " + orgs[index], nil
	})

	for i := 0; i < count; i++ {
		conn := data.InternalPairwise{}
		err2.Check(faker.FakeData(&conn))
		conns[i] = conn
	}
	sort.Slice(conns, func(i, j int) bool {
		return conns[i].CreatedMs < conns[j].CreatedMs
	})
	fmt.Println("var connections = []InternalPairwise{")
	for i := 0; i < len(conns); i++ {
		fmt.Printf("	")
		printObject(&(conns)[i], (conns)[i], true)
	}
	fmt.Println("}")
	return
}
