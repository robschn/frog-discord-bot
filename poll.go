package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v9"
)

type storeID struct {
	ChannelID string
	MessageID string
}

type emojiType map[string]string

func listenForPoll(s *discordgo.Session, e *discordgo.MessageCreate) {

	// moviemonday channel
	if e.ChannelID == "833899631330852934" {

		// check for poll role
		botRollCheck := false
		for _, role := range e.Member.Roles {

			if role == "767986320927883266" {
				botRollCheck = true
			}
		}

		if botRollCheck {

			// check for poll command
			if strings.Contains(e.Content, "!poll movie") {
				pollMovie(s, e)
			}

			if strings.Contains(e.Content, "!poll count") {
				countMovie(s)
			}

			if strings.Contains(e.Content, "!poll upload") {
				uploadMovie(s, e)
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

func pollMovie(s *discordgo.Session, e *discordgo.MessageCreate) {

	// connect to redis
	ctx, client := redisClient()

	votingChannel := "866455048108113941"

	// send where to vote
	s.ChannelMessageSend(e.ChannelID, fmt.Sprintf("MovieMonday™️ voting started in <#%v>", votingChannel))

	// grab 3 unwatched movies
	pickedMovies := client.SRandMemberN(ctx, "unwatched", 3).Val()

	emojiMessage := `
MovieMonday™️ voting is starting!
🧡 - **%s**
💛 - **%s**
💚 - **%s**

Please click on the emoji below to vote!
`

	votingMessage := fmt.Sprintf(emojiMessage, pickedMovies[0], pickedMovies[1], pickedMovies[2])

	// send message to channel
	messageInfo, _ := s.ChannelMessageSend(votingChannel, votingMessage)

	storeInfo := storeID{
		ChannelID: messageInfo.ChannelID,
		MessageID: messageInfo.ID,
	}

	// add emojis to message
	emojiHash := emojiType{
		"🧡": pickedMovies[0],
		"💛": pickedMovies[1],
		"💚": pickedMovies[2],
	}

	for i := range emojiHash {
		s.MessageReactionAdd(storeInfo.ChannelID, storeInfo.MessageID, i)
	}

	// store movies to vote for
	client.HSet(ctx, "voting", emojiHash)

	// store voting message info
	client.HSet(ctx, "message", storeInfo)

	// close redis connection
	client.Close()
}

func countMovie(s *discordgo.Session) {

	// connect to redis
	ctx, client := redisClient()

	// grab picked movies info
	emojiHash := client.HGetAll(ctx, "voting").Val()

	messageInfo := storeID{
		ChannelID: client.HGet(ctx, "message", "ChannelID").Val(),
		MessageID: client.HGet(ctx, "message", "MessageID").Val(),
	}

	emojiMessage, _ := s.ChannelMessage(messageInfo.ChannelID, messageInfo.MessageID)

	// compare counts to return the highest in 0 index
	for i := 1; i < len(emojiMessage.Reactions); i++ {
		if emojiMessage.Reactions[0].Count < emojiMessage.Reactions[i].Count {
			emojiMessage.Reactions[0] = emojiMessage.Reactions[i]
		}
	}

	winnerMovie := emojiHash[emojiMessage.Reactions[0].Emoji.Name]
	winnerMessage := fmt.Sprintf("The MovieMonday winner is **%s** !", winnerMovie)
	s.ChannelMessageSend(messageInfo.ChannelID, winnerMessage)

	// check for Demo
	if *DemoMode {
		s.ChannelMessageSend(messageInfo.ChannelID, "*Demo mode enabled, database will not be affected.*")
	} else {
		// move winnerMovie to watched
		client.SMove(ctx, "unwatched", "watched", winnerMovie)
	}

	// close redis connection
	client.Close()
}

func uploadMovie(s *discordgo.Session, e *discordgo.MessageCreate) {
	// connect to redis
	ctx, client := redisClient()

	// trim the command from message
	movieName := strings.TrimPrefix(e.Content, "!poll upload ")
	addMovieMessage := fmt.Sprintf("<@%v> Adding **%v** to MovieMonday Db..", e.Author.ID, movieName)

	// send message if there is no bot command
	if !(strings.Contains(addMovieMessage, "!poll")) {

		s.ChannelMessageSend(e.ChannelID, addMovieMessage)

		// check for Demo
		if *DemoMode {
			s.ChannelMessageSend(e.ChannelID, "*Demo mode enabled, database will not be affected.*")
		} else {
			client.SAdd(ctx, "unwatched", movieName)
		}

		s.ChannelMessageSend(e.ChannelID, "Done!")

		// close redis connection
		client.Close()
	}
}
