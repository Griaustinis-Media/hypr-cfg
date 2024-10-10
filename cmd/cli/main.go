package main

import (
  "os"
  "log"
)

func main() {
  args := os.Args[1:]

  var cfgPath string
  if len(args) > 0 {
    cfgPath = args[0]
  } else {
    cfgPath = "~/.config/hypr/hyprland.conf"
  }

  log.Printf("Loading config from path: '%s'", cfgPath)
}
