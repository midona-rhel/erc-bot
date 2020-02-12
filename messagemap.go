package main

import (
	"sync"

	"github.com/bwmarrin/discordgo"
)

type messageMap struct {
	sync.RWMutex
	messages map[string]discordgo.MessageCreate
}

func (m *messageMap) getMessage(key string) (discordgo.MessageCreate, bool) {
	m.RLock()
	c, ok := m.messages[key]
	m.RUnlock()
	return c, ok
}

func (m *messageMap) setMessage(key string, c *discordgo.MessageCreate) {
	m.Lock()
	m.messages[key] = *c
	m.Unlock()
}
