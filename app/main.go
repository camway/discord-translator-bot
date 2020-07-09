package main

import (
  "github.com/bwmarrin/discordgo"
  "fmt"
  "os"
  "os/signal"
  "syscall"
  "context"
  "google.golang.org/api/option"
  translate "cloud.google.com/go/translate/apiv3"
  translatepb "google.golang.org/genproto/googleapis/cloud/translate/v3"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func debugError(err error) {
  fmt.Println("Error: " + err.Error())
}

func watchMessages() {
  fmt.Println("watchMessages...")
  // Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func main() {
  var host = os.Getenv("MONGO_HOST")
  var db = os.Getenv("MONGO_DATABASE")
  var user = os.Getenv("MONGO_USERNAME")
  var pass = os.Getenv("MONGO_PASSWORD")
  var port = os.Getenv("MONGO_PORT")

  var connectionString = "mongodb://" + user + ":" + pass + "@" + host + ":" + port

  client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
  if err != nil {
    panic(err)
  }
  ctx, _ := context.WithTimeout(context.TODO(), 10*time.Second)
  err = client.Connect(ctx)
  if err != nil {
    panic(err)
  }
  NewManager(client, db)
  defer client.Disconnect(ctx)

  watchMessages()
}

func translateMessage(message string, fromLanguageCode string, toLanguageCode string) (string, error) {
  ctx := context.TODO()
  c, err := translate.NewTranslationClient(ctx, option.WithCredentialsFile("/endbot.json"))
  if err != nil {
    fmt.Println(err.Error())
    return "", err
  }

  req := &translatepb.TranslateTextRequest{
    Parent: fmt.Sprintf("projects/%s/locations/global", os.Getenv("GOOGLE_API_PROJECT_ID")),
    Contents: []string{message},
    SourceLanguageCode: fromLanguageCode,
    TargetLanguageCode: toLanguageCode,
    MimeType: "text/plain",
  }
  resp, err := c.TranslateText(ctx, req)
  if err != nil {
    fmt.Println(err.Error())
    return "", err
  }
  for _, translation := range resp.GetTranslations() {
    return translation.GetTranslatedText(), nil
  }

  fmt.Println("Fell off")
  return "", nil
}
