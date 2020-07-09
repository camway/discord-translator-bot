package main

import (
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
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

  pingMessageHandler(bot, s, m)

  commandListMessageHandler(bot, s, m)
  languageListMessageHandler(bot, s, m)

  listGroupsMessageHandler(bot, s, m)

  createGroupMessageHandler(bot, s, m)
  joinGroupMessageHandler(bot, s, m)
  leaveGroupMessageHandler(bot, s, m)
  deleteGroupMessageHandler(bot, s, m)

  translateMessageHandler(bot, s, m)
}
