# recv

A simple message command library for [DiscordGo](https://github.com/bwmarrin/discordgo), inspired by [discord.py](https://github.com/Rapptz/discord.py).

[![GoDoc](https://godoc.org/github.com/quakecodes/recv?status.svg)](https://pkg.go.dev/github.com/quakecodes/recv)

## Features

- Command aliases
- Argument converters: convert a string input into output of any type
- Command checks (for permissions etc)
- Command usage generation
- Easy to implement in an existing DiscordGo project

## Installation

It is assumed that you have Go 1.17 installed.

```
go get github.com/quakecodes/recv
```

## Usage

A full scale project can be found in the `_example` folder, including examples for checks and converters

---

### Create a new command router

```go
router := recv.NewCommandRouter()
```

### Add a command

```go
router.AddCommand(recv.Command{
  Name:        "Ping",
  Description: "Sends \"pong\"",
  Callback: func(c *recv.CommandCtx) {
    c.Session.ChannelMessageSend(c.Message.ChannelID, "pong")
  },
})
```

### Process commands in a message

```go
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
  res, err := router.ProcessCommands(".", s, m)
  if err != nil {
    fmt.Println(err)
    return
  } else if res.Command == nil {
    s.ChannelMessageSend(m.ChannelID, "command not found")
  }
  go res.Callback()
}
```
