package events

import "github.com/bwmarrin/discordgo"

type Event interface {
	Handlers(s *discordgo.Session)
}
