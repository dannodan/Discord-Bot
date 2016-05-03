package main

import (
//   "fmt"
  "strings"
  "math"
  // "encoding/json"
  "github.com/bwmarrin/discordgo"
)

type Party struct {
  ID        string            `json: "id"`
  Members   map[string]bool   `json: "members"`
}

type MemberList struct {
  ID        string                  `json: "id"`
  Leader    bool                    `json: "leader"`
  Data      map[string]interface{}  `json: "data"`
}

func createParty(command string, s *discordgo.Session, channelID string, user *discordgo.User) {
  player, err := readFromDatabase("players", user.ID)
  if err != nil {
    s.ChannelMessageSend(channelID, "`You don't have a character, please use $generate first`")
    return
  }
  if player["Party"] != "" {
    s.ChannelMessageSend(channelID, "`You are already in a party, leave it first`")
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
  party := Party{ partyName[0], make(map[string]bool) }
  party.Members[player["ID"].(string)] = true
  writeToDatabase("parties", party.ID, party)
  updated := map[string]string{"Party":party.ID}
  updateToDatabase("players", user.ID, updated)
  s.ChannelMessageSend(channelID, "`"+user.Username+" created the '"+party.ID+"' Party`")
}

func showParty(s *discordgo.Session, channelID string, user *discordgo.User) {
  player, err := readFromDatabase("players", user.ID)
  if err != nil {
    s.ChannelMessageSend(channelID, "`You haven't created a character, please use $generate first`")
    return
  }
  if player["Party"] == "" {
    s.ChannelMessageSend(channelID, "`You don't belong to any party`")
    return
  }
  party, err := readFromDatabase("parties", player["Party"].(string))
  if err != nil {
    s.ChannelMessageSend(channelID, "`There was an error showing the party`")
    return
  }
  index := 1.0
  partyMembers := make([]string,0)
  for key, value := range party["Members"].(map[string]interface{}) {
    aux, err := readFromDatabase("players", key)
    if err != nil {
      panic(err)
    }
    if value.(bool) {
      partyMembers = append([]string{"The Leader is "+aux["Name"].(string)+"\n\nThe members are:\n"}, strings.Join(partyMembers, ""))
    } else if math.Mod(index, 3.0) == 0 {
      partyMembers = append(partyMembers, ""+aux["Name"].(string)+"\n")
      index = index + 1.0
    } else {
      partyMembers = append(partyMembers, ""+aux["Name"].(string)+"\t")
      index = index + 1.0
    }
  }
  partyMembers = append([]string{ "```"+party["ID"].(string)+" Party\n\n" }, strings.Join(partyMembers, ""))
  partyMembers = append(partyMembers, "```")
  message := strings.Join(partyMembers, "")
  s.ChannelMessageSend(channelID, message)
}

func leaveParty(s *discordgo.Session, channelID string, user *discordgo.User) {
  player, err := readFromDatabase("players", user.ID)
  if err != nil {
    s.ChannelMessageSend(channelID, "`You don't have a character, please use $generate first`")
    return
  }
  if player["Party"] == "" {
    s.ChannelMessageSend(channelID, "`You don't have a party to leave`")
    return
  }
  party, err := readFromDatabase("parties", player["Party"].(string))
  if err != nil {
    s.ChannelMessageSend(channelID, "`There was an error leaving from the party`")
    return
  }
  message := "`"+user.Username+" left the '"+party["ID"].(string)+"' Party`"
  partyMembers := party["Members"].(map[string]interface{})
  if partyMembers[player["ID"].(string)] == true {
    delete(partyMembers, player["ID"].(string))
    if len(partyMembers) <= 0 {
      member := map[string]string{"Party":""}
      updateToDatabase("players", user.ID, member)
      deleteFromDatabase("parties", party["ID"].(string))
      s.ChannelMessageSend(channelID, "`The "+party["ID"].(string)+" Party has disbanded`")
      return
    }
    for key, _ := range partyMembers {
      newLeader, err := readFromDatabase("players", key)
      if err != nil {
        panic(err)
      }
      message = message+"\n`"+newLeader["Name"].(string)+" is the new leader of the Party`"
      partyMembers[key] = true
      break
    }
  } else {
    delete(partyMembers, player["ID"].(string))
  }
  party["Members"] = partyMembers
  writeToDatabase("parties", party["ID"].(string), party)
  updated := map[string]string{"Party":""}
  updateToDatabase("players", user.ID, updated)
  s.ChannelMessageSend(channelID, message)
}

func disbandParty(s *discordgo.Session, channelID string, user *discordgo.User) {
  player, err := readFromDatabase("players", user.ID)
  if err != nil {
    s.ChannelMessageSend(channelID, "`You don't have a character, please use $generate first`")
    return
  }
  if player["Party"] == "" {
    s.ChannelMessageSend(channelID, "`You don't have a party to disband`")
    return
  }
  party, err := readFromDatabase("parties", player["Party"].(string))
  if err != nil {
    s.ChannelMessageSend(channelID, "`There was an error disbanding the party`")
    return
  }
  cleanParty(party["ID"].(string))
  deleteFromDatabase("parties", party["ID"].(string))
  s.ChannelMessageSend(channelID, "`The "+party["ID"].(string)+" Party has disbanded`")
}

func cleanParty(partyID string) {
  party, err := readFromDatabase("parties", partyID)
  if err != nil {
    panic(err)
  }
  var updated map[string]string
  for key, _ := range party["Members"].(map[string]interface{}) {
    updated = map[string]string{"Party":""}
    updateToDatabase("players", key, updated)
  }
}

func inviteToParty(s *discordgo.Session, channelID string, user *discordgo.User, mentions []*discordgo.User) {
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
  party, err := readFromDatabase("parties", player["Party"].(string))
  if err != nil {
    s.ChannelMessageSend(channelID, "`There was an error inviting to the Party`")
    return
  }
  partyMembers := party["Members"].(map[string]interface{})
  value, ok := partyMembers[player["ID"].(string)];  if value == false || ok == false {
    s.ChannelMessageSend(channelID, "`You are not the party leader`")
    return
  }
  partyMembers[mentions[0].ID] = false
  writeToDatabase("parties", party["ID"].(string), party)
  updated := map[string]string{"Party":party["ID"].(string)}
  updateToDatabase("players", mentions[0].ID, updated)
  s.ChannelMessageSend(channelID, "`Added "+mentions[0].Username+" to the '"+party["ID"].(string)+"' Party`")
}
