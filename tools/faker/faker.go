package faker

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/golang/glog"

	"github.com/findy-network/findy-agent-vault/tools/data/model"

	"github.com/bxcodec/faker/v3"
	"github.com/lainio/err2"
)

const (
	eventsCountFactor = 10
)

func initFaker(c *model.Items) {
	_ = faker.AddProvider("eventPairwiseId", func(v reflect.Value) (interface{}, error) {
		return c.RandomID(), nil
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

func Run(c, e *model.Items) ([]model.InternalPairwise, []model.InternalEvent, model.InternalUser) {
	defer err2.Catch(func(err error) {
		panic(err)
	})
	initFaker(c)

	connCount := 5
	conns, err := fakeConnections(connCount, true)
	err2.Check(err)

	for index := range conns {
		c.Append(&conns[index])
	}

	events, err := fakeAndPrintEvents(connCount*eventsCountFactor, true)
	err2.Check(err)

	for index := range events {
		e.Append(&events[index])
	}

	user, err := fakeUser(true)
	err2.Check(err)

	glog.Infof("Generated %d connections and %d events for user %s", len(conns), len(events), user.Name)

	return conns, events, user
}
