package pkg

import (
  "github.com/rivo/tview"
  "github.com/gdamore/tcell/v2"
)

func buildMainFrame() *tview.Frame {
	return tview.NewFrame(tview.NewBox()).SetBorders(2, 2, 2, 2, 4, 4).AddText("HyprCfg", true, tview.AlignCenter, tcell.ColorWhite).AddText("Footer middle", false, tview.AlignCenter, tcell.ColorGreen)
}

func BuildApp(gconf GroupedConfig) *tview.Application {
  frame := buildMainFrame()
  app := tview.NewApplication() 
  app.SetRoot(frame, true).EnableMouse(true)
  return app 
}
