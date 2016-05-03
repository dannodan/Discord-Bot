package main

import (
	"fmt"
  	"strings"
	"github.com/bwmarrin/discordgo"
)

func main() {

	// Create a new Discord session using the provided login information.
	// Use discordgo.New(Token) to just use a token for login.
	dg, err := discordgo.New("MTc2MTUyNTk5OTExODU4MTc4.Cgb3lw.fD_D5yhrHgFOgp_G3CE5VpZlJGs")
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	// Test
	// Register messageCreate as a callback for the messageCreate events.
	dg.AddHandler(messageCreate)

	// Open the websocket and begin listening.
	dg.Open()

	fmt.Println("GM Bot is now running.  Press CTRL-C to exit.")
	// Simple way to keep program running until CTRL-C is pressed.
	<-make(chan struct{})
	return
}

func verifyUpdateArgs(args string, s *discordgo.Session, chID string, user *discordgo.User) {
	arguments := strings.Split(strings.TrimSpace(args), " ")
	if len(arguments) < 2 || len(arguments) > 2 {
		s.ChannelMessageSend(chID, "`You need only the Stat you wish to update and the number of points to allocate`")
		return
	}
	switch {
		case strings.EqualFold(arguments[0],"STR") || strings.EqualFold(arguments[0],"Strength"):
			allocateStatPoints(s, chID, user, "Strength", strings.TrimSpace(arguments[1]))
		case strings.EqualFold(arguments[0],"INT") || strings.EqualFold(arguments[0],"Intelligence"):
			allocateStatPoints(s, chID, user, "Intelligence", strings.TrimSpace(arguments[1]))
		case strings.EqualFold(arguments[0],"DEX") || strings.EqualFold(arguments[0],"Dexterity"):
			allocateStatPoints(s, chID, user, "Dexterity", strings.TrimSpace(arguments[1]))
		case strings.EqualFold(arguments[0],"VIT") || strings.EqualFold(arguments[0],"Vitality"):
			allocateStatPoints(s, chID, user, "Vitality", strings.TrimSpace(arguments[1]))
		case strings.EqualFold(arguments[0],"SPR") || strings.EqualFold(arguments[0],"Spirit"):
			allocateStatPoints(s, chID, user, "Spirit", strings.TrimSpace(arguments[1]))
		case strings.EqualFold(arguments[0],"AGI") || strings.EqualFold(arguments[0],"Agility"):
			allocateStatPoints(s, chID, user, "Agility", strings.TrimSpace(arguments[1]))
		default:
			s.ChannelMessageSend(chID, "`Not a valid status`")
	}
}

// func verifyPartyInviteArgs(args string, s *discordgo.Session, chID string, user *discordgo.User, mentions []*discordgo.User) {

// }

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

  // Print message to stdout.
	chID := m.ChannelID
	content := m.Content
	mentions := m.Mentions
	user := m.Author

  if strings.HasPrefix(content, "$") {
    command := strings.SplitN(strings.TrimPrefix(content, "$"), " ", 2)
    switch {
			case strings.EqualFold(command[0],"stats"):
        playerStats(s, chID, user)
			case strings.EqualFold(command[0],"generate"):
        generatePlayer(s, chID, user)
			case strings.EqualFold(command[0],"allocate"):
				if len(command) < 2 {
					s.ChannelMessageSend(chID, "`You need to have the Stat you wish to update and the number of points to allocate`")
				} else {
					verifyUpdateArgs(command[1], s, chID, user)
				}
			case strings.EqualFold(command[0], "pcreate"):
			if len(command) < 2 {
				s.ChannelMessageSend(chID, "`You need to specify a Party Name`")
			} else {
				createParty(command[1], s, chID, user)
			}
			case strings.EqualFold(command[0], "pinvite"):
			if len(command) < 2 {
				s.ChannelMessageSend(chID, "`You need someone to invite to the party`")
			} else {
				inviteToParty(s, chID, user, mentions)
			}
			case strings.EqualFold(command[0], "pleave"):
				leaveParty(s, chID, user)
			case strings.EqualFold(command[0], "pdisband"):
				disbandParty(s, chID, user)
			case strings.EqualFold(command[0], "test"):
				testing(user)
			case strings.EqualFold(command[0], "party"):
				showParty(s, chID, user)
			case strings.EqualFold(command[0], "commands"):
				s.ChannelMessageSend(chID, "```$generate\t$stats\t$allocate\t$pcreate\t$party\t$pinvite\t$pleave\t$pdisband\t$commands```")
      default:
        s.ChannelMessageSend(chID, "`Not a valid command`")
    }
  }
}
