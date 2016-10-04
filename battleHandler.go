package main

import (
  "fmt"
  "github.com/bwmarrin/discordgo"
  "sort"
  "strconv"
  "math/rand"
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

func defineTurnOrder(party map[string]interface{}, enemies map[string]interface{}) ([]Fighter) {
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
  writeToDatabase("battles", currentBattle.ID, currentBattle)
  player, _ := readFromDatabase("players", currentTurn.ID)
  if player == nil {
    enemy, _ := readFromDatabase("enemies", currentTurn.ID)
    s.ChannelMessageSend(channelID, "`It's "+enemy["Name"].(string)+"'s turn`")
    enemyMove(s, channelID, user, enemy)
  } else {
    s.ChannelMessageSend(channelID, "`It's "+player["Name"].(string)+"'s turn`")
  }
}

func beginBattle(s *discordgo.Session, channelID string, user *discordgo.User)  {
  var currentBattle Battle
  player, err := readFromDatabase("players", user.ID)
  if err != nil {
    s.ChannelMessageSend(channelID, "`You don't have a character, please use $generate first`")
    return
  }
  if player["Party"] == "" {
    currentBattle.ID = player["ID"].(string)
    fighterOrder := make(map[string]interface{})
    fighterOrder[player["ID"].(string)] = ""
    enemies, _ := readFromDatabase("partyEncounters", strconv.Itoa(rand.Intn(1)))
    currentBattle.TurnOrder = defineTurnOrder(fighterOrder, enemies)
  } else {
    party, err := readFromDatabase("parties", player["Party"].(string))
    if err != nil {
      s.ChannelMessageSend(channelID, "`Error`")
      return
    }
    enemies, _ := readFromDatabase("partyEncounters", strconv.Itoa(rand.Intn(1)))
    currentBattle.ID = party["ID"].(string)
    currentBattle.TurnOrder = defineTurnOrder(party["Members"].(map[string]interface{}), enemies)
  }
  getTurn(s, channelID, user, currentBattle)
  return
}

func enemyMove(s *discordgo.Session, channelID string, user *discordgo.User, enemy map[string]interface{}) {
  moves := enemy["Moves"].([]interface{})
  move := moves[rand.Intn(len(moves))].(string)
  s.ChannelMessageSend(channelID, "`"+enemy["Name"].(string)+" uses "+move+"`")
  executeMove(s, channelID, user, move, enemy)
}
