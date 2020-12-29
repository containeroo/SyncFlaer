package sf

import (
	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

var slackMessages []slack.Attachment

func addSlackMessage(message string, color string) {
	newMessage := slack.Attachment{
		Color: color,
		Text:  message,
	}
	slackMessages = append(slackMessages, newMessage)
}

// SendSlackMessage sends a Slack message if configured
func SendSlackMessage() {
	if config.Notifications.Slack.WebhookURL == "" || len(slackMessages) == 0 {
		return
	}

	msg := slack.WebhookMessage{
		Username:    config.Notifications.Slack.Username,
		IconURL:     config.Notifications.Slack.IconURL,
		Channel:     config.Notifications.Slack.Channel,
		Attachments: slackMessages,
	}

	err := slack.PostWebhook(config.Notifications.Slack.WebhookURL, &msg)
	if err != nil {
		log.Errorf("Unable to send Slack message: %s", err)
	}
}
