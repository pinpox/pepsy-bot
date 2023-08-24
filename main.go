package main

// https://discord.com/api/oauth2/authorize?client_id=1144257514813276280&permissions=40666901510209&scope=bot

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var token string
var buffer = make([][]byte, 0)

func main() {

	token = os.Getenv("DISCORD_TOKEN")

	if token == "" {
		fmt.Println("No token provided. Set the DISCORD_TOKEN environment variable")
		return
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	// Register ready as a callback for the ready events.
	dg.AddHandler(ready)

	// Register messageCreate as a callback for the messageCreate events.
	dg.AddHandler(messageCreate)

	// Register guildCreate as a callback for the guildCreate events.
	dg.AddHandler(guildCreate)

	// We need information about guilds (which includes their channels),
	// messages and voice states.
	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentMessageContent

	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Pepsybot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) when the bot receives
// the "ready" event from Discord.
func ready(s *discordgo.Session, event *discordgo.Ready) {

	// Set the playing status.
	s.UpdateGameStatus(0, "With your mom")
}

// // This function will be called (due to AddHandler above) every time a new
// // guild is joined.
func guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {

	if event.Guild.Unavailable {
		return
	}

	for _, channel := range event.Guild.Channels {
		if channel.ID == event.Guild.ID {
			_, _ = s.ChannelMessageSend(channel.ID, "Bot ready.")
			return
		}
	}
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Content == "!ranks" {
		s.ChannelMessageSend(m.ChannelID, MakeSeasonTable())
	}

	if m.Content == "!trashtalk" {
		s.ChannelMessageSend(m.ChannelID, phrases[rand.Intn(len(phrases))])
	}
}

var phrases []string = []string{

	"Nice try, but you'll need more than that to catch me!",
	"Did you get your driver's license from a cereal box?",
	"You must be a parking expert because you're always behind me!",
	"I didn't know we were racing in reverse today.",
	"Are you lost, or are you just that slow?",
	"You're like my rearview mirror – always behind.",
	"Is this your first time behind the wheel?",
	"I've seen faster speeds in a school zone.",
	"You're so slow, even snails are passing you.",
	"I'm in the winner's circle; you're in the kiddie pool.",
	"Did you take a pit stop for a coffee break?",
	"I'd say 'nice try,' but it wasn't even close.",
	"Are you racing, or just taking a scenic drive?",
	"You drive like my grandma on a Sunday outing.",
	"I'd give you a head start, but it wouldn't help.",
	"Is this a race or a Sunday drive?",
	"I'm in the fast lane; you're in the 'Oh no, I'm losing' lane.",
	"Do you need a roadmap to find the finish line?",
	"You drive like you have a blindfold on.",
	"I'm so far ahead; I can't even see you in my rearview mirror.",
	"Did you mistake the accelerator for the brake pedal again?",
	"Your car is a museum piece; it belongs in the slow lane.",
	"I'm lapping you, and I don't mean with my tongue.",
	"Even turtles are giving you side-eye right now.",
	"You're about as fast as molasses in January.",
	"My grandma drives faster than you, and she's 90!",
	"Is your car powered by hope and dreams?",
	"You should enter a 'going slow' competition instead.",
	"I didn't realize they let amateurs on the track.",
	"You must be allergic to speed.",
	"I'll send you a postcard from the winner's circle.",
	"You're so slow; you make rush hour traffic look fast.",
	"I've seen faster speed bumps.",
	"I'm setting records; you're setting cruise control.",
	"At this rate, you'll finish next week.",
	"Your car has a 'turbo' button, right? Or is it 'snail' mode?",
	"Do they give out medals for last place?",
	"You're so slow; even the tortoise is laughing.",
	"You're driving like you're paid by the hour.",
	"I hope you brought a good book; you'll be back there a while.",
	"You're not racing; you're just taking a leisurely drive.",
	"I've got a one-word review of your driving: Ouch!",
	"Your driving is so bad; even GPS can't help you.",
	"I'm not driving; I'm giving a masterclass in speed.",
	"You're so slow; you're in a different time zone.",
	"They say practice makes perfect, but not in your case.",
	"I'd ask if you're trying, but it's hard to tell.",
	"You're so slow; I'm getting déjà vu.",
	"You're like a speed bump with headlights.",
	"I'm breaking speed limits; you're breaking records... for slowness.",
	"Do you need a pit stop for a pep talk?",
	"Your car must run on 'slow motion' fuel.",
	"I'm so fast; I'm already in next week's race.",
	"Is there a 'reverse' race happening I don't know about?",
	"You're like a snail in a drag race.",
	"I've seen quicker reflexes in a sloth.",
	"You're slower than rush hour traffic.",
	"I'd slow down to let you catch up, but that's impossible.",
	"I'm leaving you in the dust, literally.",
	"I'm setting records; you're setting personal bests... for slowness.",
	"Did you get your driver's license in a cereal box?",
	"You're so slow; it's like watching paint dry.",
	"I'm in the zone; you're in the 'no chance' zone.",
	"You're so slow; you make traffic cones look fast.",
	"I'm driving circles around you, quite literally.",
	"I'm so fast; I'm aging slower than you.",
	"You're driving like you're lost in the desert.",
	"You're slower than a computer on dial-up.",
	"I'm in the lead; you're in the 'just woke up' phase.",
	"I'd say you need a pit stop, but I'm too busy winning.",
	"You're driving so slowly; even your car is yawning.",
	"I'm in first place; you're in the 'is the race over yet?' place.",
	"Are you racing or participating in a parade?",
	"I've got a trophy waiting for me; what about you?",
	"I didn't know they had a 'going nowhere' category in this race.",
	"You're so slow; even sloths are saying, 'speed up.'",
	"I'm breaking records; you're breaking the sound barrier...in reverse.",
	"You drive like you're trying to save fuel in a video game.",
	"I'm winning so hard; it's almost unfair.",
	"You're like a speed bump with a driver's license.",
	"Is your car powered by determination or just slowness?",
	"I'm at the finish line; you're still at the starting line.",
	"You must have an anchor attached to your car.",
	"I'm setting the pace; you're setting a record... for last.",
	"You're driving so slow; even the GPS is confused.",
}
