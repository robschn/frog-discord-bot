package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

func listenForPoll(s *discordgo.Session, e *discordgo.MessageCreate) {
	if e.ChannelID == "864719238279462912" {
		// check for poll command
		if e.Content == "!poll movie" {
			// send message to user asking for the names of the 3 movies
			pickedMovies := jsonToPoll(3)

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
			emojiList := [3]string{"🧡", "💛", "💚"}
			for _, i := range emojiList {
				s.MessageReactionAdd(messageInfo.ChannelID, messageInfo.ID, i)
			}

			// sleep for a set time
			time.Sleep(10 * time.Second)

			// count emojis and pick winner
			emojiCheck, _ := s.ChannelMessage(messageInfo.ChannelID, messageInfo.ID)

			for _, e := range emojiCheck.Reactions {
				fmt.Println(e.Count, e.Emoji.Name)
				// if i > i+1 {

				// }
			}

		}
	}
}

func jsonToPoll(limit int) []string {
	// initalize struct for JSON
	type Movie struct {
		Name    string `json:"name"`
		Year    int    `json:"year"`
		Watched bool   `json:"watched"`
	}

	type Movies struct {
		Movies []Movie `json:"movies"`
	}

	// auth to GitHub Gist API

	// grab current lists
	moviesJson, err := os.Open("movies.json")
	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(moviesJson)

	var movies Movies

	// add JSON to movies list
	json.Unmarshal(byteValue, &movies)

	var unwatchedMovies []Movie

	// exclude watched moviess
	for _, mov := range movies.Movies {
		if !(mov.Watched) {
			unwatchedMovies = append(unwatchedMovies, mov)
		}
	}

	var randIndex []int

	// get random indexes
	rand.Seed(time.Now().UnixNano())
	p := rand.Perm(len(unwatchedMovies))
	randIndex = append(randIndex, p...)

	var pickedMovies []string

	// loop over the random indexes from 0 to limit
	for _, i := range randIndex[0:limit] {
		movieInfo := fmt.Sprintf("%s (%v)", unwatchedMovies[i].Name, unwatchedMovies[i].Year)
		pickedMovies = append(pickedMovies, movieInfo)
	}

	return pickedMovies
}
