package main

import (
  "fmt"
  "github.com/bwmarrin/discordgo"
)

type Battle struct {
  ID            string `json: "id"`
  TurnOrder   []string `json: "turnorder"`
}

type Fighter struct {
  ID        string  `json: "id"`
  Speed     string  `json: "speed"`
}

func defineTurnOrder(party map[string]interface{}, enemies map[string]string) ([]Fighter) {
  turns := make([]Fighter, 0)
  for key := range party {
    fmt.Println(key)
    player, err := readFromDatabase("players", key)
    if err != nil {
      return nil
    }
    fmt.Println(player)
    turns = append(turns, Fighter{ ID : player["ID"].(string), Speed : player["Agility"].(string)})
  }
  for key := range enemies {
    enemy, err := readFromDatabase("enemies", key)
    if err != nil {
      return nil
    }
    turns = append(turns, Fighter{ ID : enemy["ID"].(string), Speed : enemy["Agility"].(string)})
  }
  return turns
}

func beginBattle(s *discordgo.Session, channelID string, user *discordgo.User)  {
  player, err := readFromDatabase("players", user.ID)
  if err != nil {
    s.ChannelMessageSend(channelID, "`You don't have a character, please use $generate first`")
    return
  }
  if player["Party"] == "" {

  } else {
    party, err := readFromDatabase("parties", player["Party"].(string))
    if err != nil {
      s.ChannelMessageSend(channelID, "`Error`")
      return
    }
    order := defineTurnOrder(party["Members"].(map[string]interface{}), map[string]string{"1":"Monster"})
    fmt.Println(order)
  }
  return
}
