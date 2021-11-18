package main

import (
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/quakecodes/recv"
)

var (
	intConverter = &recv.Converter{
		Name: "Int",
		Callback: func(c *recv.ConverterCtx) (interface{}, bool) {
			num, err := strconv.Atoi(c.Input)
			if err != nil {
				c.Session.ChannelMessageSend(c.Message.ChannelID, fmt.Sprintf("failed to convert \"%s\" into an integer", c.Input))
				return nil, false
			}
			return num, true
		},
	}

	adminCheck = recv.CommandCheck{
		Name: "Is Admin",
		Callback: func(c *recv.Ctx) bool {
			perms, err := c.Session.UserChannelPermissions(c.Message.Author.ID, c.Message.ChannelID)
			if err != nil || perms&discordgo.PermissionAdministrator == 0 {
				c.Session.ChannelMessageSend(c.Message.ChannelID, "missing permissions: you are not an admin")
				return false
			}
			return true
		},
	}
)

func init() {
	// simple command with no arguments
	router.AddCommand(recv.Command{
		Name:        "Ping",
		Description: "Sends \"pong\"",
		Callback: func(c *recv.CommandCtx) {
			c.Session.ChannelMessageSend(c.Message.ChannelID, "pong")
		},
	})

	// simple command with a single argument with no converter
	router.AddCommand(recv.Command{
		Name:        "Repeat",
		Description: "Repeats the provided text",
		Args: []recv.CommandArg{
			{
				Name:        "text",
				Description: "The text that you want to repeat",
			},
		},
		Callback: func(c *recv.CommandCtx) {
			c.Session.ChannelMessageSend(c.Message.ChannelID, fmt.Sprint(c.Args[0].(string)))
		},
	})

	// more complex command with multiple arguments and converters
	router.AddCommand(recv.Command{
		Name:        "Add",
		Description: "Adds up to 3 numbers together",
		Args: []recv.CommandArg{
			{
				Name:        "num1",
				Description: "The first number to number",
				Converter:   intConverter,
			},
			{
				Name:        "num2",
				Description: "The second number to number",
				Converter:   intConverter,
				Optional:    true,
			},
			{
				Name:        "num3",
				Description: "The third number to number",
				Optional:    true,
				Converter:   intConverter,
			},
		},
		NoJoin: true,
		Callback: func(c *recv.CommandCtx) {
			sum := 0
			for _, num := range c.Args {
				sum += num.(int)
			}
			c.Session.ChannelMessageSend(c.Message.ChannelID, fmt.Sprint(sum))
		},
	})

	// command that fetches the usage for a specified command - demonstrates use of GetCommand()
	router.AddCommand(recv.Command{
		Name:        "Usage",
		Description: "Fetches the usage for a command",
		Args: []recv.CommandArg{
			{
				Name:        "command",
				Description: "The command to fetch usage for",
				Converter: &recv.Converter{
					Name: "Command",
					Callback: func(c *recv.ConverterCtx) (interface{}, bool) {
						comm, ok := router.GetCommand(c.Input)
						if !ok {
							c.Session.ChannelMessageSend(c.Message.ChannelID, fmt.Sprintf("command \"%s\" not found", c.Input))
							return nil, false
						}
						return comm, true
					},
				},
			},
		},
		Callback: func(c *recv.CommandCtx) {
			comm := c.Args[0].(recv.Command)
			c.Session.ChannelMessageSend(c.Message.ChannelID, fmt.Sprintf("usage for: \"%s\"\n```%s%s```", comm.Name, prefix, comm.Usage()))
		},
	})

	// command that runs only if the user is an admin - this requires guild members intent
	router.AddCommand(recv.Command{
		Name:        "adminonly",
		Description: "Can only be run by admins (hopefully)",
		Checks: []recv.CommandCheck{
			adminCheck,
		},
		Callback: func(c *recv.CommandCtx) {
			c.Session.ChannelMessageSend(c.Message.ChannelID, "hi admin")
		},
	})

	// command that only runs if the user's name is QuaKe and they have admin
	router.AddCommand(recv.Command{
		Name:        "quakeadminonly",
		Description: "quake admin only: otherwise go away",
		Checks: []recv.CommandCheck{
			adminCheck,
			{
				Name: "Is QuaKe",
				Callback: func(c *recv.Ctx) bool {
					if !(c.Message.Author.Username == "QuaKe") {
						c.Session.ChannelMessageSend(c.Message.ChannelID, "go away")
						return false
					}
					return true
				},
			},
		},
		Callback: func(c *recv.CommandCtx) {
			c.Session.ChannelMessageSend(c.Message.ChannelID, "hello quake")
		},
	})
}
