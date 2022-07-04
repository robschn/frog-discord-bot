package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var messagePollAuthor string

func listenForPoll(s *discordgo.Session, e *discordgo.MessageCreate) {
	if e.ChannelID == "864719238279462912" {
		if e.Content == "!poll movie" {
			// send message to user asking for the names of the 3 movies
			messagePollAuthor := e.Author.ID
			messageResponse := fmt.Sprintf("<@%s> please type in 3 movie names separated by commas.", messagePollAuthor)
			s.ChannelMessageSend(e.ChannelID, messageResponse)
		}
		if e.Author.ID == messagePollAuthor {
			s.ChannelMessageSend(e.ChannelID, "Hey")
		}
	}
}
