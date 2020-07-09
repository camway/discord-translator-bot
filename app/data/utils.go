package data

func getLanguagesInGroups(groups []Group, excludeLanguage string) (languageList []string) {
  var languages map[string]bool = make(map[string]bool)

  for _, g := range groups {
    for _, c := range g.Channels {
      if c.Language != excludeLanguage {
        languages[c.Language] = true
      }
    }
  }

  for k := range languages {
    languageList = append(languageList, k)
  }
  return
}
