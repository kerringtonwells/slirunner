package commands

import (
	"context"
	"fmt"
	"os"
  "strings"
	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/lager/lagerctx"
	"github.com/kerringtonwells/slirunner/probes"
)

type onceCommand struct {
	Target          string `long:"target" required:"true"`
	PipelinesPrefix string `long:"prefix" default:"slirunner-"`

	Username     string `long:"username"      short:"u" required:"true"`
	Password     string `long:"password"      short:"p" required:"true"`
	ConcourseUrl string `long:"concourse-url" short:"c" required:"true"`
	InsecureTls  bool   `long:"insecure-tls"  short:"k" required:"false" description:"Skip tls verification"`

	LdapAuth bool   `long:"ldapauth"      short:"l" required:"false" description:"LDAP boolean if using ldap auth"`
	LdapTeam string `long:"ldapteam"      short:"m" required:"false" description:"LDAP team if using ldap auth"`
	WorkerPool string `long:"workerpool"      short:"w" required:"true" description:"worker pool for concourse pipelines"`
	Harbor_url string `long:"harbor_url"      short:"r" required:"true" description:"repository url"`
	Debug string `long:"debug"      short:"d" required:"true" description:"debug"`
}

func (c *onceCommand) Execute(args []string) (err error) {
	var logs string
	if strings.Contains(c.Debug, "true") {
	    logs = "set -o xtrace"
	}else {
      logs = ""
	}

  //singleConcourseUrl := strings.Split(c.ConcourseUrl, " ")
	//singleTarget := strings.Split(c.Target, " ")
	//singleWorkerPool := strings.Split(c.WorkerPool, " ")
	//counter := 0
	//var stringCounter string
  //for i := range singleConcourseUrl {
	ctx, cancel := context.WithCancel(context.Background())
	go onTerminationSignal(cancel)
	logger := lager.NewLogger("my-app")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.INFO))

	ctx = lagerctx.NewContext(ctx, logger)
		err = probes.NewAll(
			c.Target,
			c.Username, c.Password,
			c.ConcourseUrl,
			c.PipelinesPrefix,
			c.InsecureTls,
			c.LdapAuth, c.LdapTeam,
			c.WorkerPool,
			c.Harbor_url,
			logs,
			//logs,
		).Run(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
    }
	return
}
