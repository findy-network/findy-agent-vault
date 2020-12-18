package utils

import (
	"flag"
	"fmt"

	"github.com/golang/glog"
	"github.com/lainio/err2"
)

const (
	logLevelMedium = 2
	logLevelLow    = 3
)

func LogMed() glog.Verbose { return glog.V(logLevelMedium) }
func LogLow() glog.Verbose { return glog.V(logLevelLow) }

func SetLogDefaults() {
	defer err2.Catch(func(err error) {
		fmt.Println("ERROR:", err)
	})
	err2.Check(flag.Set("logtostderr", "true"))
	err2.Check(flag.Set("stderrthreshold", "WARNING"))
	err2.Check(flag.Set("v", "3"))
	flag.Parse()
}
