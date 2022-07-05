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

var messagePollAuthor string

func listenForPoll(s *discordgo.Session, e *discordgo.MessageCreate) {
	if e.ChannelID == "864719238279462912" {
		// check to see if messagePollAuthor is set
		if e.Author.ID == messagePollAuthor {
			jsonToPoll(3)
		}
		// check for poll command
		if e.Content == "!poll movie" {
			// send message to user asking for the names of the 3 movies
			messagePollAuthor = e.Author.ID
			messageResponse := fmt.Sprintf("<@%s> please type in 3 movie names separated by commas.", messagePollAuthor)
			s.ChannelMessageSend(e.ChannelID, messageResponse)
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
