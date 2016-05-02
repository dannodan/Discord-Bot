package main

import (
  // "fmt"
	// "encoding/json"
	"strconv"
	"github.com/bwmarrin/discordgo"
)

type Player struct {
	ID              string	`json: "id"`
	Name            string	`json: "name"`
	Level						string	`json: "level"`
	Experience			string  `json: "experience"`
	Next						string	`json: "next"`
	Health					string	`json: "health"`
	MaxHealth				string	`json: "maxHealth"`
	Strength        string	`json: "strength"`
	Intelligence    string	`json: "intelligence"`
	Vitality        string	`json: "vitality"`
	Spirit          string	`json: "spirit"`
	Dexterity       string	`json: "dexterity"`
	Agility         string	`json: "agility"`
	FreePoints			string	`json: "freePoints"`
  Party           string  `json: "party"`
}

func generatePlayer(s *discordgo.Session, channelID string, user *discordgo.User) {
  player := Player{user.ID, user.Username, "1", "0", "50", "20", "20", "0", "0" ,"0" ,"0" ,"0" ,"0", "8", ""}
	writeToDatabase("players", player.ID, player)
	message := 	"```Character Created:\nName: "+player.Name+"\tLvl: "+player.Level+
							" ("+player.Experience+"/"+player.Next+")\nHP: "+player.Health+
							" / "+player.MaxHealth+"\nSTR: "+player.Strength+
							"\tVIT: "+player.Vitality+"\nINT: "+player.Intelligence+
							"\tSPR: "+player.Spirit+"\nDEX: "+player.Dexterity+
							"\tAGI: "+player.Agility+"\nYou have "+player.FreePoints+" Stat points```"
	s.ChannelMessageSend(channelID, message)
}

func playerStats(s *discordgo.Session, channelID string, user *discordgo.User) {
	player, err := readFromDatabase("players", user.ID)
  if err != nil {
    s.ChannelMessageSend(channelID, "`There isn't a character with this ID, please use $generate first`")
		return
  }
  message :=	"```Status for:\nName: "+player["Name"].(string)+"\tLvl: "+player["Level"].(string)+
							" ("+player["Experience"].(string)+"/"+player["Next"].(string)+")\nHP: "+player["Health"].(string)+
							" / "+player["MaxHealth"].(string)+"\nSTR: "+player["Strength"].(string)+
              "\tVIT: "+player["Vitality"].(string)+"\nINT: "+player["Intelligence"].(string)+
              "\tSPR: "+player["Spirit"].(string)+"\nDEX: "+player["Dexterity"].(string)+
              "\tAGI: "+player["Agility"].(string)+"\nYou have "+player["FreePoints"].(string)+" Stat points```"

  s.ChannelMessageSend(channelID, message)
}

func allocateStatPoints(s *discordgo.Session, channelID string, user *discordgo.User, argument, quantity string) {
	quantityInt, err := strconv.Atoi(quantity)
	if err != nil {
		s.ChannelMessageSend(channelID, "`Invalid value for points, please enter a valid number`")
		return
	}
	player, err := readFromDatabase("players", user.ID)
	if err != nil {
    s.ChannelMessageSend(channelID, "`There isn't a character with this ID, please use $generate first`")
    return
  }
  freePointInt, err := strconv.Atoi(player["FreePoints"].(string))
	if quantityInt > freePointInt {
		s.ChannelMessageSend(channelID, "`Not enough Stat Points to allocate`")
		return
	}
	playerStatus, err := strconv.Atoi(player[argument].(string))
	if err != nil {
		s.ChannelMessageSend(channelID, "`Couldn't allocate points`")
		return
	}
  player[argument] = strconv.Itoa(playerStatus + quantityInt)
  player["FreePoints"] = strconv.Itoa(freePointInt - quantityInt)
  writeToDatabase("players", player["ID"].(string), player)
  message := 	"```Status for:\nName: "+player["Name"].(string)+"\tLvl: "+player["Level"].(string)+
              " ("+player["Experience"].(string)+"/"+player["Next"].(string)+")\nHP: "+player["Health"].(string)+
							" / "+player["MaxHealth"].(string)+"\nSTR: "+player["Strength"].(string)+
              "\tVIT: "+player["Vitality"].(string)+"\nINT: "+player["Intelligence"].(string)+
              "\tSPR: "+player["Spirit"].(string)+"\nDEX: "+player["Dexterity"].(string)+
              "\tAGI: "+player["Agility"].(string)+"\nYou have "+player["FreePoints"].(string)+" Stat points```"

  s.ChannelMessageSend(channelID, message)
}

func testing(user *discordgo.User)  {
  test := map[string]string{}
  test["Party"] = "Test"
  updateToDatabase("players", user.ID, test)
}
