package main

import (
  "github.com/bwmarrin/discordgo"
  "fmt"
  "strings"
  "sync"
)

/*
!list languages
!list groups
!create GROUPNAME
!join GROUPNAME LANGUAGECODE
!leave GROUPNAME
!delete GROUPNAME
*/

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
  fmt.Println("Message received: " + m.Content)

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
    fmt.Println("Early out")
		return
	}

  pingMessageHandler(s, m)

  commandListMessageHandler(s, m)
  languageListMessageHandler(s, m)

  listGroupsMessageHandler(s, m)

  createGroupMessageHandler(s, m)
  joinGroupMessageHandler(s, m)
  leaveGroupMessageHandler(s, m)
  deleteGroupMessageHandler(s, m)

  translateMessageHandler(s, m)
}

func translateMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
  // Don't try to translate commands
  if strings.HasPrefix(m.Content, "!") {
    return
  }
  var channel = dataManager.GetChannelByID(m.ChannelID)
  if channel == nil {
    return
  }
  // groups, languages
  groups, languages := dataManager.ChannelGroups(channel, m.GuildID)
  // Lang to Translation
  var translations = make(map[string]string)
  var wg sync.WaitGroup
  for _, lang := range languages {
    wg.Add(1)
    go func(l string) {
      fmt.Println(fmt.Sprintf("Getting translation for %s from %s to %s", m.Content, channel.Language, l))
      var translation, err = translateMessage(m.Content, channel.Language, l)
      if err == nil {
        translations[l] = translation
      } else {
        fmt.Println("Error while translating...")
        fmt.Println(err.Error())
      }
      wg.Done()
    }(lang)
  }
  wg.Wait()

  var channels = make(map[string]Channel)
  for _, g := range groups {
    for _, c := range g.Channels {
      channels[c.ID] = c
    }
  }

  var username string
  if m.Member.Nick != "" {
    username = m.Member.Nick
  } else {
    username = m.Author.Username
  }


  for _, c := range channels {
    if c.ID == channel.ID {
      continue
    } else if c.Language == channel.Language {
      s.ChannelMessageSend(c.ID, fmt.Sprintf("%s: %s", username, m.Content))
    } else {
      if val, ok := translations[c.Language]; ok {
        s.ChannelMessageSend(c.ID, fmt.Sprintf("%s: %s", username, val))
      } else {
        fmt.Println("Couldn't find translation for: ", c.Language)
      }
    }
  }
}

func pingMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
  // If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
    fmt.Println("sending Pong!")
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
    fmt.Println("sending Ping!")
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}

const commandList = `!commands
!list languages
!list groups
!create GROUPNAME
!join GROUPNAME LANGUAGECODE
!leave GROUPNAME
!delete GROUPNAME`

// !commands
func commandListMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
  if m.Content != "!commands" { return }

  s.ChannelMessageSend(m.ChannelID, commandList)
}

// !list languages
func languageListMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
  if m.Content != "!list languages" { return }

  s.ChannelMessageSend(m.ChannelID, LanguageList())
}

// !list groups
func listGroupsMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
  if m.Content != "!list groups" { return }

  s.ChannelMessageSend(m.ChannelID, dataManager.ListGroups())
}

// !create GROUPNAME
func createGroupMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
  if !strings.HasPrefix(m.Content, "!create ") { return }

  var parts = strings.SplitN(m.Content, " ", 2)
  if len(parts) == 1 {
    s.ChannelMessageSend(m.ChannelID, "Need a name for the group")
    return
  }

  var err = dataManager.createGroup(strings.TrimSpace(parts[1]), m.GuildID)
  if err == nil {
  s.ChannelMessageSend(m.ChannelID, "Group created")
  } else {
    debugError(err)
    s.ChannelMessageSend(m.ChannelID, "Error occured while saving the group")
  }
}

// !join GROUPNAME LANGUAGECODE
func joinGroupMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
  if !strings.HasPrefix(m.Content, "!join ") { return }
  // var contentLength = len(m.Content)
  // var prefixLength = len("!join ")

  var parts = strings.SplitN(m.Content, " ", 3)

  var group = strings.Trim(parts[1], " ")
  if len(group) == 0 {
    s.ChannelMessageSend(m.ChannelID, "Need a name for which group")
    return
  }

  if len(parts) != 3 {
    s.ChannelMessageSend(m.ChannelID, "Can't parse arguments")
    return
  }

  var code = strings.Trim(parts[2], " ")
  if len(code) == 0 {
    s.ChannelMessageSend(m.ChannelID, "Need the language code. Check: !list languages")
    return
  }

  var lang = GetLanguage(code)
  if lang == nil {
    s.ChannelMessageSend(m.ChannelID, "Can't find the language code. Check: !list languages")
    return
  }

  var channel, err = s.Channel(m.ChannelID)
  if err != nil {
    debugError(err)
    s.ChannelMessageSend(m.ChannelID, "Error locating discord channel")
    return
  }

  err = dataManager.addChannelToGroup(
    group,
    Channel{
      ID: m.ChannelID,
      Name: channel.Name,
      Language: code,
    },
  )
  if err == nil {
    s.ChannelMessageSend(m.ChannelID, "Channel added")
  } else {
    debugError(err)
    s.ChannelMessageSend(m.ChannelID, "Error while adding the channel")
  }
}

// !leave GROUPNAME
func leaveGroupMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
  if !strings.HasPrefix(m.Content, "!leave ") { return }
  // var contentLength = len(m.Content)
  var prefixLength = len("!leave ")

  var group = strings.Trim(m.Content[prefixLength:], " ")
  if len(group) == 0 {
    s.ChannelMessageSend(m.ChannelID, "Need a name for which group")
    return
  }

  err := dataManager.removeChannelFromGroup(group, m.ChannelID)
  if err == nil {
    s.ChannelMessageSend(m.ChannelID, "Channel removed")
  } else {
    s.ChannelMessageSend(m.ChannelID, "Error while removing the channel")
  }
}

// !delete GROUPNAME
func deleteGroupMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
  if !strings.HasPrefix(m.Content, "!delete ") { return }
  var prefixLength = len("!delete ")

  var group = strings.Trim(m.Content[prefixLength:], " ")
  if len(group) == 0 {
    s.ChannelMessageSend(m.ChannelID, "Need a name for which group")
    return
  }

  g := dataManager.GetGroupByName(group)
  if g == nil {
    s.ChannelMessageSend(m.ChannelID, "Couldn't find the group")
  }

  err := dataManager.deleteGroup(g.ID)
  if err == nil {
    s.ChannelMessageSend(m.ChannelID, "Group deleted")
  } else {
    s.ChannelMessageSend(m.ChannelID, "Error while deleting the group")
  }
}
