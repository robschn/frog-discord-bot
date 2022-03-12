package main

func main() {
	// authenticate and register the handler functions
	auth()

	go runSatanOff()
}
