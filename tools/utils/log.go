package utils

import (
	"flag"
	"fmt"

	"github.com/lainio/err2"
)

func SetLogDefaults() {
	defer err2.Catch(func(err error) {
		fmt.Println("ERROR:", err)
	})
	err2.Check(flag.Set("logtostderr", "true"))
	err2.Check(flag.Set("stderrthreshold", "WARNING"))
	err2.Check(flag.Set("v", "3"))
	flag.Parse()
}
