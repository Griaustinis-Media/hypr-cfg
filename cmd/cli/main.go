package main

import (
  "os"
  "log"

  "github.com/rivo/tview"
  "github.com/gdamore/tcell/v2"

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

  app := tview.NewApplication().EnableMouse(true)
	frame := tview.NewFrame(tview.NewBox().SetBackgroundColor(tcell.ColorBlue)).
		SetBorders(2, 2, 2, 2, 4, 4).
		AddText("Header left", true, tview.AlignLeft, tcell.ColorWhite).
		AddText("Header middle", true, tview.AlignCenter, tcell.ColorWhite).
		AddText("Header right", true, tview.AlignRight, tcell.ColorWhite).
		AddText("Header second middle", true, tview.AlignCenter, tcell.ColorRed).
		AddText("Footer middle", false, tview.AlignCenter, tcell.ColorGreen).
		AddText("Footer second middle", false, tview.AlignCenter, tcell.ColorGreen)
	if err := app.SetRoot(frame, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

  cfg.Show()
}
