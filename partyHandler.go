package main

import (
  // "fmt"
  "strings"
  "encoding/json"
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
  // player := dbread.(map[string]interface{})
  party := Party{ partyName[0], user.ID, make([]string, 0) }
  party.Members = append(party.Members, user.ID)
  writeToDatabase("parties", party.ID, party)
  updated := map[string]string{"Party":party.ID}
  updateToDatabase("players", user.ID, updated)
  s.ChannelMessageSend(channelID, "`"+user.Username+" created the '"+party.ID+"' Party`")
}

func inviteToParty(s *discordgo.Session, channelID string, user *discordgo.User, mentions []*discordgo.User) {
  party := Party{}
  player, err := readFromDatabase("players", user.ID)
  if err != nil {
    s.ChannelMessageSend(channelID, "`You don't have a character, please use $generate first`")
    return
  }
  _, err = readFromDatabase("players", mentions[0].ID)
  if err != nil {
    s.ChannelMessageSend(channelID, "`The user you are trying to invite doesn't have a character, please tell them to use $generate first`")
    return
  }
  if player["Party"].(string) == "" {
    s.ChannelMessageSend(channelID, "`You don't have a party, please create one with $pcreate 'PartyName'`")
    return
  }
  dbparty, err := readFromDatabase("parties", player["Party"].(string))
  if err != nil {
    s.ChannelMessageSend(channelID, "`There was an error inviting to the Party`")
    return
  }
  aux, err := json.Marshal(dbparty)
  if err != nil {
    panic(err)
  }
  if err := json.Unmarshal(aux, &party); err != nil {
    panic(err)
  }
  if player["ID"].(string) != party.Leader {
    s.ChannelMessageSend(channelID, "`You are not the party leader`")
    return
  }
  party.Members = append(party.Members, mentions[0].ID)
  writeToDatabase("parties", party.ID, party)
  updated := map[string]string{"Party":party.ID}
  updateToDatabase("players", mentions[0].ID, updated)
  s.ChannelMessageSend(channelID, "`Added "+mentions[0].Username+" to the '"+party.ID+"' Party`")
}
