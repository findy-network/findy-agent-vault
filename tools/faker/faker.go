package faker

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	utils "github.com/findy-network/findy-agent-vault/tools/tools"

	"github.com/golang/glog"

	"github.com/findy-network/findy-agent-vault/tools/data/model"

	"github.com/bxcodec/faker/v3"
	"github.com/lainio/err2"
)

const (
	eventsCountFactor = 10
	msgsCountFactor   = 5
	credsCountFactor  = 2
)

func initFaker(c *model.Items) {
	_ = faker.AddProvider("organisationLabel", func(v reflect.Value) (interface{}, error) {
		orgs := []string{"Bank", "Ltd", "Agency", "Company", "United"}
		index := utils.Random(len(orgs))
		return faker.LastName() + " " + orgs[index], nil
	})

	_ = faker.AddProvider("pairwiseIdPtr", func(v reflect.Value) (interface{}, error) {
		id := c.RandomID()
		return id, nil
	})

	_ = faker.AddProvider("pairwiseId", func(v reflect.Value) (interface{}, error) {
		id := c.RandomID()
		return *id, nil
	})

	_ = faker.AddProvider("created", func(v reflect.Value) (interface{}, error) {
		t := faker.UnixTime() * int64(time.Microsecond)
		return t, nil
	})
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

func Run(p, e, m, c *model.Items) *model.InternalUser {
	defer err2.Catch(func(err error) {
		panic(err)
	})
	initFaker(p)

	connCount := 5
	conns, err := FakeConnections(connCount, true)
	err2.Check(err)
	for index := range conns {
		p.Append(&conns[index])
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

	creds, err := FakeCredentials(connCount * credsCountFactor)
	err2.Check(err)
	for index := range creds {
		c.Append(&creds[index])
	}

	user, err := fakeUser(true)
	err2.Check(err)

	glog.Infof("Generated %d connections, %d messages, %d creds "+
		"and %d events for user %s", len(conns), len(msgs), len(events), len(creds), user.Name)
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

func FakeCredentials(count int) (creds []model.InternalCredential, err error) {
	defer err2.Return(&err)
	creds = make([]model.InternalCredential, count)

	for i := 0; i < count; i++ {
		c := model.InternalCredential{}
		err2.Check(faker.FakeData(&c))
		creds[i] = c
	}
	sort.Slice(creds, func(i, j int) bool {
		return creds[i].CreatedMs < creds[j].CreatedMs
	})
	return
}

func FakeProofs(count int) (creds []model.InternalProof, err error) {
	defer err2.Return(&err)
	creds = make([]model.InternalProof, count)

	for i := 0; i < count; i++ {
		c := model.InternalProof{}
		err2.Check(faker.FakeData(&c))
		creds[i] = c
	}
	sort.Slice(creds, func(i, j int) bool {
		return creds[i].CreatedMs < creds[j].CreatedMs
	})
	return
}
