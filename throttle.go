package main

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/sync/semaphore"
)

func newThrottledChannelUserTokenMap() *throttledChannelUserTokenMap {
	return &throttledChannelUserTokenMap{
		tokens: map[string]*semaphore.Weighted{},
	}
}

// throttledChannelUserTokenMap represents a map of users + throttled channels and the number of tokens they have
// remaining.
type throttledChannelUserTokenMap struct {
	sync.RWMutex
	tokens map[string]*semaphore.Weighted
}

// userCanPost tests if a user is able to post given the max amount of tokens a user has.
func (c *throttledChannelUserTokenMap) userCanPost(key string, maxTokens int, t time.Duration) bool {
	c.RLock()
	tokenCount := c.tokens[key]
	c.RUnlock()
	if tokenCount == nil {
		tokenCount = c.initUserPostTokens(key, maxTokens)
	}
	value := tokenCount.TryAcquire(1)
	if value {
		time.AfterFunc(t, func() {
			tokenCount.Release(1)
		})
	}
	return value
}

// initUserPostTokens adds a new semaphore to the map with the correct value and returns said value.
func (c *throttledChannelUserTokenMap) initUserPostTokens(key string, maxTokens int) *semaphore.Weighted {
	c.Lock()
	defer c.Unlock()
	tokenCount := c.tokens[key]
	if tokenCount != nil {
		return tokenCount
	}
	semaphore := semaphore.NewWeighted(int64(maxTokens))
	c.tokens[key] = semaphore
	return semaphore
}

func (b *Bot) handleThrottle(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}
	content, err := m.ContentWithMoreMentionsReplaced(s)
	if err != nil {
		log.Error(err)
		return
	}
	regex := regexp.MustCompile("\\<(.*?)\\>")
	content = string(regex.ReplaceAll([]byte(content), []byte("")))
	for _, c := range b.config.Throttle {
		if m.ChannelID == c.ChannelID {
			message := ""
			if c.CharLimit > 0 && c.CharLimit < len(content) {
				message = buildCharLimitResponse(c.CharLimit, len(content))

			} else if c.NewlineLimit > 0 && c.NewlineLimit < strings.Count(m.Content, "\n") {
				message = buildNewlineLimitResponse(c.NewlineLimit, strings.Count(m.Content, "\n"))

			} else if !b.throttledChannels.userCanPost(m.Author.ID+m.ChannelID, c.MaxTokens, time.Duration(c.TokenInterval)*time.Second) {
				message = "Hello " + getUserName(m.Author, m.Member) + ", this is an automated message. Your latest message on ERC has been deleted due to our 24-hour repost rule. \n\nYou will be able to post again in the LFG/LFM channel in " + (time.Duration(c.TokenInterval) * time.Second).String() + ". \n\nIf you have deleted your message by accident, the admin team cannot help you lift the limit.\n\nIf your previous post was removed by a moderator for breaking the rules then you will have to wait for the 24 hours to pass. Please do not repost on another account to get around this. For more information, please read #announcements and the pins in the relevant channels."
			}
			if message != "" {
				err := s.ChannelMessageDelete(m.ChannelID, m.ID)
				if err != nil {
					b.logMessageDeleteError(c.ChannelID, m.ID, err)
				}
				b.pmUser(m.Author.ID, message)
				b.logThrottleUser(m)
			}
			break
		}
	}
}

func (b *Bot) pmUser(userID, content string) {
	c, err := b.session.UserChannelCreate(userID)
	if err != nil {
		b.logFailedToCreateChannel(userID, err)
		return
	}
	_, err = b.session.ChannelMessageSend(c.ID, content)
	if err != nil {
		b.logMessageSendError(c.ID, err)
	}
}

func buildCharLimitResponse(i, j int) string {
	return fmt.Sprintf("Your message was deleted because it is too long, the limit is %d characters "+
		"while your message is %d characters long.", i, j)
}

func buildNewlineLimitResponse(i, j int) string {
	return fmt.Sprintf("Your message was deleted because it has too many newlines, the limit is %d "+
		"while your message has %d newlines.", i, j)
}
