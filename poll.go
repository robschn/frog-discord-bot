package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v9"
)

func listenForPoll(s *discordgo.Session, e *discordgo.MessageCreate) {
	if e.ChannelID == "864719238279462912" {
		// check for poll command
		if e.Content == "!poll movie" {
			// send message to user asking for the names of the 3 movies
			pickedMovies := fetchMovies(3)

			// format to add 7 days to current date
			nextWeek := time.Now().AddDate(0, 0, 6)
			formatWeek := fmt.Sprintf("%s, %v\n", nextWeek.Month(), nextWeek.Day())
			emojiMessage := `(A)here

MovieMonday™️ voting for %s
🧡 - %s
💛 - %s
💚 - %s

Please click on the emoji below to vote!
Voting ends at midnight on Sunday.`

			votingMessage := fmt.Sprintf(emojiMessage, formatWeek, pickedMovies[0], pickedMovies[1], pickedMovies[2])

			// send message to channel
			messageInfo, _ := s.ChannelMessageSend(e.ChannelID, votingMessage)

			// add emojis to message
			emojiHash := map[string]string{
				"🧡": pickedMovies[0].Value,
				"💛": pickedMovies[1].Value,
				"💚": pickedMovies[2].Value,
			}

			for i := range emojiHash {
				s.MessageReactionAdd(messageInfo.ChannelID, messageInfo.ID, i)
			}

			// sleep for a set time
			time.Sleep(10 * time.Second)

			// grab message info
			emojiCheck, _ := s.ChannelMessage(messageInfo.ChannelID, messageInfo.ID)

			// compare counts to return the highest in 0 index
			for i := 1; i < len(emojiCheck.Reactions); i++ {
				if emojiCheck.Reactions[0].Count < emojiCheck.Reactions[i].Count {
					emojiCheck.Reactions[0] = emojiCheck.Reactions[i]
				}
			}
			winnerMovie := fmt.Sprintf("The MovieMonday winner is %s !", emojiHash[emojiCheck.Reactions[0].Emoji.Name])
			s.ChannelMessageSend(messageInfo.ChannelID, winnerMovie)
		}
		if e.Content == "!poll upload" {

		}
	}
}

func fetchMovies(limit int) []redis.KeyValue {

	ctx := context.TODO()

	// connect to redis database
	r := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       0,
	})

	return r.HRandFieldWithValues(ctx, "unwatched", limit).Val()
}
