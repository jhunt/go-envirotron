package main

import (
  "fmt"
  env "github.com/jhunt/go-envirotron"
)

type Config struct {
  URL      string `env:"THING_URL"`
  Username string `env:"THING_USERNAME"`
  Password string `env:"THING_PASSWORD"`
}

func main() {
  c := Config{}
  env.Override(&c)

  fmt.Printf("connecting to %s, as %s\n", c.URL, c.Username)
}
