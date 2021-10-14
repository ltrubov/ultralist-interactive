package ultralist

import (
	//"io"
	"sort"
	"strconv"
	"strings"
	//"text/tabwriter"
	"time"

	//"github.com/cheynewallace/tabby"
	//"github.com/fatih/color"

	"github.com/nsf/termbox-go"
	"github.com/mattn/go-runewidth"
	//"unicode/utf8"
)

type Panel struct {
	text     string
	Title    string
	Zerox    int
	Panwid   int
	cx int
	cy int
}

const coldef = termbox.ColorDefault

func NewPanel(tt string, zx, pw int) *Panel {
	p := &Panel{Zerox: zx, Panwid: pw, Title: tt}
	p.cx = 0
	p.cy = 2
	return p
}

func (p *Panel) Reset() {
	_, h := termbox.Size()

	fill(p.Zerox, 0, p.Panwid, h-3, termbox.Cell{Ch: ' '})
	tbprint(p.Zerox + p.Panwid/2 - len(p.Title)/2, 0, coldef, coldef, p.Title)
}

func (p *Panel) Draw(text string, color termbox.Attribute, sh bool) {
	for _, c := range text {
		rw := runewidth.RuneWidth(c)
		if c == '\n' {
			p.cy++
			p.cx = 0
		} else {
			if p.cx + rw >= p.Panwid {
				if sh {
					termbox.SetCell(p.Zerox + p.cx, p.cy, '…', color, coldef)
					// p.cy++
					// p.cx = 0
					break
				} else {
					termbox.SetCell(p.Zerox + p.cx, p.cy, '-', color, coldef)
					p.cy++
					p.cx = 2
				}
			}
			termbox.SetCell(p.Zerox + p.cx, p.cy, c, color, coldef)
			p.cx += rw
		}
	}
}

func (p *Panel) AddLine() {
	p.cy++
	p.cx = 0
}

func (p *Panel) DrawTab(text string, color termbox.Attribute, ts, tl int) {
	if ts + tl > p.Panwid {
		tl = p.Panwid - ts
	}

	var te = ts + tl
	for _, c := range text {
		rw := runewidth.RuneWidth(c)
		if ts + rw > te {
			termbox.SetCell(p.Zerox + ts, p.cy, '…', color, coldef)
			break
		}

		termbox.SetCell(p.Zerox + ts, p.cy, c, color, coldef)
		ts += rw
	}
}

// func (p Panel) Write(pa []byte) (n int, err error) {
// 	return len(pa),nil
// }

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}

func fill(x, y, w, h int, cell termbox.Cell) {
	for ly := 0; ly < h; ly++ {
		for lx := 0; lx < w; lx++ {
			termbox.SetCell(x+lx, y+ly, cell.Ch, cell.Fg, cell.Bg)
		}
	}
}

// PanelPrinter is the default struct of this file
type PanelPrinter struct {
	panel 					*Panel
	UnicodeSupport 	bool
	shorthand 			bool
}

// NewPanelPrinter creates a new screeen printer.
func NewPanelPrinter(unicodeSupport bool, p *Panel, sh bool) *PanelPrinter {
	formatter := &PanelPrinter{panel: p, UnicodeSupport: unicodeSupport, shorthand: sh}
	return formatter
}

const date_length = 12
const check_length = 5

// Print prints the output of ultralist to the panel.
func (f *PanelPrinter) Print(groupedTodos *GroupedTodos, printNotes bool, showStatus bool) {
	f.panel.Reset()

	var id_length = longestID(groupedTodos) + 2

	var keys []string
	for key := range groupedTodos.Groups {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	//tabby := tabby.NewCustom(tabwriter.NewWriter(f.panel, 0, 0, 2, ' ', 0))
	//tabby.AddLine()
	for _, key := range keys {
		//tabby.AddLine(cyan.Sprint(key))
		f.panel.Draw(key, termbox.ColorCyan, f.shorthand)
		f.panel.AddLine()
		for _, todo := range groupedTodos.Groups[key] {
			f.printTodo(todo, id_length)
		}
		//tabby.AddLine()
		f.panel.AddLine()
	}
	//tabby.Print()
}

func longestID(groupedTodos *GroupedTodos) int {
	var res,nl = 0,0
	for key := range groupedTodos.Groups {
		for _, todo := range groupedTodos.Groups[key] {
			nl = len(strconv.Itoa(todo.ID))
			if res < nl {
				res = nl
			}
		}
	}
	return res
}

func (f *PanelPrinter) printTodo(todo *Todo, id_length int) {
	var tdc = f.todoColor(todo)

	f.panel.DrawTab(strconv.Itoa(todo.ID), termbox.ColorYellow, 0, id_length)
	f.panel.DrawTab(f.formatCompleted(todo.Completed), termbox.ColorWhite, id_length, check_length)
	f.panel.DrawTab(f.formatDue(todo.Due), tdc, id_length + check_length, date_length)
	f.panel.cx = id_length + check_length + date_length
	f.panel.Draw(todo.Subject, tdc, f.shorthand)
	f.panel.AddLine()
}

func (f *PanelPrinter) todoColor(td *Todo) termbox.Attribute {
	if td.Completed {
		return termbox.ColorGreen
	} else {
		dueTime, _ := time.Parse(DATE_FORMAT, td.Due)
		if isPastDue(dueTime) {
			return termbox.ColorRed
		} else {
			return termbox.ColorBlue
		}
	}
}

func (f *PanelPrinter) formatID(ID int, isPriority bool) string {
	if isPriority {
		return yellowBold.Sprint(strconv.Itoa(ID))
	}
	return yellow.Sprint(strconv.Itoa(ID))
}

func (f *PanelPrinter) formatCompleted(completed bool) string {
	if completed {
		if f.UnicodeSupport {
			//return white.Sprint("[✔]")
			return "[✔]"
		}
		return "[x]"// white.Sprint("[x]")
	}
	return "[ ]"// white.Sprint("[ ]")
}

func (f *PanelPrinter) formatDue(due string) string {
	if due == "" {
		return "          "
	}
	dueTime, _ := time.Parse(DATE_FORMAT, due)
	return f.printRawDue(dueTime)

	// if isPriority {
	// 	return f.printPriorityDue(dueTime, completed)
	// }
	// return f.printDue(dueTime, completed)
}

func (f *PanelPrinter) formatStatus(status string, isPriority bool) string {
	if status == "" {
		return green.Sprint("          ")
	}

	if len(status) < 10 {
		for x := len(status); x <= 10; x++ {
			status += " "
		}
	}

	statusRune := []rune(status)

	if isPriority {
		return greenBold.Sprintf("%-10v", string(statusRune[0:10]))
	}
	return green.Sprintf("%-10s", string(statusRune[0:10]))
}

func (f *PanelPrinter) formatInformation(todo *Todo) string {
	var information []string
	if todo.IsPriority {
		information = append(information, "*")
	} else {
		information = append(information, " ")
	}
	if todo.HasNotes() {
		information = append(information, "N")
	} else {
		information = append(information, " ")
	}

	return white.Sprint(strings.Join(information, ""))
}

func (f *PanelPrinter) printRawDue(due time.Time) string {
	if isToday(due) {
		return "today     "
	} else if isTomorrow(due) {
		return "tomorrow  "
	}
	return due.Format("Mon Jan 02")
}

func (f *PanelPrinter) printDue(due time.Time, completed bool) string {
	if isToday(due) {
		return blue.Sprint("today     ")
	} else if isTomorrow(due) {
		return blue.Sprint("tomorrow  ")
	} else if isPastDue(due) && !completed {
		return red.Sprint(due.Format("Mon Jan 02"))
	}
	return blue.Sprint(due.Format("Mon Jan 02"))
}

func (f *PanelPrinter) printPriorityDue(due time.Time, completed bool) string {
	if isToday(due) {
		return blueBold.Sprint("today     ")
	} else if isTomorrow(due) {
		return blueBold.Sprint("tomorrow  ")
	} else if isPastDue(due) && !completed {
		return redBold.Sprint(due.Format("Mon Jan 02"))
	}
	return blueBold.Sprint(due.Format("Mon Jan 02"))
}

func (f *PanelPrinter) formatSubject(subject string, isPriority bool) string {
	splitted := strings.Split(subject, " ")

	if isPriority {
		return f.printPrioritySubject(splitted)
	}
	return f.printSubject(splitted)
}

func (f *PanelPrinter) printPrioritySubject(splitted []string) string {
	coloredWords := []string{}
	for _, word := range splitted {
		if projectRegex.MatchString(word) {
			coloredWords = append(coloredWords, magentaBold.Sprint(word))
		} else if contextRegex.MatchString(word) {
			coloredWords = append(coloredWords, redBold.Sprint(word))
		} else {
			coloredWords = append(coloredWords, whiteBold.Sprint(word))
		}
	}
	return strings.Join(coloredWords, " ")
}

func (f *PanelPrinter) printSubject(splitted []string) string {
	coloredWords := []string{}
	for _, word := range splitted {
		if projectRegex.MatchString(word) {
			coloredWords = append(coloredWords, magenta.Sprint(word))
		} else if contextRegex.MatchString(word) {
			coloredWords = append(coloredWords, red.Sprint(word))
		} else {
			coloredWords = append(coloredWords, white.Sprint(word))
		}
	}
	return strings.Join(coloredWords, " ")
}
