package main

import (
	"encoding/json"
	"strconv"
	"github.com/bwmarrin/discordgo"
)

type Player struct {
	ID              string	`json: "id"`
	Name            string	`json: "name"`
	Level						string	`json: "level"`
	Experience			string  `json: "experience"`
	Next						string	`json: "next"`
	Strength        string	`json: "strength"`
	Intelligence    string	`json: "stringligence"`
	Dexterity       string	`json: "dexterity"`
	Charisma        string	`json: "charisma"`
	Wisdom          string	`json: "wisdom"`
	Luck            string	`json: "luck"`
	FreePoints			string	`json: "freePoints"`
}

func playerStats(s *discordgo.Session, channelID string, user *discordgo.User) {
  message := ""
	dbread, err := readFromDatabase("players", user.ID)
  if err != nil {
    message = "`There isn't a character with this ID, please use $generate first`"
  } else {
    player := dbread.(map[string]interface{})
    message = 	"```Status for:\nName: "+player["Name"].(string)+"\tLvl: "+player["Level"].(string)+
                " ("+player["Experience"].(string)+"/"+player["Next"].(string)+")\nSTR: "+player["Strength"].(string)+
                "\tCHA: "+player["Charisma"].(string)+"\nINT: "+player["Intelligence"].(string)+
                "\tWIS: "+player["Wisdom"].(string)+"\nDEX: "+player["Dexterity"].(string)+
                "\tLUK: "+player["Luck"].(string)+"\nYou have "+player["FreePoints"].(string)+" Stat points```"
  }
  s.ChannelMessageSend(channelID, message)
}

func generatePlayer(s *discordgo.Session, channelID string, user *discordgo.User) {
  player := Player{user.ID, user.Username, "1", "0", "50", "0", "0" ,"0" ,"0" ,"0" ,"0", "8"}
	writeToDatabase("players", player.ID, player)
	message := 	"```Character Created:\nName: "+player.Name+"\tLvl: "+player.Level+
							" ("+player.Experience+"/"+player.Next+")\nSTR: "+player.Strength+
							"\tCHA: "+player.Charisma+"\nINT: "+player.Intelligence+
							"\tWIS: "+player.Wisdom+"\nDEX: "+player.Dexterity+
							"\tLUK: "+player.Luck+"\nYou have "+player.FreePoints+" Stat points```"
	s.ChannelMessageSend(channelID, message)
}

func updatePlayer(s *discordgo.Session, channelID string, user *discordgo.User, argument, quantity string) {
	updatedPlayer := Player{}
	quantityInt, _ := strconv.Atoi(quantity)
	dbread, err := readFromDatabase("players", user.ID)
  if err != nil {
    s.ChannelMessageSend(channelID, "`There isn't a character with this ID, please use $generate first`")
    return
  }
  player := dbread.(map[string]interface{})
  freePointInt, err := strconv.Atoi(player["FreePoints"].(string))
  if quantityInt <= freePointInt {
    player[argument] = quantity
    freePointInt = freePointInt - quantityInt
    player["FreePoints"] = strconv.Itoa(freePointInt)
  } else {
    s.ChannelMessageSend(channelID, "`Not enough Stat Points to allocate`")
    return
  }
  aux, err := json.Marshal(player)
  if err != nil {
    panic(err)
  }
  if err := json.Unmarshal(aux, &updatedPlayer); err != nil {
    panic(err)
  }
  writeToDatabase("players", player["ID"].(string), updatedPlayer)
  message := 	"```Status for:\nName: "+player["Name"].(string)+"\tLvl: "+player["Level"].(string)+
              " ("+player["Experience"].(string)+"/"+player["Next"].(string)+")\nSTR: "+player["Strength"].(string)+
              "\tCHA: "+player["Charisma"].(string)+"\nINT: "+player["Intelligence"].(string)+
              "\tWIS: "+player["Wisdom"].(string)+"\nDEX: "+player["Dexterity"].(string)+
              "\tLUK: "+player["Luck"].(string)+"\nYou have "+player["FreePoints"].(string)+" Stat points```"

  s.ChannelMessageSend(channelID, message)
}