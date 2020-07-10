package lang

import (
  "strings"
)

func LanguageList() (list string) {
  for _, language := range languageCodes {
    list = list + language.ToString() + "\n"
  }
  return
}

func GetLanguage(str string) *Language {
  for _, lang := range languageCodes {
    if str == lang.Name || str == lang.Code || strings.HasPrefix(str, lang.Name) {
      return &lang
    }
  }
  return nil
}

var languageCodes []Language = []Language{
  {"Amharic","am"},
  {"Arabic","ar"},
  {"Basque","eu"},
  {"Bengali","bn"},
  {"Bulgarian","bg"},
  {"Catalan","ca"},
  {"Cherokee","chr"},
  {"Chinese-PRC","zh-CN"},
  {"Chinese-Taiwan","zh-TW"},
  {"Croatian","hr"},
  {"Czech","cs"},
  {"Danish","da"},
  {"Dutch","nl"},
  {"English-US","en"},
  {"English-UK","en-GB"},
  {"Estonian","et"},
  {"Filipino","fil"},
  {"Finnish","fi"},
  {"French","fr"},
  {"German","de"},
  {"Greek","el"},
  {"Gujarati","gu"},
  {"Hebrew","iw"},
  {"Hindi","hi"},
  {"Hungarian","hu"},
  {"Icelandic","is"},
  {"Indonesian","id"},
  {"Italian","it"},
  {"Japanese","ja"},
  {"Kannada","kn"},
  {"Korean","ko"},
  {"Latvian","lv"},
  {"Lithuanian","lt"},
  {"Malay","ms"},
  {"Malayalam","ml"},
  {"Marathi","mr"},
  {"Norwegian","no"},
  {"Polish","pl"},
  {"Portuguese-Brazil","pt-BR"},
  {"Portuguese-Portugal","pt-PT"},
  {"Romanian","ro"},
  {"Russian","ru"},
  {"Serbian","sr"},
  {"Slovak","sk"},
  {"Slovenian","sl"},
  {"Spanish","es"},
  {"Swahili","sw"},
  {"Swedish","sv"},
  {"Tamil","ta"},
  {"Telugu","te"},
  {"Thai","th"},
  {"Turkish","tr"},
  {"Urdu","ur"},
  {"Ukrainian","uk"},
  {"Vietnamese","vi"},
  {"Welsh","cy"},
}
