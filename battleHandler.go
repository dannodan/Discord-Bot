package main

import (
  "fmt"
  "github.com/bwmarrin/discordgo"
  "sort"
  "strconv"
)

type Battle struct {
  ID             string `json: "id"`
  Current        string `json: "current"`
  TurnOrder   []Fighter `json: "turnorder"`
}

type Fighter struct {
  ID        string  `json: "id"`
  Speed     string  `json: "speed"`
}

type BySpeed []Fighter

func (a BySpeed) Len() int           { return len(a) }
func (a BySpeed) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BySpeed) Less(i, j int) bool { first, _ := strconv.Atoi(a[i].Speed)
                                       scnd, _ := strconv.Atoi(a[j].Speed)
                                       return first > scnd }

func defineTurnOrder(party map[string]interface{}, enemies map[string]string) ([]Fighter) {
  turns := make([]Fighter, 0)
  for key := range party {
    player, err := readFromDatabase("players", key)
    if err != nil {
      return nil
    }
    turns = append(turns, Fighter{ ID : player["ID"].(string), Speed : player["Agility"].(string)})
  }
  for key := range enemies {
    enemy, err := readFromDatabase("enemies", key)
    if err != nil {
      return nil
    }
    turns = append(turns, Fighter{ ID : enemy["ID"].(string), Speed : enemy["Agility"].(string)})
  }
  sort.Sort(BySpeed(turns))
  return turns
}

func getTurn(s *discordgo.Session, channelID string, user *discordgo.User, currentBattle Battle) {
  currentTurn := currentBattle.TurnOrder[0]
  currentBattle.TurnOrder = append(currentBattle.TurnOrder[1:], currentBattle.TurnOrder[0])
  currentBattle.Current = currentTurn.ID
  player, _ := readFromDatabase("players", currentTurn.ID)
  s.ChannelMessageSend(channelID, "`It's "+player["Name"].(string)+"'s turn`")
  if player["Party"] == "" {
    
  } else {
    writeToDatabase("battles", player["Party"].(string), currentBattle)
  }
  fmt.Println(currentBattle.TurnOrder)
}

func beginBattle(s *discordgo.Session, channelID string, user *discordgo.User)  {
//   order := make([]Fighter, 0)
  var currentBattle Battle
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
    currentBattle.ID = party["ID"].(string)
    currentBattle.TurnOrder = defineTurnOrder(party["Members"].(map[string]interface{}), map[string]string{"1":"Monster"})
  }
  getTurn(s, channelID, user, currentBattle)
  return
}
