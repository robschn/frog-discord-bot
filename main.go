package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Error creating Discord session,", err)
		return
	}

	// Register the functions with for their callback events.
	dg.AddHandler(reactionAdd)
	dg.AddHandler(reactionRemove)

	// Set Identity Intent for MessageReactions
	dg.Identify.Intents = discordgo.IntentsGuildMessageReactions

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

var emojiMap = map[string]string{
	"🦋": "761015648494288926",
	"🐝": "761016805229461505",
	"🪲": "761016870286655509",
}

// Capture reactions added
func reactionAdd(s *discordgo.Session, e *discordgo.MessageReactionAdd) {
	value, ok := emojiMap[e.Emoji.Name]
	if ok {
		if e.MessageID == "866078211628073011" {
			// if added do something in Roles
			s.GuildMemberRoleAdd(e.GuildID, e.UserID, value)
		}
	}
}

// Capture reactions removed
func reactionRemove(s *discordgo.Session, e *discordgo.MessageReactionRemove) {
	value, ok := emojiMap[e.Emoji.Name]
	if ok {
		if e.MessageID == "866078211628073011" {
			// if added do something in Roles
			s.GuildMemberRoleAdd(e.GuildID, e.UserID, value)
		}
	}
}
