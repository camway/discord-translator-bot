package main

import (
  "fmt"

  "github.com/bwmarrin/discordgo"
  "translatorbot/data"
  "translatorbot/lang"
)

func NewBot(
  token string,
  dataManager *data.Manager,
  translator *lang.Translator,
) *Bot {
  return &Bot{
    token: token,
    DataManager: dataManager,
    Translator: translator,
  }
}

type Bot struct {
  dg *discordgo.Session
  token string
  DataManager *data.Manager
  Translator *lang.Translator
}

func (b *Bot) Connect() (err error) {
  dg, err := discordgo.New("Bot " + b.token)
	if err != nil {
		return
	}

  err = dg.Open()
  if err != nil {
    return
  }

  dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
    messageCreate(b, s, m)
  })

  b.dg = dg

  return
}

func (b *Bot) Close() {
  b.dg.Close()
}

func messageCreate(
  bot *Bot,
  s *discordgo.Session,
  m *discordgo.MessageCreate,
) {
  fmt.Println(fmt.Sprintf("Received: `%s`\n", m.Content))
	if m.Author.ID == s.State.User.ID {
		return
	}

  if handled := commandListMessageHandler(bot, s, m); handled {
  } else if handled := languageListMessageHandler(bot, s, m); handled {
  } else if handled := listGroupsMessageHandler(bot, s, m); handled {
  } else if handled := createGroupMessageHandler(bot, s, m); handled {
  } else if handled := joinGroupMessageHandler(bot, s, m); handled {
  } else if handled := leaveGroupMessageHandler(bot, s, m); handled {
  } else if handled := deleteGroupMessageHandler(bot, s, m); handled {
  } else if handled := translateMessageHandler(bot, s, m); handled {
  } else {
    fmt.Printf("Unhandled...\n")
  }
}
