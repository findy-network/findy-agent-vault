package faker

import (
	"fmt"
	"sort"

	"github.com/bxcodec/faker/v3"
	data "github.com/findy-network/findy-agent-vault/tools/data/model"
	"github.com/lainio/err2"
)

func FakeConnections(count int, skipPrint bool) (conns []data.InternalPairwise, err error) {
	defer err2.Return(&err)
	conns = make([]data.InternalPairwise, count)

	for i := 0; i < count; i++ {
		conn := data.InternalPairwise{}
		err2.Check(faker.FakeData(&conn))
		conns[i] = conn
	}
	sort.Slice(conns, func(i, j int) bool {
		return conns[i].CreatedMs < conns[j].CreatedMs
	})
	if !skipPrint {
		fmt.Println("var connections = []InternalPairwise{")
		for i := 0; i < len(conns); i++ {
			fmt.Printf("	")
			printObject(&(conns)[i], (conns)[i], true)
		}
		fmt.Println("}")
	}
	return
}
