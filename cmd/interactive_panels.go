package cmd

import (
  "time"
  "fmt"

  "github.com/nsf/termbox-go"
  "github.com/ultralist/ultralist/ultralist"
)

//var left_panel, center_panel, right_panel ultralist.Panel

var gallery_titles = [3]string{"TODAY", "THIS WEEK", "THIS MONTH"}
var gallery_cmds = [3]string{
  "due:agenda group:project",
  "dueafter:tod duebefore:sun group:project",
  fmt.Sprintf("duebefore:%s01 dueafter:sat group:project", next_month_str())}


// var full_panel ultralist.Panel

func next_month_str() string {
  return time.Now().AddDate(0,1,0).Month().String()[:3]
}

func redraw_panel_gallery() {

  w, h := termbox.Size()
  pw := w/3

  //separate columns
  var i = 0
  for i < h-3 {
    termbox.SetCell(pw-1, i, '|', termbox.ColorDarkGray, termbox.ColorDefault)
    termbox.SetCell(2*pw-1, i, '|', termbox.ColorDarkGray, termbox.ColorDefault)
    i++
  }

  for i,pt := range gallery_titles {
    panel := ultralist.NewPanel(pt, i*pw, pw-1)
    ultralist.NewAppForPanel(true, panel, true).ListTodos(gallery_cmds[i], true, true)
  }
}

type FullPanelConfig struct {
  CurrCmd      string
  LastCmd      string
  Active       bool
}

var fpc = FullPanelConfig{CurrCmd: "", LastCmd: "", Active: false}

const helptext = `
This is the interactive version of the ultralist tool.

It allows one to see the tasks for today, the current week, and the current month at a glance.

The three groups of tasks are displayed in three panels across the screen. At the bottom, the edit box accepts

a) regular ultralist commands (other than help)
b) interactive commands

The interactive commands include:
a) 'quit' and 'exit' to return to the shell
b) 'glance' to go into glance (3 panel) mode
c) 'help' to see this help

Ultralist commands modifying the todos state (adding, completing, deleting, renaming tasks) will redisplay the current
state with data updates.

The 'uhelp' command will run ultralist help, and display the help in a single panel covering the entire width of the
terminal. The 'glance' command will return to the three-panel state.

Running an ultralist 'list' command will display the results of that command in a single screenwide panel. Afterwards,
modification commands will result in this command being re-run, redisplaying the tasks in the fullscreen panel. Run
the 'glance' command to return to the 3-panel view.
`

func redraw_full_panel() {
  w, _ := termbox.Size()

  panel := ultralist.NewPanel("HELP", 0, w)
  panel.Reset()
  panel.Draw(helptext, termbox.ColorDefault, true)
}

