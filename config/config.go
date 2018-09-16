package config

import (
  "log"
  "github.com/BurntSushi/toml"
)

// DB server and credentials
type Config struct {
  Addrs []string
  Database string
  Username string
  Password string
}

// Read and pharse config file
func (c *Config) Read() {
  if _, err := toml.DecodeFile("config.toml", &c); err != nil {
    log.Fatal(err)
  }
}
