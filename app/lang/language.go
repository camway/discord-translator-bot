package lang

import (
  "fmt"

  "golang.org/x/text/language"
  "golang.org/x/text/language/display"
)

type Language struct {
  Name string
  Code string
}

func (l *Language) Tag() language.Tag {
  return language.Make(l.Code)
}

func (l *Language) ToString() string {
  en := display.English.Tags()

  return fmt.Sprintf(
    "%-20s %-8s %s",
    en.Name(l.Tag()),
    fmt.Sprintf("(%s)", l.Code),
    display.Self.Name(l.Tag()),
  )
}
