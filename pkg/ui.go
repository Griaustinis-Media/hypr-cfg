package pkg

import (
  "fmt"
  "strings"
  "github.com/rivo/tview"
  "github.com/gdamore/tcell/v2"
)

var selectedLine *RawLine

func buildMainFrame(tree *tview.TreeView) *tview.Frame {
	return tview.NewFrame(tree).AddText("HyprCfg", true, tview.AlignCenter, tcell.ColorWhite).AddText("Footer middle", false, tview.AlignCenter, tcell.ColorGreen)
}

func addBranch(name string, target *tview.TreeNode) *tview.TreeNode {
  node := tview.NewTreeNode(name).SetSelectable(true).SetExpanded(false)

  target.AddChild(node)

  return node
}

func addLeaf(ref *RawLine, key string, values []string, parent *tview.TreeNode) {
  var text string
  if len(values) > 0 {
    text = fmt.Sprintf("%s => '%s'", key, strings.Join(values, ", "))
  } else {
    text = key
  }
  node := tview.NewTreeNode(text).SetSelectable(true)
  node.SetReference(ref)
  parent.AddChild(node)
}

func openEditView(pager *tview.Pages, form *tview.Form, ref *RawLine) {
  // Form set value
  item := form.GetFormItemByLabel("Value")
  _, v := ref.Parse(false)
  item.(*tview.InputField).SetText(v.ToString())
  selectedLine = ref
  pager.ShowPage("edit")
}

func buildEditForm(form *tview.Form, pager *tview.Pages) {
  form.AddInputField("Value", "", 40, nil, nil)

  form.AddButton("Save", func() {
    item := form.GetFormItemByLabel("Value")
    text := item.(*tview.InputField).GetText()
    selectedLine.Update(text)
    pager.SwitchToPage("home")
  })

  form.AddButton("Cancel", func() {
    pager.SwitchToPage("home")
  })
}

func BuildApp(gconf GroupedConfig) *tview.Application {
  app := tview.NewApplication() 

  root := tview.NewTreeNode(gconf.Source)

  variablesBranch := addBranch("Variables", root)
  for _, v := range gconf.Variables {
    addLeaf(v.Ref, v.Key.Key, v.Value.Values, variablesBranch)
  }

  monitorsBranch := addBranch("Monitors", root)
  for _, v := range gconf.Monitors {
    addLeaf(v.Ref, v.Key.Key, v.Value.Values, monitorsBranch)
  }

  envBranch := addBranch("Env Variables", root)
  for _, v := range gconf.Env {
    addLeaf(v.Ref, v.Key.Key, v.Value.Values, envBranch)
  }

  autostartBranch := addBranch("Autostart", root)
  for _, v := range gconf.Autostart {
    addLeaf(v.Ref, v.Key.Key, v.Value.Values, autostartBranch)
  }

  bindingsBranch := addBranch("Bindings", root)
  for _, v := range gconf.Bindings {
    addLeaf(v.Ref, v.Key.Key, v.Value.Values, bindingsBranch)
  }

  tree := tview.NewTreeView().SetRoot(root).SetCurrentNode(root)

  frame := buildMainFrame(tree)
  form := tview.NewForm()
  pager := tview.NewPages().AddPage("home", frame, true, true).AddPage("edit", form, true, false)
  buildEditForm(form, pager)

  tree.SetSelectedFunc(func(node *tview.TreeNode) {
    ref := node.GetReference()
    if ref != nil {
      openEditView(pager, form, ref.(*RawLine))
    } else {
      node.SetExpanded(!node.IsExpanded())
    }
  })

  app.SetRoot(pager, true).EnableMouse(true).EnablePaste(true)
  return app 
}
