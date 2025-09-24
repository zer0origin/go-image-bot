package main

import (
	"SmpDiscordBot/events"
	"SmpDiscordBot/internal/eventsystem"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

var (
	token = os.Getenv("BOT_TOKEN")
)

func main() {
	client, err := discordgo.New("Bot " + token)
	if err != nil {
		panic("failed: " + err.Error())
	}

	commands := []*discordgo.ApplicationCommand{{
		Name:        "write",
		Description: "Send a message as the bot user account!",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "message",
				Description: "Write a message for the application to send.",
				Required:    true,
			},
		},
	}}

	client.AddHandler(onReady)
	//client.AddHandler(onInteraction)

	eventsystem.SetupHandlers(client)

	drop := &events.DropEvent{}
	eventsystem.LoadEvent(drop, drop.OnMessage)

	client.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAllWithoutPrivileged)

	err = client.Open()
	if err != nil {
		log.Printf("Cannot open the session: %v", err)
	}
	defer client.Close()

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := client.ApplicationCommandCreate(client.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	//log.Println("Removing commands...")
	//for _, v := range registeredCommands {
	//	err := client.ApplicationCommandDelete(client.State.User.ID, "1188337965789884426", v.ID)
	//	if err != nil {
	//		log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
	//	}
	//}

	log.Println("Gracefully shutting down.")
}

func onReady(s *discordgo.Session, e *discordgo.Ready) {
	log.Println("Bot is ready")
}
