package faker

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/findy-network/findy-agent-vault/tools/utils"

	"github.com/golang/glog"

	"github.com/findy-network/findy-agent-vault/tools/data/model"

	"github.com/bxcodec/faker/v3"
	"github.com/lainio/err2"
)

const (
	eventsCountFactor = 10
	msgsCountFactor   = 5
)

func initFaker(c *model.Items) {
	defer err2.Catch(func(err error) {
		panic(err)
	})

	err2.Check(faker.AddProvider("organisationLabel", func(v reflect.Value) (interface{}, error) {
		orgs := []string{"Bank", "Ltd", "Agency", "Company", "United"}
		index := utils.Random(len(orgs))
		f := fakeLastName{}
		_ = faker.FakeData(&f)
		return f.Name + " " + orgs[index], nil
	}))

	err2.Check(faker.AddProvider("pairwiseIdPtr", func(v reflect.Value) (interface{}, error) {
		id := c.RandomID()
		return id, nil
	}))

	err2.Check(faker.AddProvider("pairwiseId", func(v reflect.Value) (interface{}, error) {
		id := c.RandomID()
		return *id, nil
	}))
}

func printObject(objectPtr, object interface{}, printComma bool) {
	t := reflect.TypeOf(object)
	s := reflect.ValueOf(objectPtr).Elem()
	fmt.Printf("{\n")
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		if !strings.HasPrefix(t.Field(i).Name, "Skip") {
			if i != 0 {
				fmt.Printf("\n")
			}
			fmt.Printf("\t\t")
			if f.Type().String() == "string" {
				fmt.Printf("\"%s\"", f.Interface())
			} else if f.Type().String() == "int64" {
				fmt.Printf("%d", f.Interface())
			} else if f.Type().String() == "bool" {
				fmt.Printf("%t", f.Interface())
			} else {
				fmt.Printf("%s", f.Interface())
			}
			fmt.Printf(",")
		}
	}

	fmt.Print("\n\t}")
	if printComma {
		fmt.Print(",")
	}
	fmt.Print("\n")
}

func Run(c, e, m *model.Items) *model.InternalUser {
	defer err2.Catch(func(err error) {
		panic(err)
	})
	initFaker(c)

	connCount := 5
	conns, err := FakeConnections(connCount, true)
	err2.Check(err)
	for index := range conns {
		c.Append(&conns[index])
	}

	msgs, err := FakeMessages(connCount * msgsCountFactor)
	err2.Check(err)
	for index := range msgs {
		m.Append(&msgs[index])
	}

	events, err := fakeAndPrintEvents(connCount*eventsCountFactor, true)
	err2.Check(err)
	for index := range events {
		e.Append(&events[index])
	}

	user, err := fakeUser(true)
	err2.Check(err)

	glog.Infof("Generated %d connections, %d messages and %d events for user %s", len(conns), len(msgs), len(events), user.Name)
	return &user
}

func FakeMessages(count int) (msgs []model.InternalMessage, err error) {
	defer err2.Return(&err)
	msgs = make([]model.InternalMessage, count)

	for i := 0; i < count; i++ {
		m := model.InternalMessage{}
		err2.Check(faker.FakeData(&m))
		msgs[i] = m
	}
	sort.Slice(msgs, func(i, j int) bool {
		return msgs[i].CreatedMs < msgs[j].CreatedMs
	})
	return
}
