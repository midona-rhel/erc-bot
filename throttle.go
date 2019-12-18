package main

import "github.com/bwmarrin/discordgo"

import "sync"

import "time"

import "fmt"

func newChannelTimestamps() *ChannelTimestamps {
	return &ChannelTimestamps{
		timestamps: map[string]time.Time{},
	}
}

// ChannelTimestamps represents a map of the last accepted post in a channel of a user.
type ChannelTimestamps struct {
	sync.RWMutex
	timestamps map[string]time.Time
}

func (c *ChannelTimestamps) get(s string) time.Time {
	c.RLock()
	defer c.RUnlock()
	return c.timestamps[s]
}

func (c *ChannelTimestamps) set(s string, t time.Time) {
	c.Lock()
	c.timestamps[s] = t
	c.Unlock()
}

func newUserLastPost() *UserLastPost {
	return &UserLastPost{
		users: map[string]*ChannelTimestamps{},
	}
}

// UserLastPost represents a map of the users to channeltimestamps
type UserLastPost struct {
	sync.RWMutex
	users map[string]*ChannelTimestamps
}

func (u *UserLastPost) get(s string) *ChannelTimestamps {
	u.RLock()
	defer u.RUnlock()
	return u.users[s]
}

func (u *UserLastPost) set(s string, c *ChannelTimestamps) {
	u.Lock()
	u.users[s] = c
	u.Unlock()
}

func (b *Bot) handleThrottle(m *discordgo.MessageCreate, s *discordgo.Session) {
	// if a channel is throttled.
	for _, c := range b.config.Throttle {
		if m.ChannelID == c.ChannelID {
			// If the user has an existing entry in memory check if the user has a previous post in said channel, if he
			// does and it is within the probation period delete said post, else initialize said entry.
			userposts := b.lastposts.get(m.Author.ID)
			if userposts != nil {
				t := userposts.get(m.ChannelID)
				// Checks if the user is posting within the correct time, dumb timecode below, fix later
				now := time.Now()
				then := now.Add(time.Duration(-c.TokenInterval) * time.Second)
				if t.After(then) {
					err := s.ChannelMessageDelete(m.ChannelID, m.ID)
					if err != nil {
						fmt.Println(err)
					} else {
						fmt.Println("Removed user message")
					}
					return
				}
				userposts.set(m.ChannelID, time.Now())
			} else {
				post := newChannelTimestamps()
				post.set(m.ChannelID, time.Now())
			}
		}
	}
}
