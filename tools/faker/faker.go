package faker

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/bxcodec/faker/v3"
	"github.com/findy-network/findy-agent-api/tools/data"

	"github.com/lainio/err2"
)

const (
	eventsCountFactor = 10
)

func InitFaker() {
	_ = faker.AddProvider("eventPairwiseId", func(v reflect.Value) (interface{}, error) {
		return data.State.Connections.RandomID(), nil
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

func Run() {
	defer err2.Catch(func(err error) {
		fmt.Println("ERROR:", err)
	})

	InitFaker()

	connCount := 5

	conns, err := fakeConnections(connCount)
	err2.Check(err)

	fakeAndPrintEvents(connCount*eventsCountFactor, conns)

	_, _ = fakeUser()
}
