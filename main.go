package main

import (
	"os"

	"github.com/bwmarrin/discordgo"
)

func main() {

	discordToken := os.Getenv("DISCORD_TOKEN")

	discordgo.New(discordToken)

}
