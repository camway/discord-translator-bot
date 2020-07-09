package lang

import (
  "fmt"
  "os"
  "context"
  "errors"
  "google.golang.org/api/option"
  translate "cloud.google.com/go/translate/apiv3"
  translatepb "google.golang.org/genproto/googleapis/cloud/translate/v3"
)

func NewTranslator(credentialsFile string) (translator *Translator, err error) {
  c, err := translate.NewTranslationClient(context.TODO(), option.WithCredentialsFile(credentialsFile))
  if err != nil {
    return nil, err
  }

  translator = &Translator{client: c}

  return
}

type Translator struct {
  client *translate.TranslationClient
}

func (t *Translator) TranslateMessage(message string, fromLanguageCode string, toLanguageCode string) (string, error) {
  req := &translatepb.TranslateTextRequest{
    Parent: fmt.Sprintf("projects/%s/locations/global", os.Getenv("GOOGLE_API_PROJECT_ID")),
    Contents: []string{message},
    SourceLanguageCode: fromLanguageCode,
    TargetLanguageCode: toLanguageCode,
    MimeType: "text/plain",
  }
  resp, err := t.client.TranslateText(context.TODO(), req)
  if err != nil {
    fmt.Println(err.Error())
    return "", err
  }
  for _, translation := range resp.GetTranslations() {
    return translation.GetTranslatedText(), nil
  }

  return "", errors.New("No translation found")
}
