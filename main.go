package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Error creating Discord session,", err)
		return
	}

	// Register the functions with for their callback events.
	dg.AddHandler(reactionAdd)
	dg.AddHandler(reactionRemove)

	// Set Identity Intent for MessageReactions
	dg.Identify.Intents = discordgo.IntentsGuildMessageReactions

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

var pronounMap = map[string]string{
	"🦋": "761015648494288926", // they/them
	"🐝": "761016805229461505", // she/her
	"🪲": "761016870286655509", // he/him
}

var cosmoMap = map[string]string{
	"♈": "874347037167058964", // Aries
	"♉": "874347209401974824", // Taurus
	"♊": "874347353979645993", // Gemini
	"♋": "874347521978286141", // Cancer
	"♌": "874347602487955476", // Leo
	"♍": "874348150893183066", // Virgo
	"♎": "874348259563413566", // Libra
	"♏": "874348351221551164", // Scorpio
	"♐": "874348449674461185", // Sagittarius
	"♑": "874348499142078495", // Capricorn
	"♒": "874348617777967106", // Aquarius
	"♓": "874348705539563531", // Pisces
}

var pronounsMessage = "866078211628073011"
var signsMessage = "874352686408007690"

func emojiCheck(emoji string) (string, bool) {
	value1, found1 := pronounMap[emoji]
	value2, found2 := cosmoMap[emoji]

	if found1 {
		return value1, true
	} else if found2 {
		return value2, true
	} else {
		return "None", false
	}
}

// Capture reactions added
func reactionAdd(s *discordgo.Session, e *discordgo.MessageReactionAdd) {

	// Check to see if emoji is our list
	value, ok := emojiCheck(e.Emoji.Name)
	if ok {
		// 				  Message ID for pronouns				 Message ID for pronouns
		if e.MessageID == pronounsMessage || e.MessageID == signsMessage {
			// if added change Role
			s.GuildMemberRoleAdd(e.GuildID, e.UserID, value)
		}
	}
}

// Capture reactions removed
func reactionRemove(s *discordgo.Session, e *discordgo.MessageReactionRemove) {

	// Check to see if emoji is our list
	value, ok := emojiCheck(e.Emoji.Name)
	if ok {
		if e.MessageID == pronounsMessage || e.MessageID == signsMessage {
			// if removed change Role
			s.GuildMemberRoleRemove(e.GuildID, e.UserID, value)
		}
	}
}
