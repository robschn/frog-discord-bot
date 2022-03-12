package main

import "time"

// https://github.com/jamestjw/lyrical/blob/master/poll/poll.go

func runSatanOff() {
	// check to see if we can run poll
	if startSatanOff("Time example here") {
		// get user list
		registerContestants()

		// run voting
		winnerUser := runVoting()

		// change the role from returned voting
		changeRole(winnerUser)

		time.Sleep(144 * time.Hour)

	} else {
		time.Sleep(30 * time.Minute)
	}
}

// send message that Satan off is starting on Sunday
func startSatanOff(timeToStart string) bool {
	// start when time condition is met
	var shouldStart bool
	// return bool for whether the contest should start
	return shouldStart
}

// grab participants
func registerContestants() []string {
	// send message for users to register
	// current king
	// these will be users that show in poll
	// possibly allow to choose their own emoji
	// cannot vote

	// green checkmark to play

	// no green checkmark then tell us who reamins King

	// Leave voting up for 30 minutues

	// Warn users at 2 mins
	// display who is playing

	// Tell chat that voting is closed and display contestants
	// remove FrogBot from users

	// Wish good luck to participants and chat
	var registeredUsers []string
	// return a slice of registered user
	return registeredUsers
}

// Around 5PM open voting
func runVoting() []string {
	// display participants

	// For every person, assign them to an emoji
	// Have users assign emojis
	// !myEmote :emote:

	// after an hour, count the emojis and display winner
	// maybe pin post
	// if tie, the each tie wins

	// change role to:
	// Literally Lucifer - One winner
	// Co-Satan - Two or more winners
	// should leverage the role.go we made
	// we're returning a slice for more than one winner
	var userWinner []string
	//
	return userWinner
}
