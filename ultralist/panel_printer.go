package ultralist

import (
	//"io"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/cheynewallace/tabby"
	//"github.com/fatih/color"

	"github.com/nsf/termbox-go"
	"github.com/mattn/go-runewidth"
	"unicode/utf8"
)

// var (
// 	blue        = color.New(0, color.FgBlue)
// 	blueBold    = color.New(color.Bold, color.FgBlue)
// 	green       = color.New(0, color.FgGreen)
// 	greenBold   = color.New(color.Bold, color.FgGreen)
// 	cyan        = color.New(0, color.FgCyan)
// 	cyanBold    = color.New(color.Bold, color.FgCyan)
// 	magenta     = color.New(0, color.FgMagenta)
// 	magentaBold = color.New(color.Bold, color.FgMagenta)
// 	red         = color.New(0, color.FgRed)
// 	redBold     = color.New(color.Bold, color.FgRed)
// 	white       = color.New(0, color.FgWhite)
// 	whiteBold   = color.New(color.Bold, color.FgWhite)
// 	yellow      = color.New(0, color.FgYellow)
// 	yellowBold  = color.New(color.Bold, color.FgYellow)
// )

type Panel struct {
	text     string
	Title    string
	Zerox    int
	Panwid   int
	active   bool
}

func (p *Panel) Draw(text string) {
	const coldef = termbox.ColorDefault
	_, h := termbox.Size()
	tbprint(p.Zerox + p.Panwid/2 - len(p.Title)/2, 1, coldef, coldef, p.Title)

	y := 3
	for y < h-3 {
		tbprint(p.Zerox, y, coldef, coldef, text)
		y++
	}
}

func (p Panel) Write(pa []byte) (n int, err error) {
	cx,cy := 0,3
	_, h := termbox.Size()
	const coldef = termbox.ColorDefault

	tbprint(p.Zerox + p.Panwid/2 - len(p.Title)/2, 1, coldef, coldef, p.Title)


	t := pa
	//lx := 0
	//tabstop := 0
	for {
		//rx := lx - eb.line_voffset
		if len(t) == 0 {
			break
		}

// 		if lx == tabstop {
// 			tabstop += tabstop_length
// 		}
//
// 		if rx >= w {
// 			termbox.SetCell(x+w-1, y, arrowRight,
// 				colred, coldef)
// 			break
// 		}

		r, size := utf8.DecodeRune(t)
		rw := runewidth.RuneWidth(r)
		if cx + rw > p.Panwid {
			cx = 2
			cy++
		}

		if r == '\n' {
			cy++
// 			for ; lx < tabstop; lx++ {
// 				rx = lx - eb.line_voffset
// 				if rx >= w {
// 					goto next
// 				}
//
// 				if rx >= 0 {
// 					termbox.SetCell(x+rx, y, ' ', coldef, coldef)
// 				}
// 			}
		} else if cy >= h-3 {
			break
		} else {
			termbox.SetCell(p.Zerox+cx, cy, r, coldef, coldef)
			cx += rw
		}
	//next:
		t = t[size:]
	}

	//tbprint(p.Zerox + cx, cy, coldef, coldef, string(pa))


	return len(pa),nil
}

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}

// PanelPrinter is the default struct of this file
type PanelPrinter struct {
	//Writer         *io.Writer
	panel Panel
	UnicodeSupport bool
}

// NewPanelPrinter creates a new screeen printer.
func NewPanelPrinter(unicodeSupport bool, p Panel) *PanelPrinter {
	// w := new(io.Writer)
	// formatter := &PanelPrinter{Writer: w, UnicodeSupport: unicodeSupport}
	formatter := &PanelPrinter{panel: p, UnicodeSupport: unicodeSupport}
	return formatter
}

// Print prints the output of ultralist to the panel.
func (f *PanelPrinter) Print(groupedTodos *GroupedTodos, printNotes bool, showStatus bool) {
	var keys []string
	for key := range groupedTodos.Groups {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	//tabby := tabby.NewCustom(tabwriter.NewWriter(color.Output, 0, 0, 2, ' ', 0))
	tabby := tabby.NewCustom(tabwriter.NewWriter(f.panel, 0, 0, 2, ' ', 0))
	tabby.AddLine()
	for _, key := range keys {
		tabby.AddLine(cyan.Sprint(key))
		for _, todo := range groupedTodos.Groups[key] {
			f.printTodo(tabby, todo, printNotes, showStatus)
		}
		tabby.AddLine()
	}
	tabby.Print()
}

func (f *PanelPrinter) printTodo(tabby *tabby.Tabby, todo *Todo, printNotes bool, showStatus bool) {
	if showStatus {
		tabby.AddLine(
			f.formatID(todo.ID, todo.IsPriority),
			f.formatCompleted(todo.Completed),
			f.formatDue(todo.Due, todo.IsPriority, todo.Completed),
			f.formatStatus(todo.Status, todo.IsPriority),
			f.formatSubject(todo.Subject, todo.IsPriority))
	} else {
		tabby.AddLine(
			f.formatID(todo.ID, todo.IsPriority),
			f.formatCompleted(todo.Completed),
			f.formatDue(todo.Due, todo.IsPriority, todo.Completed),
			f.formatStatus(todo.Status, todo.IsPriority),
			f.formatSubject(todo.Subject, todo.IsPriority))
	}

	if printNotes {
		for nid, note := range todo.Notes {
			tabby.AddLine(
				"  "+cyan.Sprint(strconv.Itoa(nid)),
				white.Sprint(""),
				white.Sprint(""),
				white.Sprint(""),
				white.Sprint(""),
				white.Sprint(note))
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
			return white.Sprint("[✔]")
		}
		return white.Sprint("[x]")
	}
	return white.Sprint("[ ]")
}

func (f *PanelPrinter) formatDue(due string, isPriority bool, completed bool) string {
	if due == "" {
		return white.Sprint("          ")
	}
	dueTime, _ := time.Parse(DATE_FORMAT, due)

	if isPriority {
		return f.printPriorityDue(dueTime, completed)
	}
	return f.printDue(dueTime, completed)
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
