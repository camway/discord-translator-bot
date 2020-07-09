package data

type Channel struct {
  ID string `bson:"id"`
  Name string `bson:"name"`
  Language string `bson:"language"`
}

func (c *Channel) ToString() string {
  return c.Name + " (" + c.Language + ")"
}
