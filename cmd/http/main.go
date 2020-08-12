package main

import (
	"flag"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/stripe-archive/simple-powershell-dsc/dsc"
	localconfig "github.com/stripe-archive/simple-powershell-dsc/dsc/config/local"
	localreport "github.com/stripe-archive/simple-powershell-dsc/dsc/report/local"
	localstatus "github.com/stripe-archive/simple-powershell-dsc/dsc/status/local"
)

var (
	listenAddress string
)

func init() {
	flag.StringVar(&listenAddress, "addr", "localhost:8000", "listen address for the server")
}

func main() {
	flag.Parse()

	log := logrus.New()
	log.Level = logrus.DebugLevel

	config := localconfig.New("test/config")
	report := localreport.New("test/reports")
	status, err := localstatus.New(config, "test/status")
	if err != nil {
		log.WithError(err).Fatal("error creating NodeStatus")
	}

	mgr := dsc.NewManager(config, report, status, dsc.WithLogger(log))
	log.WithField("address", listenAddress).Info("server started")
	if err := http.ListenAndServe(listenAddress, mgr); err != nil {
		log.WithError(err).Fatal("error in server")
	}
}
