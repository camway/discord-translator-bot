package db

type Config struct {
  Host string
  Database string
  Username string
  Password string
  Port string
}

func (c *Config) ConnectionString() string {
  return "mongodb://" +
    c.Username +
    ":" +
    c.Password +
    "@" +
    c.Host +
    ":" +
    c.Port
}
