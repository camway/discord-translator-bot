package data

import (
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"
  "context"
  "sync"
  "time"
  "errors"
  "strconv"
  "translatorbot/db"
)

const groupCollection = "groups"

func NewManager(config db.Config) (dataManager *Manager, err error) {
  dataManager = &Manager{config: config}

  err = dataManager.connect()
  if err != nil {
    return
  }

  err = dataManager.Load()
  if err != nil {
    return
  }
  return
}

type Manager struct {
  config db.Config
  mclient *mongo.Client
  Groups []Group
  mutex sync.Mutex
}

func (m *Manager) connect() error {
  client, err := mongo.NewClient(options.Client().ApplyURI(m.config.ConnectionString()))
  if err != nil {
    return err
  }
  ctx, _ := context.WithTimeout(context.TODO(), 10*time.Second)
  err = client.Connect(ctx)
  if err != nil {
    return err
  }

  m.mclient = client

  return nil
}

func (m *Manager) Close() {
  m.mclient.Disconnect(context.TODO())
}

func (m *Manager) Load() error {
  groups, err := m.GetGroups()
  if err != nil {
    return err
  }
  m.Groups = groups
  return nil
}

func (m *Manager) GroupsInServer(id string) (groups []Group) {
  for _, g := range m.Groups {
    if g.ServerID == id {
      groups = append(groups, g)
    }
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

func (m *Manager) GetGroups() ([]Group, error) {
  col := m.mclient.Database(m.config.Database).Collection(groupCollection)
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

func (m *Manager) CreateGroup(name string, serverID string) error {
  var group = Group{
    ID: primitive.NewObjectID(),
    Name: name,
    ServerID: serverID,
    Channels: []Channel{},
  }
  col := m.mclient.Database(m.config.Database).Collection(groupCollection)
  _, err := col.InsertOne(context.TODO(), group)
  if err != nil {
    return err
  }
  m.Groups = append(m.Groups, group)
  m.Load()
  return nil
}

func (m *Manager) SaveGroup(group *Group) error {
  col := m.mclient.Database(m.config.Database).Collection(groupCollection)
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

func (m *Manager) DeleteGroup(groupID primitive.ObjectID) error {
  col := m.mclient.Database(m.config.Database).Collection(groupCollection)
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

func (m *Manager) AddChannelToGroup(groupName string, channel Channel) (err error) {
  group := m.GetGroupByName(groupName)
  if group == nil {
    return errors.New("Couldn't find group")
  }

  group.Channels = append(group.Channels, channel)

  err = m.SaveGroup(group)
  if err != nil {
    return
  }
  m.Load()
  return
}

func (m *Manager) RemoveChannelFromGroup(groupName string, channelID string) (err error) {
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

  err = m.SaveGroup(group)
  if err != nil {
    return
  }
  m.Load()
  return
}
