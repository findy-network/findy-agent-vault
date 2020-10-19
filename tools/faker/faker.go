package faker

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/lainio/err2"
)

func printObject(objectPtr interface{}, object interface{}, printComma bool) {
	t := reflect.TypeOf(object)
	s := reflect.ValueOf(objectPtr).Elem()
	fmt.Printf("{")
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		if !strings.HasPrefix(t.Field(i).Name, "Skip") {

			if i != 0 {
				fmt.Printf(",")
			}
			if f.Type().String() == "string" {
				fmt.Printf("\"%s\"", f.Interface())
			} else if f.Type().String() == "int64" {
				fmt.Printf("%d", f.Interface())
			} else if f.Type().String() == "bool" {
				fmt.Printf("%t", f.Interface())
			} else {
				fmt.Printf("%s", f.Interface())
			}

		}
	}

	fmt.Print("}")
	if printComma {
		fmt.Print(",")
	}
	fmt.Print("\n")
}

func Run() {
	defer err2.Catch(func(err error) {
		fmt.Println("ERROR:", err)
	})

	connCount := 5

	conns, err := fakeConnections(connCount)
	err2.Check(err)

	fakeAndPrintEvents(connCount*10, conns)

	_, _ = fakeUser()
}
