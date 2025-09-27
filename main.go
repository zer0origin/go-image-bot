package main

import (
	"SmpDiscordBot/events"
	"SmpDiscordBot/internal/imagemapper"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	token = os.Getenv("BOT_TOKEN")
)

func main() {
	client, err := discordgo.New("Bot " + token)
	if err != nil {
		panic("failed: " + err.Error())
	}

	commands := []*discordgo.ApplicationCommand{{
		Name:        "drop",
		Description: "Sends a drop to the target channel!",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "guess",
				Description: "What should the users be trying to guess?",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionChannel,
				Name:        "channel",
				Description: "Write a message for the application to send.",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "font",
				Description: "What font should the bot use? /fonts to view available fonts",
				Required:    false,
			},
		},
	},
		{
			Name:        "fonts",
			Description: "Display fonts available to the bot",
		}}

	client.AddHandler(onReady)
	client.AddHandler(onInteraction)

	client.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAllWithoutPrivileged)

	err = client.Open()
	if err != nil {
		log.Printf("Cannot open the session: %v", err)
	}
	defer client.Close()

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := client.ApplicationCommandCreate(client.State.User.ID, "1188337965789884426", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	log.Println("Removing commands...")
	for _, v := range registeredCommands {
		err := client.ApplicationCommandDelete(client.State.User.ID, "1188337965789884426", v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}

	log.Println("Gracefully shutting down.")
}

func onReady(s *discordgo.Session, e *discordgo.Ready) {
	log.Println("Bot is ready")
}

func onInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name == "drop" {
		guess := ""
		channel := i.ChannelID
		font := ""

		options := i.ApplicationCommandData().Options
		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
		for _, opt := range options {
			optionMap[opt.Name] = opt
		}

		if opt, ok := optionMap["guess"]; ok {
			guess = opt.StringValue()
		}

		if opt, ok := optionMap["channel"]; ok {
			channel = opt.StringValue()
		}

		if opt, ok := optionMap["font"]; ok {
			stringValue := opt.StringValue()
			if !imagemapper.HasFont(stringValue) {
				log.Printf("tried to load invalid font %s\n", stringValue)
				return
			}

			font = stringValue
		}

		event := events.DropEvent{
			Message:   guess,
			ChannelID: channel,
		}

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "Drop command sent!",
			},
		})
		if err != nil {
			log.Println(err.Error())
			return
		}

		event.CreateDrop(s, font)
	}

	if i.ApplicationCommandData().Name == "fonts" {
		fonts := imagemapper.FontList()

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: strings.Join(fonts, ", "),
			},
		})
		if err != nil {
			log.Println(err.Error())
			return
		}
	}
}
