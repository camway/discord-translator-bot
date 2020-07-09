package data

import (
  "go.mongodb.org/mongo-driver/bson/primitive"
)

type Group struct {
  ID primitive.ObjectID `bson:"_id"`
  Name string `bson:"name"`
  ServerID string `bson:"serverId"`
  Channels []Channel `bson:"channels"`
}

func (g *Group) ToString() string {
  var output = g.Name + "\n"
  for index := range g.Channels {
    output = output + "  " + g.Channels[index].Name + "\n"
  }
  return output
}
