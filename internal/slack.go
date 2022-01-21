package sf

import (
	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

type SlackHandler struct {
	messages []slack.Attachment
}

func (s *SlackHandler) AddSlackMessage(message, color string) {
	newMessage := slack.Attachment{
		Text:  message,
		Color: color,
	}
	s.messages = append(s.messages, newMessage)
}

func (s *SlackHandler) SendSlackMessage(config Configuration) {
	if config.Notifications.Slack.WebhookURL == "" || len(s.messages) == 0 {
		return
	}

	msg := slack.WebhookMessage{
		Username:    config.Notifications.Slack.Username,
		IconURL:     config.Notifications.Slack.IconURL,
		Channel:     config.Notifications.Slack.Channel,
		Attachments: s.messages,
	}

	err := slack.PostWebhook(config.Notifications.Slack.WebhookURL, &msg)
	if err != nil {
		log.Errorf("Error sending slack message: %s", err)
	}
}

// NewSlackHandler creates a new slack handler
func NewSlackHandler() *SlackHandler {
	return &SlackHandler{}
}
