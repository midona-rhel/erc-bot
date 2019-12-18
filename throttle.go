package main

import (
	"github.com/bwmarrin/discordgo"

	"sync"
	"time"
)

func newUserDebtMap() *userDebtMap {
	return &userDebtMap{
		tokens: map[string]int{},
	}
}

// userDebtMap represents a map of users and the number of tokens they have remaining.
type userDebtMap struct {
	sync.RWMutex
	tokens map[string]int
}

func (c *userDebtMap) get(s string) int {
	c.RLock()
	defer c.RUnlock()
	return c.tokens[s]
}

func (c *userDebtMap) set(s string, i int) {
	c.Lock()
	c.tokens[s] = i
	c.Unlock()
}

func newThrottledChannelsMap() *throttledChannelsMap {
	return &throttledChannelsMap{
		channels: map[string]*userDebtMap{},
	}
}

// throttledChannelsMap represents a map of the channels to user message timestamps.
type throttledChannelsMap struct {
	sync.RWMutex
	channels map[string]*userDebtMap
}

func (u *throttledChannelsMap) get(s string) *userDebtMap {
	u.RLock()
	defer u.RUnlock()
	return u.channels[s]
}

func (u *throttledChannelsMap) set(s string, t *userDebtMap) {
	u.Lock()
	u.channels[s] = t
	u.Unlock()
}
func (b *Bot) handleThrottle(m *discordgo.MessageCreate, s *discordgo.Session) error {
	for _, c := range b.config.Throttle {
		if m.ChannelID == c.ChannelID {
			debt := b.throttledChannels.get(c.ChannelID)
			if debt == nil {
				debt = newUserDebtMap()
				debt.set(m.Author.ID, 0)
				go removeDebt(time.Duration(c.TokenInterval)*time.Second, debt, m.Author.ID)
			} else {
				i := debt.get(m.Author.ID)
				if i == c.MaxDebt {
					return s.ChannelMessageDelete(m.ChannelID, m.ID)
				}
				debt.set(m.Author.ID, i+1)
			}
		}
	}
	return nil
}

func removeDebt(t time.Duration, u *userDebtMap, userID string) {
	for {
		<-time.After(t)
		debt := u.get(userID)
		if debt != 0 {
			u.set(userID, debt-1)
		}
	}
}
