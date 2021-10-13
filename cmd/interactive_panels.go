package cmd

import (
  // "unicode/utf8"

  //"github.com/mattn/go-runewidth"
  "github.com/nsf/termbox-go"
  "github.com/ultralist/ultralist/ultralist"
)

// type Panel struct {
//   text     string
//   title    string
//   zerox    int
//   panwid   int
//   active   bool
// }

var left_panel, center_panel, right_panel ultralist.Panel
var full_panel ultralist.Panel

func redraw_panel_gallery() {

  w, _ := termbox.Size()
  pw := w/3

  for i,pt := range [3]string{"TODAY", "THIS WEEK", "THIS MONTH"} {
    panel := ultralist.NewPanel(pt, i*pw, pw-1)// ultralist.Panel{Zerox: i*pw, Panwid: pw, Title: pt}
    //panel.Draw("text of doom")
    ultralist.NewAppForPanel(true, panel).ListTodos("due:agenda group:project", true, true)
  }

  //full_panel = ultralist.Panel{Zerox: 0, Panwid: w, Title: "HELP"}
  //left_panel.SetupWith(0, pw, "TODAY", "Batman")
  // center_panel.SetupWith(pw, pw, "THIS WEEK", "Superman")
  // right_panel.SetupWith(2*pw, pw, "THIS MONTH", "Flash")
  // full_panel.SetupWith(0, w, "HELP", "Help")
  // full_panel.active = false

//   panels := [4]ultralist.Panel{left_panel, right_panel, center_panel, full_panel}
//
//   for _,p := range panels {
//     if !p.active {
//       p.Draw()
//     }
//   }
}

func redraw_full_panel() {
  w, _ := termbox.Size()

  panel := ultralist.Panel{Zerox: 0, Panwid: w, Title: "HELP"}
  panel.Draw("text of nice", termbox.ColorGreen)
}

// func (p *Panel) SetupWith(zx int, w int, tt string, textStr string) {
//   p.zerox = zx
//   p.panwid = w
//   p.title = tt
//   p.text = textStr
//   p.active = true
// }



