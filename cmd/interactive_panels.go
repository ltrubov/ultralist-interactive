package cmd

import (
  // "unicode/utf8"

  //"github.com/mattn/go-runewidth"
  "github.com/nsf/termbox-go"
)

type Panel struct {
  text     string
  title    string
  zerox    int
  active   bool
}

var left_panel, center_panel, right_panel Panel
var full_panel Panel

func redraw_panels() {

  w, _ := termbox.Size()
  pw := w/3

  left_panel.SetupWith(0, "TODAY", "Batman")
  center_panel.SetupWith(pw, "THIS WEEK", "Superman")
  right_panel.SetupWith(2*pw, "THIS MONTH", "Flash")
  full_panel.SetupWith(0, "HELP", "Help")
  full_panel.active = false

  panels := [4]Panel{left_panel, right_panel, center_panel, full_panel}

  for _,p := range panels {
    if !p.active {
      p.Draw()
    }
  }
}

func draw_full_panel() {
  //w, _ := termbox.Size()
  full_panel.zerox = 0
}

func (p *Panel) SetupWith(zx int, tt string, textStr string) {
  p.zerox = zx
  p.title = tt
  p.text = textStr
  p.active = true
}

func (p *Panel) Draw() {
  const coldef = termbox.ColorDefault
  w, h := termbox.Size()
  pw := w/3
  tbprint(p.zerox + pw/2 - len(p.title)/2, 1, coldef, coldef, p.title)

  y := 3
  for y < h-3 {
    tbprint(p.zerox, y, coldef, coldef, p.text)
    y++
  }
}

