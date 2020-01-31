package main

import (
	"fmt"
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
	for _, c := range b.config.Throttle {
		if m.ChannelID == c.ChannelID {
			message := ""
			if c.CharLimit > 0 && c.CharLimit < len(m.Content) {
				message = fmt.Sprintf("Your message was deleted because it was too long, the limit is %d characters while your message was %d characters long", c.CharLimit, len(m.Content))
			} else if c.NewlineLimit > 0 && c.NewlineLimit < strings.Count(m.Content, "\n") {
				message = fmt.Sprintf("Your message was deleted because it had too many newlines, the limit is %d while your message had %d newlines", c.NewlineLimit, strings.Count(m.Content, "\n"))
			} else if !b.throttledChannels.userCanPost(m.Author.ID+m.ChannelID, c.MaxTokens, time.Duration(c.TokenInterval)*time.Second) {
				message = "Your message was deleted because you are posting too soon in the channel again."
			} else {
				return
			}
			err := s.ChannelMessageDelete(m.ChannelID, m.ID)
			if err != nil {
				b.logMessageDeleteError(c.ChannelID, m.ID, err)
				return
			}
			b.logThrottleUser(m)
			b.respondAndDelete(m.Author.Mention()+" "+message, m.ChannelID, m.Author.ID, time.Second*30)
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
