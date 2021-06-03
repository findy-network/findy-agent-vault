package utils

import (
	"flag"
	"fmt"

	"github.com/golang/glog"
	"github.com/lainio/err2"
)

const (
	logLevelHigh   = 1
	logLevelMedium = 3
	logLevelLow    = 5
	logLevelTrace  = 7
)

func LogHigh() glog.Verbose  { return glog.V(logLevelHigh) }
func LogMed() glog.Verbose   { return glog.V(logLevelMedium) }
func LogLow() glog.Verbose   { return glog.V(logLevelLow) }
func LogTrace() glog.Verbose { return glog.V(logLevelTrace) }

func SetLogDefaults() {
	logParse("3")
}

func SetLogConfig(config *Configuration) {
	logParse(config.LogLevel)
}

func logParse(level string) {
	defer err2.Catch(func(err error) {
		fmt.Println("ERROR:", err)
	})
	err2.Check(flag.Set("logtostderr", "true"))
	err2.Check(flag.Set("stderrthreshold", "WARNING"))
	err2.Check(flag.Set("v", level))
	flag.Parse()
}
