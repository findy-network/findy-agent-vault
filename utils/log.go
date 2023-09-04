package utils

import (
	"flag"
	"log"

	"github.com/golang/glog"
	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
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
	defer err2.Catch(err2.Err(func(err error) {
		log.Println("ERROR:", err)
	}))
	try.To(flag.Set("logtostderr", "true"))
	try.To(flag.Set("stderrthreshold", "WARNING"))
	try.To(flag.Set("v", level))
	flag.Parse()
}
