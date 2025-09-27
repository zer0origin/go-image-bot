package main

import (
	"SmpDiscordBot/events"
	"log"
	"os"
	"os/signal"

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
			//{
			//	Type:        discordgo.ApplicationCommandOptionInteger,
			//	Description: "Pick a font for the application to use",
			//	Choices: func() []*discordgo.ApplicationCommandOptionChoice {
			//		fmt.Println("DOING SOMETHING")
			//		fonts := imagemapper.FontList()
			//		res := make([]*discordgo.ApplicationCommandOptionChoice, len(fonts))
			//		var invalidChars = regexp.MustCompile(`[^a-zA-Z0-9]`)
			//		for i, font := range fonts {
			//			name := strings.ReplaceAll(invalidChars.ReplaceAllString(font, ""), "ttf", "")
			//			if len(name) > 32 {
			//				name = name[:32]
			//				fmt.Println(len(name))
			//			}
			//
			//			element := &discordgo.ApplicationCommandOptionChoice{
			//				Name:  strings.ToLower(name),
			//				Value: i,
			//			}
			//
			//			res[i] = element
			//		}
			//
			//		return res
			//	}(),
			//	Required: false,
			//},
		},
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

		event := events.DropEvent{
			Message:   guess,
			ChannelID: channel,
		}

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				// Note: this isn't documented, but you can use that if you want to.
				// This flag just allows you to create messages visible only for the caller of the command
				// (user who triggered the command)
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "Drop command sent!",
			},
		})
		if err != nil {
			log.Println(err.Error())
			return
		}

		event.CreateDrop(s)
	}
}
