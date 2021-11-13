package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/quakecodes/recv"
)

var (
	token  string
	prefix string

	router = recv.NewCommandRouter()
)

// parses command line flags - this step is not required
func init() {
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.StringVar(&prefix, "p", ".", "Bot Prefix")
	flag.Parse()
}

// basic bot setup, listens for guild messages
func main() {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.IntentsAll

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("bot is running. press ctrl-c to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// checks if the message was sent by the bot or if the message has the correct prefix
	if m.Author.ID == s.State.User.ID || !strings.HasPrefix(m.Message.Content, prefix) {
		return
	}
	res, err := router.ProcessCommands(prefix, s, m)

	// handles errors from parsing - not all errors are included here
	if err != nil {
		switch i := err.(type) {
		case recv.MissingRequiredArgumentError:
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("missing required argument \"%s\" at position %v", i.Arg.Name, i.ArgPosition))
		}
		return
	} else if res.Command == nil {
		s.ChannelMessageSend(m.ChannelID, "couldn't find command")
		return
	}

	// runs the command handler in a seperate goroutine
	go res.Callback()
}
