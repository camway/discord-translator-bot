package main

import (
  "fmt"
  "strings"
  "sync"

  "translatorbot/data"
  "translatorbot/lang"

  "github.com/bwmarrin/discordgo"
)

const discordMaxMessageLength = 2000

func sendMessage(s *discordgo.Session, channelID string, message string) (err error) {
  return sendMessageWithWrapper(s, channelID, message, "")
}

func sendMessageWithWrapper(s *discordgo.Session, channelID string, message string, wrapper string) (err error) {
  var (
    lines = strings.Split(message, "\n")
    messages []string

    msg string
  )

  for _, l := range lines {
    if len(msg) + len(l) + 1 + (2 * len(wrapper)) > discordMaxMessageLength {
      messages = append(messages, msg)
      msg = l + "\n"
    } else {
      msg = msg + l + "\n"
    }
  }
  if msg != "" {
    messages = append(messages, msg)
  }

  fmt.Println(fmt.Sprintf("Sending messages: %d", len(messages)))
  for _, m := range messages {
    fmt.Printf("Sending (%d): %s\n", len(m), m)
    _, err = s.ChannelMessageSend(channelID, fmt.Sprintf("%s%s%s", wrapper, m, wrapper))
    if err != nil {
      return err
    }
  }

  return
}

func translateMessageHandler(bot *Bot, s *discordgo.Session, m *discordgo.MessageCreate) (handled bool) {
  // Don't try to translate commands
  if strings.HasPrefix(m.Content, "!") {
    return false
  }
  handled = true
  var channel = bot.DataManager.GetChannelByID(m.ChannelID)
  if channel == nil {
    return
  }
  // groups, languages
  groups, languages := bot.DataManager.ChannelGroups(channel, m.GuildID)
  // Lang to Translation
  var translations = make(map[string]string)
  var wg sync.WaitGroup
  for _, lang := range languages {
    wg.Add(1)
    go func(l string) {
      var translation, err = bot.Translator.TranslateMessage(m.Content, channel.Language, l)
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

  var channels = make(map[string]data.Channel)
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
      sendMessage(s, c.ID, fmt.Sprintf("%s: %s", username, m.Content))
    } else {
      if val, ok := translations[c.Language]; ok {
        sendMessage(s, c.ID, fmt.Sprintf("%s: %s", username, val))
      } else {
        fmt.Println("Couldn't find translation for: ", c.Language)
      }
    }
  }

  return
}

const commandList = `!commands
!list languages
!list groups
!create GROUPNAME
!join GROUPNAME LANGUAGECODE
!leave GROUPNAME
!delete GROUPNAME`

// !commands
func commandListMessageHandler(bot *Bot, s *discordgo.Session, m *discordgo.MessageCreate) (handled bool) {
  if m.Content != "!commands" {
    return false
  }
  handled = true

  sendMessage(s, m.ChannelID, commandList)

  return
}

// !list languages
func languageListMessageHandler(bot *Bot, s *discordgo.Session, m *discordgo.MessageCreate) (handled bool) {
  fmt.Println("languageListMessageHandler - 1")
  if m.Content != "!list languages" {
    return false
  }
  fmt.Println("languageListMessageHandler - 2")
  handled = true

  fmt.Println("languageListMessageHandler - 3")

  fmt.Println(lang.LanguageList())
  err := sendMessageWithWrapper(s, m.ChannelID, lang.LanguageList(), "```")
  if err != nil {
    fmt.Println(err.Error())
  }

  fmt.Println("languageListMessageHandler - 4")

  return
}

// !list groups
func listGroupsMessageHandler(bot *Bot, s *discordgo.Session, m *discordgo.MessageCreate) (handled bool) {
  if m.Content != "!list groups" {
    return false
  }
  handled = true

  _,err := s.ChannelMessageSend(m.ChannelID, bot.DataManager.ListGroups())
  if err != nil {
    fmt.Println(err.Error())
  }

  return
}

// !create GROUPNAME
func createGroupMessageHandler(bot *Bot, s *discordgo.Session, m *discordgo.MessageCreate) (handled bool) {
  if !strings.HasPrefix(m.Content, "!create ") {
    return false
  }
  handled = true

  var parts = strings.SplitN(m.Content, " ", 2)
  if len(parts) == 1 {
    sendMessage(s, m.ChannelID, "Need a name for the group")
    return
  }

  var err = bot.DataManager.CreateGroup(strings.TrimSpace(parts[1]), m.GuildID)
  if err == nil {
    sendMessage(s, m.ChannelID, "Group created")
  } else {
    debugError(err)
    sendMessage(s, m.ChannelID, "Error occured while saving the group")
  }

  return
}

// !join GROUPNAME LANGUAGECODE
func joinGroupMessageHandler(bot *Bot, s *discordgo.Session, m *discordgo.MessageCreate) (handled bool) {
  if !strings.HasPrefix(m.Content, "!join ") {
    return false
  }
  handled = true

  var parts = strings.SplitN(m.Content, " ", 3)

  var group = strings.Trim(parts[1], " ")
  if len(group) == 0 {
    sendMessage(s, m.ChannelID, "Need a name for which group")
    return
  }

  if len(parts) != 3 {
    sendMessage(s, m.ChannelID, "Can't parse arguments")
    return
  }

  var code = strings.Trim(parts[2], " ")
  if len(code) == 0 {
    sendMessage(s, m.ChannelID, "Need the language code. Check: !list languages")
    return
  }

  var lang = lang.GetLanguage(code)
  if lang == nil {
    sendMessage(s, m.ChannelID, "Can't find the language code. Check: !list languages")
    return
  }

  var channel, err = s.Channel(m.ChannelID)
  if err != nil {
    debugError(err)
    sendMessage(s, m.ChannelID, "Error locating discord channel")
    return
  }

  err = bot.DataManager.AddChannelToGroup(
    group,
    data.Channel{
      ID: m.ChannelID,
      Name: channel.Name,
      Language: code,
    },
  )
  if err == nil {
    sendMessage(s, m.ChannelID, "Channel added")
  } else {
    debugError(err)
    sendMessage(s, m.ChannelID, "Error while adding the channel")
  }

  return
}

// !leave GROUPNAME
func leaveGroupMessageHandler(bot *Bot, s *discordgo.Session, m *discordgo.MessageCreate) (handled bool) {
  if !strings.HasPrefix(m.Content, "!leave ") {
    return false
  }
  handled = true
  // var contentLength = len(m.Content)
  var prefixLength = len("!leave ")

  var group = strings.Trim(m.Content[prefixLength:], " ")
  if len(group) == 0 {
    sendMessage(s, m.ChannelID, "Need a name for which group")
    return
  }

  err := bot.DataManager.RemoveChannelFromGroup(group, m.ChannelID)
  if err == nil {
    sendMessage(s, m.ChannelID, "Channel removed")
  } else {
    sendMessage(s, m.ChannelID, "Error while removing the channel")
  }

  return
}

// !delete GROUPNAME
func deleteGroupMessageHandler(bot *Bot, s *discordgo.Session, m *discordgo.MessageCreate) (handled bool) {
  if !strings.HasPrefix(m.Content, "!delete ") {
    return false
  }
  handled = true
  var prefixLength = len("!delete ")

  var group = strings.Trim(m.Content[prefixLength:], " ")
  if len(group) == 0 {
    sendMessage(s, m.ChannelID, "Need a name for which group")
    return
  }

  g := bot.DataManager.GetGroupByName(group)
  if g == nil {
    sendMessage(s, m.ChannelID, "Couldn't find the group")
  }

  err := bot.DataManager.DeleteGroup(g.ID)
  if err == nil {
    sendMessage(s, m.ChannelID, "Group deleted")
  } else {
    sendMessage(s, m.ChannelID, "Error while deleting the group")
  }

  return
}
