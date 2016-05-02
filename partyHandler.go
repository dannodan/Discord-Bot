package main

import (
  "fmt"
  "strings"
  "github.com/bwmarrin/discordgo"
)

type Party struct {
  ID        string    `json: "id"`
  Leader    string    `json: "leader"`
  Members   []string  `json: "members"`
}

func createParty(command string, s *discordgo.Session, channelID string, user *discordgo.User) {
  _, err := readFromDatabase("players", user.ID)
  if err != nil {
    s.ChannelMessageSend(channelID, "`You don't have a character, please use $generate first`")
    return
  }
  partyName := strings.Split(strings.TrimSpace(command), " ")
  if len(partyName) > 1 {
    s.ChannelMessageSend(channelID, "`Party Names can't have spaces`")
    return
  }
  _, err = readFromDatabase("parties", partyName[0])
  if err == nil {
    s.ChannelMessageSend(channelID, "`There is already a party with that name`")
    return
  }
  party := Party{ partyName[0], user.ID, make([]string, 0) }
  party.Members = append(party.Members, user.ID)
  writeToDatabase("parties", party.ID, party)
  s.ChannelMessageSend(channelID, "`"+user.Username+" created the '"+party.ID+"' Party`") 
}

func inviteToParty(s *discordgo.Session, channelID string, user *discordgo.User, mentions []*discordgo.User) {
  fmt.Println(mentions[0].ID)
//   _, err := readFromDatabase("players", user.ID)
//   if err != nil {
//     s.ChannelMessageSend(channelID, "`You don't have a character, please use $generate first`")
//     return
//   }
//   partyName := strings.Split(strings.TrimSpace(command), " ")
//   if len(partyName) > 1 {
//     s.ChannelMessageSend(channelID, "`Party Names can't have spaces`")
//     return
//   }
//   _, err = readFromDatabase("parties", partyName[0])
//   if err == nil {
//     s.ChannelMessageSend(channelID, "`There is already a party with that name`")
//     return
//   }
//   party := Party{ partyName[0], user.ID, make([]string, 0) }
//   party.Members = append(party.Members, user.ID)
//   writeToDatabase("parties", party.ID, party)
//   s.ChannelMessageSend(channelID, "`"+user.Username+" created the '"+party.ID+"' Party`") 
}