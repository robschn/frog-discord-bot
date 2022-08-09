package main

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v9"
)

func listenForPoll(s *discordgo.Session, e *discordgo.MessageCreate) {

	if e.ChannelID == "864719238279462912" {
		// check for poll command
		if e.Content == "!poll movie" {

			// connect to redis database
			ctx := context.TODO()

			client := redis.NewClient(&redis.Options{
				Addr:     *RedisUrl,
				Password: *RedisPass,
				DB:       0,
			})

			// grab 3 unwatched movies
			pickedMovies := client.SRandMemberN(ctx, "unwatched", 3).Val()

			emojiMessage := `(A)here

MovieMonday™️ voting is starting!
🧡 - %s
💛 - %s
💚 - %s

Please click on the emoji below to vote!
Voting ends at midnight on Sunday.`

			votingMessage := fmt.Sprintf(emojiMessage, pickedMovies[0], pickedMovies[1], pickedMovies[2])

			// send message to channel
			messageInfo, _ := s.ChannelMessageSend(e.ChannelID, votingMessage)

			// add emojis to message
			emojiHash := map[string]interface{}{
				"🧡": pickedMovies[0],
				"💛": pickedMovies[1],
				"💚": pickedMovies[2],
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
			winnerMovie := emojiHash[emojiCheck.Reactions[0].Emoji.Name]
			winnerMessage := fmt.Sprintf("The MovieMonday winner is **%s** !", emojiHash[emojiCheck.Reactions[0].Emoji.Name])
			s.ChannelMessageSend(messageInfo.ChannelID, winnerMessage)

			// move winnerMovie to watched
			client.SMove(ctx, "unwatched", "watched", winnerMovie)

			// close redis connection
			client.Close()
		}

		if e.Content == "!poll upload" {

		}
	}
}
