package main

import (
  "os"
  "log"

  "griaustinismedia/hypr-cfg/pkg"
)

func main() {
  args := os.Args[1:]

  var cfgPath string
  if len(args) > 0 {
    cfgPath = args[0]
  } else {
    log.Fatal("Please provide config path")
  }

  log.Printf("Loading config from path: '%s'", cfgPath)

  cfg, err := pkg.ReadConfig(cfgPath)
  if err != nil {
    log.Fatal(err)
  }

  grouped := pkg.BuildGroupedConfig(cfg)

  app := pkg.BuildApp(grouped)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
