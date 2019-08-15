package main

import (
	"log"
	"os"

	jobnotify "github.com/kanata2/kubectl-jobnotify"
	"github.com/urfave/cli"
)

const (
	defaultNamespace = "default"

	exitCodeOK = iota
	exitCodeError
)

var notifierMap = map[string]jobnotify.Notifier{
	"stdout": &jobnotify.Stdout{},
	"slack":  jobnotify.NewSlack(os.Getenv("SLACK_WEBHOOK_URL")),
}

func main() {
	app := cli.NewApp()
	app.Name = "kubectl-jobnotify"
	app.Version = "0.0.1"
	app.Usage = "make you notify the completion of specified k8s job"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "job, j",
			Value: "",
			Usage: "job name you want to be notified",
		},
		cli.StringFlag{
			Name:  "namespace, n",
			Value: defaultNamespace,
			Usage: "namespace specified job belongs to",
		},
		cli.StringFlag{
			Name:  "destination, d",
			Value: "stdout",
			Usage: "notification destination",
		},
	}
	app.Action = execute

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func execute(c *cli.Context) error {
	job := c.String("job")
	if job == "" {
		return cli.NewExitError("job must be set", exitCodeError)
	}
	dest := c.String("destination")
	notifier, ok := notifierMap[dest]
	if !ok {
		return cli.NewExitError("destination not exist"+dest, exitCodeError)
	}
	if s, ok := notifier.(*jobnotify.Slack); ok {
		if !s.Valid() {
			log.Println("SLACK_WEBHOOK_URL does not set. So switch destination to STDOUT automatically.")
			notifier = notifierMap["stdout"]
		}
	}
	if err := notifier.Notify(job, c.String("namespace")); err != nil {
		return cli.NewExitError(err, exitCodeError)
	}
	return nil
}
