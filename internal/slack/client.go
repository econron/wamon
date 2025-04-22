package slack

import (
	"fmt"
	"strings"
	"time"

	"github.com/econron/wamon/internal/models"
	"github.com/slack-go/slack"
)

// Config holds the Slack configuration values
type Config struct {
	Token   string
	Channel string
	Enabled bool
}

// Client is a wrapper around the Slack API client
type Client struct {
	api    *slack.Client
	config Config
}

// NewClient creates a new Slack client
func NewClient(config Config) *Client {
	// Only check for token existence, not enabled flag
	if config.Token == "" {
		return &Client{config: config}
	}

	return &Client{
		api:    slack.New(config.Token),
		config: config,
	}
}

// SendWeeklyReport sends a summary of the last week's entries to the configured Slack channel
func (c *Client) SendWeeklyReport(entries []*models.Entry) error {
	// Only check if the api is initialized
	if c.api == nil {
		return fmt.Errorf("slack client is not properly initialized")
	}

	if len(entries) == 0 {
		return fmt.Errorf("no entries to report")
	}

	// Create the message blocks
	headerBlock := slack.NewHeaderBlock(slack.NewTextBlockObject("plain_text", "🦭 先週のワモンアザラシの記録 🦭", true, false))

	var blocks []slack.Block
	blocks = append(blocks, headerBlock)
	blocks = append(blocks, slack.NewDividerBlock())

	// Group entries by day
	entriesByDay := make(map[string][]*models.Entry)
	for _, entry := range entries {
		day := entry.CreatedAt.Format("2006-01-02")
		entriesByDay[day] = append(entriesByDay[day], entry)
	}

	// Add each day's entries
	for day, dayEntries := range entriesByDay {
		// Format the date as a section header
		dateTime, _ := time.Parse("2006-01-02", day)
		dateHeader := slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*%s*", dateTime.Format("2006/01/02 (Mon)")), false, false),
			nil, nil,
		)
		blocks = append(blocks, dateHeader)

		// Add each entry for this day
		for _, entry := range dayEntries {
			var fieldTexts []string

			fieldTexts = append(fieldTexts, fmt.Sprintf("*カテゴリ:* %s", entry.Category))

			if entry.ResearchTopic != "" {
				fieldTexts = append(fieldTexts, fmt.Sprintf("*調べたこと:* %s", entry.ResearchTopic))
			}

			if entry.ProgramTitle != "" {
				fieldTexts = append(fieldTexts, fmt.Sprintf("*書いたプログラム:* %s", entry.ProgramTitle))
			}

			fieldTexts = append(fieldTexts, fmt.Sprintf("*満足度:* %s", getStarRating(entry.Satisfaction)))

			// Create a section block for the entry
			entryText := strings.Join(fieldTexts, "\n")
			entryBlock := slack.NewSectionBlock(
				slack.NewTextBlockObject("mrkdwn", entryText, false, false),
				nil, nil,
			)
			blocks = append(blocks, entryBlock)
			blocks = append(blocks, slack.NewDividerBlock())
		}
	}

	// Add summary footer
	summaryText := fmt.Sprintf("先週は合計 *%d件* の記録がありました。次も頑張りましょう！", len(entries))
	summaryBlock := slack.NewSectionBlock(
		slack.NewTextBlockObject("mrkdwn", summaryText, false, false),
		nil, nil,
	)
	blocks = append(blocks, summaryBlock)

	// Send the message
	_, _, err := c.api.PostMessage(
		c.config.Channel,
		slack.MsgOptionBlocks(blocks...),
		slack.MsgOptionAsUser(true),
	)

	return err
}

// getStarRating returns a star rating representation of the satisfaction level
func getStarRating(satisfaction int) string {
	if satisfaction < 1 {
		satisfaction = 1
	} else if satisfaction > 5 {
		satisfaction = 5
	}

	stars := strings.Repeat("★", satisfaction)
	emptyStars := strings.Repeat("☆", 5-satisfaction)
	return stars + emptyStars
}
