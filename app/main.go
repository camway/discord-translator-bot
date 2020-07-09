package main

import (
  "fmt"
  "os"
  "os/signal"
  "syscall"

  "translatorbot/data"
  "translatorbot/db"
  "translatorbot/lang"
)

func debugError(err error) {
  fmt.Println("Error: " + err.Error())
}

func main() {
  var (
    dataManager *data.Manager
    translator *lang.Translator
    err error
  )

  dataManager, err = data.NewManager(
    db.Config{
      Host: os.Getenv("MONGO_HOST"),
      Database: os.Getenv("MONGO_DATABASE"),
      Username: os.Getenv("MONGO_USERNAME"),
      Password: os.Getenv("MONGO_PASSWORD"),
      Port: os.Getenv("MONGO_PORT"),
    },
  )
  if err != nil {
    panic(err)
  }
  defer dataManager.Close()

  translator, err = lang.NewTranslator("/endbot.json")
  if err != nil {
    panic(err)
  }

  bot := NewBot(
    os.Getenv("DISCORD_BOT_TOKEN"),
    dataManager,
    translator,
  )

  err = bot.Connect()
  if err != nil {
    panic(err)
  }
  defer bot.Close()


  sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
