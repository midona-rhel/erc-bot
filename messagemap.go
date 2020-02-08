package main

type messageMap struct {
	sync.RWM
	map[string]*discordgo.MessageCreate
}

