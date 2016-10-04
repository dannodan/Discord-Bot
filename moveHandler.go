package main

import (
  // "fmt"
  "github.com/bwmarrin/discordgo"
  "github.com/Knetic/govaluate"
  "strconv"
  "math/rand"
)

type Move struct {
  ID      string `json: "id"`
  Damage  string `json: "damage"`
  Effect  string `json: "effect"`
}

func executeMove(s *discordgo.Session, channelID string, user *discordgo.User, move string, fighter map[string]interface{})  {
  var target map[string]interface{}
  player, _ := readFromDatabase("players", user.ID)
  if player["Party"] != nil {
    party, _ := readFromDatabase("parties", player["Party"].(string))
    members := party["Members"].(map[string]interface{})
    memberArray := make([]string, len(members))
    index := 0
    for k := range members {
      memberArray[index] = k
      index++
    }
    targetID := memberArray[rand.Intn(len(memberArray))]
    target, _ = readFromDatabase("players", targetID)
  } else {
    target = player
  }
  moveType, _ := readFromDatabase("moves", move)
  expression, _ := govaluate.NewEvaluableExpression(moveType["Damage"].(string))
  parameters := make(map[string]interface{})
  parameters["stat"], _ = strconv.Atoi(fighter["Strength"].(string))
  parameters["wAtk"] = 0
  damage, _ := expression.Evaluate(parameters)
  currentHealth, _ := strconv.Atoi(target["Health"].(string))
  newHealth := currentHealth - int(damage.(float64))
  s.ChannelMessageSend(channelID, "`"+fighter["Name"].(string)+" did "+strconv.Itoa(int(damage.(float64)))+" damage to "+target["Name"].(string)+"`")
  updateToDatabase("players", target["ID"].(string) , map[string]string{"Health":strconv.Itoa(newHealth)})
}
