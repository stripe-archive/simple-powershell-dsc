package main

import (
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sirupsen/logrus"

	"github.com/stripe-archive/simple-powershell-dsc/cmd/lambda/internal/gateway"
	"github.com/stripe-archive/simple-powershell-dsc/dsc"
	s3config "github.com/stripe-archive/simple-powershell-dsc/dsc/config/s3"
	s3report "github.com/stripe-archive/simple-powershell-dsc/dsc/report/s3"
	s3status "github.com/stripe-archive/simple-powershell-dsc/dsc/status/s3"
)

func main() {
	log := logrus.New()
	log.Level = logrus.DebugLevel

	// Pick S3 buckets from environment
	var (
		configBucket  string
		reportsBucket string
		statusBucket  string
	)
	environMap := []struct {
		Out     *string
		Environ string
	}{
		{&configBucket, "CONFIG_S3_BUCKET"},
		{&reportsBucket, "REPORTS_S3_BUCKET"},
		{&statusBucket, "STATUS_S3_BUCKET"},
	}
	for _, e := range environMap {
		val, found := os.LookupEnv(e.Environ)
		if !found {
			log.WithField("var", e.Environ).Fatal("required environment variable not found")
		}
		*e.Out = val
	}

	sess, err := session.NewSession()
	if err != nil {
		log.WithError(err).Fatal("error creating AWS session")
	}
	s3api := s3.New(sess)

	config := s3config.New(configBucket, s3api)
	report := s3report.New(reportsBucket, s3api)
	status := s3status.New(config, statusBucket, s3api)

	mgr := dsc.NewManager(config, report, status, dsc.WithLogger(log))

	log.Info("server started")
	if err := gateway.ListenAndServe(":3000", mgr); err != nil {
		log.WithError(err).Fatal("error in server")
	}
}
