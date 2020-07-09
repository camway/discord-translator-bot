package main

import (
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"
  "context"
  "sync"
  "errors"
  "strconv"
)

var dataManager Manager

const groupCollection = "groups"

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

type Channel struct {
  ID string `bson:"id"`
  Name string `bson:"name"`
  Language string `bson:"language"`
}

func (c *Channel) ToString() string {
  return c.Name + " (" + c.Language + ")"
}

func NewManager(client *mongo.Client, db string) {
  dataManager = Manager{
    mclient: client,
    database: db,
  }
  err := dataManager.Load()
  if err != nil {
    panic(err)
  }
}

type Manager struct {
  mclient *mongo.Client
  database string
  Groups []Group
  mutex sync.Mutex
}

func (m *Manager) GroupsInServer(id string) (groups []Group) {
  for _, g := range m.Groups {
    if g.ServerID == id {
      groups = append(groups, g)
    }
  }
  return
}

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

func (m *Manager) GetChannelByID(channelID string) *Channel {
  for _, g := range m.Groups {
    for _, c := range g.Channels {
      if c.ID == channelID {
        return &c
      }
    }
  }
  return nil
}

func (m *Manager) ChannelGroups(channel *Channel, serverID string) ([]Group, []string) {
  var groupsInServer = m.GroupsInServer(serverID)
  var groups []Group
  for _, g := range groupsInServer {
    for _, c := range g.Channels {
      if c.ID == channel.ID {
        groups = append(groups, g)
        break;
      }
    }
  }
  var languages = getLanguagesInGroups(groups, channel.Language)

  return groups, languages
}

func (m *Manager) Load() error {
  groups, err := m.GetGroups()
  if err != nil {
    return err
  }
  m.Groups = groups
  return nil
}

func (m *Manager) GetGroups() ([]Group, error) {
  col := m.mclient.Database(m.database).Collection(groupCollection)
  // Pass these options to the Find method
  findOptions := options.Find()
  findOptions.SetLimit(100)

  // Here's an array in which you can store the decoded documents
  var groups []Group

  // Passing bson.D{{}} as the filter matches all documents in the collection
  cur, err := col.Find(context.TODO(), bson.D{{}}, findOptions)
  if err != nil {
    return nil, err
  }
  defer cur.Close(context.TODO())

  // Finding multiple documents returns a cursor
  // Iterating through the cursor allows us to decode documents one at a time
  for cur.Next(context.TODO()) {
    // create a value into which the single document can be decoded
    var elem Group
    err := cur.Decode(&elem)
    if err != nil {
      return nil, err
    }

    groups = append(groups, elem)
  }

  if err := cur.Err(); err != nil {
    return nil, err
  }

  return groups, nil
}

func (m *Manager) GetGroupByName(name string) *Group {
  for _, group := range m.Groups {
    return &group
  }
  return nil
}

func (m *Manager) ListGroups() (out string) {
  if len(m.Groups) == 0 {
    out = "No groups created. Use: !create myGroup"
  }
  for _, group := range m.Groups {
    out += group.Name + "(" + strconv.Itoa(len(group.Channels)) + ")" + "\n"
    for _, channel := range group.Channels {
      out += "  " + channel.ToString() + "\n"
    }
  }
  return
}

func (m *Manager) createGroup(name string, serverID string) error {
  var group = Group{
    ID: primitive.NewObjectID(),
    Name: name,
    ServerID: serverID,
    Channels: []Channel{},
  }
  col := m.mclient.Database(m.database).Collection(groupCollection)
  _, err := col.InsertOne(context.TODO(), group)
  if err != nil {
    return err
  }
  m.Groups = append(m.Groups, group)
  m.Load()
  return nil
}

func (m *Manager) saveGroup(group *Group) error {
  col := m.mclient.Database(m.database).Collection(groupCollection)
  _, err := col.UpdateOne(
    context.TODO(),
    bson.M{"_id": group.ID},
    bson.M{"$set": group},
  )
  if err != nil {
    return err
  }
  m.Load()
  return nil
}

func (m *Manager) deleteGroup(groupID primitive.ObjectID) error {
  col := m.mclient.Database(m.database).Collection(groupCollection)
  _, err := col.DeleteOne(context.TODO(), bson.M{"_id": groupID})
  if err != nil {
    return err
  }

  var foundAt = -1
  for i := range m.Groups {
    if m.Groups[i].ID == groupID {
      foundAt = i
    }
  }
  m.Groups = append(m.Groups[foundAt:], m.Groups[foundAt+1:]...)
  m.Load()
  return nil
}

func (m *Manager) addChannelToGroup(groupName string, channel Channel) (err error) {
  group := m.GetGroupByName(groupName)
  if group == nil {
    return errors.New("Couldn't find group")
  }

  group.Channels = append(group.Channels, channel)

  err = m.saveGroup(group)
  if err != nil {
    return
  }
  m.Load()
  return
}

func (m *Manager) removeChannelFromGroup(groupName string, channelID string) (err error) {
  group := m.GetGroupByName(groupName)
  if group == nil {
    return errors.New("Couldn't find group")
  }

  var foundAt = -1
  for index, c := range group.Channels {
    if c.ID == channelID {
      foundAt = index
    }
  }

  if foundAt < 0 {
    return errors.New("Couldn't find channel")
  }

  group.Channels = append(group.Channels[:foundAt], group.Channels[foundAt+1:]...)

  err = m.saveGroup(group)
  if err != nil {
    return
  }
  m.Load()
  return
}
