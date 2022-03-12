package main

import "time"

func main() {
	// authenticate and register the handler functions
	auth()

	// check to see if we can run poll
	if startSatanOff("Time example here") {
		// get user list
		registerContestants()

		// run voting
		winnerUser := runVoting()

		// change the role from returned voting
		changeRole(winnerUser)

	} else {
		time.Sleep(30 * time.Minute)
	}
}
