package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v9"
)

func listenForPoll(s *discordgo.Session, e *discordgo.MessageCreate) {

	if e.ChannelID == "833899631330852934" {
		// check for poll command
		if strings.Contains(e.Content, "!poll movie") {

			// initalize hours sleep map
			var hoursSleep int

			// grab time amount
			if strings.ContainsAny(e.Content, "1234567890.") {
				stringHours := strings.TrimPrefix(e.Content, "!poll movie ")
				hoursSleep, _ = strconv.Atoi(stringHours)
				if hoursSleep > 24 {
					defaultMessage := fmt.Sprintf("%v hours is too long! Defaulting to 24 hours", hoursSleep)
					hoursSleep = 24
					s.ChannelMessageSend(e.ChannelID, defaultMessage)
				}

			} else {
				hoursSleep = 4
			}

			// connect to redis
			ctx, client := redisClient()

			// grab 3 unwatched movies
			pickedMovies := client.SRandMemberN(ctx, "unwatched", 3).Val()

			emojiMessage := `(A)here

MovieMonday™️ voting is starting!
🧡 - %s
💛 - %s
💚 - %s

Please click on the emoji below to vote!
Voting ends in %v hours.`

			votingMessage := fmt.Sprintf(emojiMessage, pickedMovies[0], pickedMovies[1], pickedMovies[2], hoursSleep)

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

			// sleep for time
			time.Sleep(time.Duration(hoursSleep) * time.Second)

			// grab message info
			emojiCheck, _ := s.ChannelMessage(messageInfo.ChannelID, messageInfo.ID)

			// compare counts to return the highest in 0 index
			for i := 1; i < len(emojiCheck.Reactions); i++ {
				if emojiCheck.Reactions[0].Count < emojiCheck.Reactions[i].Count {
					emojiCheck.Reactions[0] = emojiCheck.Reactions[i]
				}
			}
			winnerMovie := emojiHash[emojiCheck.Reactions[0].Emoji.Name]
			winnerMessage := fmt.Sprintf("The MovieMonday winner is **%s** !", winnerMovie)
			s.ChannelMessageSend(messageInfo.ChannelID, winnerMessage)

			// move winnerMovie to watched
			client.SMove(ctx, "unwatched", "watched", winnerMovie)

			// close redis connection
			client.Close()
		}

		if strings.Contains(e.Content, "!poll upload") {
			// trim the command from message
			movieName := strings.TrimPrefix(e.Content, "!poll upload ")
			addMovieMessage := fmt.Sprintf("<@%v> Adding **%v** to MovieMonday Db..", e.Author.ID, movieName)

			// send message if there is no bot command
			if !(strings.Contains(addMovieMessage, "!poll")) {

				s.ChannelMessageSend(e.ChannelID, addMovieMessage)

				// upload to redis
				ctx, client := redisClient()
				client.SAdd(ctx, "unwatched", movieName)

				s.ChannelMessageSend(e.ChannelID, "Done!")

				// close redis connection
				client.Close()
			}
		}
	}
}

func redisClient() (context.Context, *redis.Client) {
	// connect to redis database
	ctx := context.TODO()

	client := redis.NewClient(&redis.Options{
		Addr:     *RedisUrl,
		Password: *RedisPass,
		DB:       0,
	})

	return ctx, client
}
