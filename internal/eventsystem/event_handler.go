package eventsystem

import (
	"SmpDiscordBot/events"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	bindings = make(map[string]func(s *discordgo.Session, e any))
)

func SetupHandlers(s *discordgo.Session) {
	s.AddHandler(onMessage)
}

func LoadEvent[T any](event events.Event, listener func(s *discordgo.Session, e T)) {
	if event == nil {
		log.Print("Cannot load nil event")
		return
	}

	handlers := event.Handlers()

	for _, handler := range handlers {
		tStr := reflect.TypeOf(handler).String()
		tStrParts := strings.Split(reflect.TypeOf(handler).String(), ",")
		tStr = tStrParts[len(tStrParts)-1]
		tStr = strings.ToLower(tStr)

		evtName := strings.Trim(reflect.TypeOf(event).String(), " *()")
		evtNameParts := strings.Split(evtName, ".")
		evtName = evtNameParts[len(evtNameParts)-1]
		evtName = strings.ToLower(evtName)

		lsnType := strings.Trim(tStr, " *()")
		lsnType = strings.ReplaceAll(lsnType, "discordgo.", "")
		lsnType = strings.ToLower(lsnType)

		bind := fmt.Sprintf("%s.%d.%s", evtName, event.UID(), lsnType)
		RegisterListener(bind, listener)
		log.Printf("Loaded %s with binding %s", evtName, bind)
	}
}

func RegisterListener[T any](bind string, handler func(s *discordgo.Session, e T)) {
	parts := strings.Split(bind, ".")
	if len(parts) != 3 {
		log.Printf("invalid event binding for %s\n", bind)
		return
	}

	wrapped := func(s *discordgo.Session, e any) {
		evt, ok := e.(T)
		if !ok {
			log.Printf("Event type mismatch for %s\n", bind)
			return
		}
		handler(s, evt)
	}
	log.Printf("Registered Listener for bind %s", bind)
	bindings[bind] = wrapped
}

func FireEvent(bind string, s *discordgo.Session, e any) {
	parts := strings.Split(strings.ToLower(bind), ".")
	if len(parts) != 3 {
		log.Printf("invalid event binding for %s\n", strings.ToLower(bind))
		return
	}

	if parts[0] == "*" {
		for s2, f := range bindings {
			if strings.Contains(s2, strings.ToLower(parts[2])) {
				f(s, e)
			}
		}

		return
	}

	if parts[1] == "*" {
		for s2, f := range bindings {
			if strings.Contains(s2, strings.ToLower(parts[0])) && strings.Contains(s2, strings.ToLower(parts[2])) {
				f(s, e)
			}
		}

		return
	}

	if f, ok := bindings[bind]; ok {
		f(s, e)
	}
}

func onMessage(s *discordgo.Session, e *discordgo.MessageCreate) {
	if !e.Author.Bot {
		log.Printf("EventHandler - %s:%s\n", e.Author, e.Content)
	}

	FireEvent("*.*.MessageCreate", s, e)
}
