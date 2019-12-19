package main

import (
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
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
	for _, c := range b.config.Throttle {
		if m.ChannelID == c.ChannelID {
			if b.throttledChannels.userCanPost(m.Author.ID+m.ChannelID, c.MaxTokens, time.Duration(c.TokenInterval)*time.Second) {
			}
			log.WithFields(logrus.Fields{
				"userID":    m.Author.ID,
				"channelID": c.ChannelID,
			}).Info("throttled user")
			err := s.ChannelMessageDelete(m.ChannelID, m.ID)
			if err != nil {
				log.WithFields(logrus.Fields{
					"userID":    m.Author.ID,
					"channelID": c.ChannelID,
					"messageID": m.ID,
					"action":    "failed to delete user post",
				}).Error(err)
			}
		}
	}
}
