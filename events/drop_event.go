package events

import (
	"SmpDiscordBot/internal/imagemapper"
	"bytes"
	"fmt"
	"log"
	"sync"

	"github.com/bwmarrin/discordgo"
)

var (
	lock        = sync.Mutex{}
	id   uint16 = 0
)

type DropEvent struct {
	ChannelID string
	Message   string
}

func (t *DropEvent) CreateDrop(s *discordgo.Session, font string) {
	messageToFind := t.Message
	buffer, err := imagemapper.RenderMessageIntoImage(messageToFind, font)
	if err != nil {
		log.Printf("Failed to render image: %v\n", err)
		return
	}

	_, err = s.ChannelMessageSendComplex(t.ChannelID, &discordgo.MessageSend{
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

	t.waitForMessage(s)
}

func (t *DropEvent) waitForMessage(s *discordgo.Session) {
	resultFound := make(chan bool, 1)
	triesBeforeHint := 3
	attempt := 0

	handler := s.AddHandler(func(s *discordgo.Session, e *discordgo.MessageCreate) {
		lock.Lock()
		defer lock.Unlock()

		if e.Author.Bot || e.ChannelID != t.ChannelID {
			return
		}

		if e.Content == t.Message {
			resultFound <- true
			err := s.MessageReactionAdd(e.ChannelID, e.ID, "✅")
			if err != nil {
				log.Printf(err.Error())
				return
			}

			_, err = s.ChannelMessageSendEmbed(e.ChannelID, &discordgo.MessageEmbed{
				Title: "Drop Claimed",
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:  "By User",
						Value: fmt.Sprintf("%s", e.Author),
					},
					{
						Name:  "With Key",
						Value: fmt.Sprintf("%s", e.Content),
					},
					{
						Name:  "Reward",
						Value: "10 gems",
					},
				},

				Description: "~~---------------------------------------~~",
			})
			if err != nil {
				log.Printf(err.Error())
				return
			}

			return
		}

		attempt++
		if attempt%triesBeforeHint == 0 {
			_, err := s.ChannelMessageSendEmbed(e.ChannelID, &discordgo.MessageEmbed{
				Title: "Drop Hint",
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:  "Message Contents",
						Value: fmt.Sprintf("%s", e.Content),
					},
				},
				Description: fmt.Sprintf("%v was `%d` characters away from the answer!", e.Author.Mention(), levenshteinDP(e.Content, t.Message)),
				Footer: &discordgo.MessageEmbedFooter{
					Text: "Message Contents is the contents that the user sent.",
				},
			})
			if err != nil {
				log.Printf(err.Error())
				return
			}
		}

		err := s.MessageReactionAdd(e.ChannelID, e.ID, "❌")
		if err != nil {
			log.Printf(err.Error())
			return
		}

		fmt.Printf("%s failed to claim image with the input: %s\n", e.Author, e.Content)
	})

	<-resultFound
	fmt.Println("Message Claimed!")
	handler()
}

func levenshteinDP(s, e string) int {
	m, n := len(s), len(e)

	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	for i := 0; i <= m; i++ {
		dp[i][0] = i
	}
	for j := 0; j <= n; j++ {
		dp[0][j] = j
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if s[i-1] == e[j-1] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				dp[i][j] = 1 + minThreeArg(
					dp[i][j-1],
					dp[i-1][j],
					dp[i-1][j-1],
				)
			}
		}
	}

	return dp[m][n]
}

func minThreeArg(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}
