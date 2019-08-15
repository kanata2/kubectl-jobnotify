package jobnotify

import (
	"fmt"
	"log"
	"time"

	slack "github.com/monochromegane/slack-incoming-webhooks"
)

const defaultLayout = "2006/01/02 15:04:05"

type Notifier interface {
	Notify(job, namespace string) error
}

type Stdout struct{}

func (s *Stdout) Notify(job, namespace string) error {
	result, err := watch(job, namespace)
	if err != nil {
		return err
	}
	log.Printf(
		"%s complete. start: %s, end: %s, succeeded: %d, failed: %d\n",
		result.name,
		result.startedAt.Format(defaultLayout),
		result.completedAt.Format(defaultLayout),
		result.completed,
		result.failed,
	)
	return nil
}

type Slack struct {
	webhookURL string
}

func NewSlack(webhookURL string) *Slack {
	return &Slack{
		webhookURL: webhookURL,
	}
}

func (s *Slack) Valid() bool {
	return s.webhookURL != ""
}

func (s *Slack) Notify(job, namespace string) error {
	result, err := watch(job, namespace)
	if err != nil {
		return err
	}
	client := slack.Client{WebhookURL: s.webhookURL}
	return client.Post(&slack.Payload{
		Text: "Kubernetes Job Result",
		Attachments: []*slack.Attachment{
			&slack.Attachment{
				Title: result.name,
				Fields: []*slack.Field{
					&slack.Field{
						Title: "start datetime",
						Value: result.startedAt.Format(defaultLayout),
						Short: true,
					},
					&slack.Field{
						Title: "complete datetime",
						Value: result.completedAt.Format(defaultLayout),
						Short: true,
					},
					&slack.Field{
						Title: "success",
						Value: fmt.Sprint(result.completed),
						Short: true,
					},
					&slack.Field{
						Title: "failed",
						Value: fmt.Sprint(result.failed),
						Short: true,
					},
				},
				Timestamp: time.Now().Unix(),
			},
		},
	})
	return nil
}
