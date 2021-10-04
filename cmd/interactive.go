package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
  "bufio"
  "fmt"
  "os"
  "os/exec"
  "github.com/nsf/termbox-go"
)

func init() {
	var (
		unicodeSupport bool
		colorSupport   bool
		listNotes      bool
		showStatus     bool
		interactiveCmdDesc    = "Launch interactive version."
		interactiveCmdExample = `ultralist interactive`
		interactiveCmdLongDesc = `Launches an interactive version of the program, controllable internally.`
	)

	var interactiveCmd = &cobra.Command{
		Use:     "interactive",
		Aliases: []string{"i"},
		Example: interactiveCmdExample,
		Long:    interactiveCmdLongDesc,
		Short:   interactiveCmdDesc,
		Run: func(cmd *cobra.Command, args []string) {
			ultralist.NewAppWithPrintOptions(unicodeSupport, colorSupport).ListTodos(strings.Join(args, " "), listNotes, showStatus)
      interactiveRunLoop()
		},
	}

	rootCmd.AddCommand(interactiveCmd)
	interactiveCmd.Flags().BoolVarP(&unicodeSupport, "unicode", "", true, "Allows unicode support in Ultralist output")
	interactiveCmd.Flags().BoolVarP(&colorSupport, "color", "", true, "Allows color in Ultralist output")
	interactiveCmd.Flags().BoolVarP(&listNotes, "notes", "", false, "Show a todo's notes when listing. ")
	interactiveCmd.Flags().BoolVarP(&showStatus, "status", "", false, "Show a todo's status")
}

func interactiveRunLoop() {
  termbox.Init()
  //termboxSample()
  termboxSample2()

  reader := bufio.NewReader(os.Stdin)
  for {
    fmt.Print("$ ")
    cmdString, err := reader.ReadString('\n')
    if err != nil {
      fmt.Fprintln(os.Stderr, err)
    }
    err = runCommand(cmdString)
    if err != nil {
      fmt.Fprintln(os.Stderr, err)
    }
  }
}

func runCommand(commandStr string) error {
  commandStr = strings.TrimSuffix(commandStr, "\n")
  arrCommandStr := strings.Fields(commandStr)
  switch arrCommandStr[0] {
  case "exit", "quit":
    os.Exit(0)
    // add another case here for custom commands.
  }
  cmd := exec.Command(arrCommandStr[0], arrCommandStr[1:]...)
  cmd.Stderr = os.Stderr
  cmd.Stdout = os.Stdout
  return cmd.Run()
}

func termboxSample() {
  var i, j int
  var fg, bg termbox.Attribute
  var colorRange []termbox.Attribute = []termbox.Attribute{
    termbox.ColorDefault,
    termbox.ColorBlack,
    termbox.ColorRed,
    termbox.ColorGreen,
    termbox.ColorYellow,
    termbox.ColorBlue,
    termbox.ColorMagenta,
    termbox.ColorCyan,
    termbox.ColorWhite,
    termbox.ColorDarkGray,
    termbox.ColorLightRed,
    termbox.ColorLightGreen,
    termbox.ColorLightYellow,
    termbox.ColorLightBlue,
    termbox.ColorLightMagenta,
    termbox.ColorLightCyan,
    termbox.ColorLightGray,
  }

  var row, col int
  var text string
  for i, fg = range colorRange {
    for j, bg = range colorRange {
      row = i + 1
      col = j * 8
      text = fmt.Sprintf(" %02d/%02d ", fg, bg)
      tbprint2(col, row+0, fg, bg, text)
      /*text = fmt.Sprintf(" on ")
      tbprint2(col, row+1, fg, bg, text)
      text = fmt.Sprintf(" %2d ", bg)
      tbprint2(col, row+2, fg, bg, text)*/
      //fmt.Println(text, col, row)
    }
  }
  for j, bg = range colorRange {
    tbprint2(j*8, 0, termbox.ColorDefault, bg, "       ")
    tbprint2(j*8, i+2, termbox.ColorDefault, bg, "       ")
  }

  tbprint2(15, i+4, termbox.ColorDefault, termbox.ColorDefault,
    "Press any key to close...")
  termbox.Flush()
  termbox.PollEvent()
  termbox.Close()
}

func termboxSample2() {
  var keep_going = true

  defer termbox.Close()
  termbox.SetInputMode(termbox.InputEsc)

  redraw_all()
mainloop:
  for {
    switch ev := termbox.PollEvent(); ev.Type {
    case termbox.EventKey:
      switch ev.Key {
      case termbox.KeyEsc:
        break mainloop
      case termbox.KeyArrowLeft, termbox.KeyCtrlB:
        edit_box.MoveCursorOneRuneBackward()
      case termbox.KeyArrowRight, termbox.KeyCtrlF:
        edit_box.MoveCursorOneRuneForward()
      case termbox.KeyBackspace, termbox.KeyBackspace2:
        edit_box.DeleteRuneBackward()
      case termbox.KeyDelete, termbox.KeyCtrlD:
        edit_box.DeleteRuneForward()
      case termbox.KeyTab:
        edit_box.InsertRune('\t')
      case termbox.KeySpace:
        edit_box.InsertRune(' ')
      case termbox.KeyCtrlK:
        edit_box.DeleteTheRestOfTheLine()
      case termbox.KeyHome, termbox.KeyCtrlA:
        edit_box.MoveCursorToBeginningOfTheLine()
      case termbox.KeyEnd, termbox.KeyCtrlE:
        edit_box.MoveCursorToEndOfTheLine()
      case termbox.KeyEnter://, termbox.KeyReturn:
        keep_going = process_editbox_text()
        if !keep_going {
          break mainloop
        }
      default:
        if ev.Ch != 0 {
          edit_box.InsertRune(ev.Ch)
        }
      }
    case termbox.EventError:
      panic(ev.Err)
    }
    redraw_all()
  }
}

func process_editbox_text() bool {
  t := edit_box.TextString()
  trimmed_downcased := strings.ToLower(strings.TrimSpace(t))
  switch trimmed_downcased {
    case "exit", "quit":
      return false
      // add another case here for custom commands.
    default:
      edit_box.MoveCursorToBeginningOfTheLine()
      edit_box.DeleteTheRestOfTheLine()
    }
    return true

}


func tbprint2(x, y int, fg, bg termbox.Attribute, msg string) {
  for _, c := range msg {
    termbox.SetCell(x, y, c, fg, bg)
    x += 1
  }
}
