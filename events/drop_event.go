package events

import (
	"SmpDiscordBot/internal/imagemapper"
	"bytes"
	"log"
	"sync"

	"github.com/bwmarrin/discordgo"
)

var (
	lock        = sync.Mutex{}
	id   uint16 = 0
)

type DropEvent struct {
	id uint16
}

func (t *DropEvent) UID() uint16 {
	lock.Lock()
	defer lock.Unlock()
	if t.id == 0 {
		id++
		t.id = id
	}

	return t.id
}

func (t *DropEvent) Handlers() []any {
	return []any{t.OnMessage}
}

func (t *DropEvent) CreateNewDrop() {

}

func (t *DropEvent) OnMessage(s *discordgo.Session, e *discordgo.MessageCreate) {
	if e.Author.Bot {
		return
	}

	buffer, err := imagemapper.RenderMessageIntoImage(e.Content)
	if err != nil {
		log.Printf("Failed to render image: %v\n", err)
	}

	_, err = s.ChannelMessageSendComplex(e.ChannelID, &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title:       "Drop",
			Description: "Type the text in the image to collect the drop!",
			Image: &discordgo.MessageEmbedImage{
				URL: "attachment://drop.png",
			},
			Color: 0xff00b7,
		},
		Files: []*discordgo.File{
			{
				Name:   "drop.png",
				Reader: bytes.NewReader(buffer.Bytes()),
			},
		},
	})
}
