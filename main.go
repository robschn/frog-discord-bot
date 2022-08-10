package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	BotToken  = flag.String("t", "", "Discord bot access token")
	RedisUrl  = flag.String("c", "", "Redis database URL to include port")
	RedisPass = flag.String("p", "", "Redis database password")
	DemoMode  = flag.Bool("demo", false, "Enable demo mode to leave database untouched")
)

func init() { flag.Parse() }

func main() {

	// auth to discord
	dg, err := discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	// Register the functions with for their callback events.
	dg.AddHandler(reactionAdd)
	dg.AddHandler(reactionRemove)
	dg.AddHandler(listenForPoll)

	// Set Identity Intent for functions
	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildMessageReactions

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
