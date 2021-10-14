package cmd

import (
	"strings"
  "strconv"

	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
  //"bufio"
  //"fmt"
  "os"
  //"os/exec"
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
			//ultralist.NewAppWithPrintOptions(unicodeSupport, colorSupport).ListTodos(strings.Join(args, " "), listNotes, showStatus)
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
  termboxSample2()
  os.Exit(0)
}

// func runCommand(commandStr string) error {
//   commandStr = strings.TrimSuffix(commandStr, "\n")
//   arrCommandStr := strings.Fields(commandStr)
//   switch arrCommandStr[0] {
//   case "exit", "quit":
//     os.Exit(0)
//     // add another case here for custom commands.
//   }
//   cmd := exec.Command(arrCommandStr[0], arrCommandStr[1:]...)
//   cmd.Stderr = os.Stderr
//   cmd.Stdout = os.Stdout
//   return cmd.Run()
// }

// func termboxSample() {
//   var i, j int
//   var fg, bg termbox.Attribute
//   var colorRange []termbox.Attribute = []termbox.Attribute{
//     termbox.ColorDefault,
//     termbox.ColorBlack,
//     termbox.ColorRed,
//     termbox.ColorGreen,
//     termbox.ColorYellow,
//     termbox.ColorBlue,
//     termbox.ColorMagenta,
//     termbox.ColorCyan,
//     termbox.ColorWhite,
//     termbox.ColorDarkGray,
//     termbox.ColorLightRed,
//     termbox.ColorLightGreen,
//     termbox.ColorLightYellow,
//     termbox.ColorLightBlue,
//     termbox.ColorLightMagenta,
//     termbox.ColorLightCyan,
//     termbox.ColorLightGray,
//   }
//
//   var row, col int
//   var text string
//   for i, fg = range colorRange {
//     for j, bg = range colorRange {
//       row = i + 1
//       col = j * 8
//       text = fmt.Sprintf(" %02d/%02d ", fg, bg)
//       tbprint2(col, row+0, fg, bg, text)
//       /*text = fmt.Sprintf(" on ")
//       tbprint2(col, row+1, fg, bg, text)
//       text = fmt.Sprintf(" %2d ", bg)
//       tbprint2(col, row+2, fg, bg, text)*/
//       //fmt.Println(text, col, row)
//     }
//   }
//   for j, bg = range colorRange {
//     tbprint2(j*8, 0, termbox.ColorDefault, bg, "       ")
//     tbprint2(j*8, i+2, termbox.ColorDefault, bg, "       ")
//   }
//
//   tbprint2(15, i+4, termbox.ColorDefault, termbox.ColorDefault,
//     "Press any key to close...")
//   termbox.Flush()
//   termbox.PollEvent()
//   termbox.Close()
// }

func termboxSample2() {
  var command_result = 0

  defer termbox.Close()
  termbox.SetInputMode(termbox.InputEsc)

  redraw_all()
mainloop:
  for {
    switch ev := termbox.PollEvent(); ev.Type {
    case termbox.EventKey:
      if edit_box.clear_if_need {
        edit_box.ClearBox()
        //continue
      }
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
        command_result = process_editbox_text()
        if command_result < 0 {
          break mainloop
        } else if command_result > 0 {
          edit_box.clear_if_need = true
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

func process_editbox_text() int {
  t := edit_box.TextString()
  trimmed_downcased := strings.ToLower(strings.TrimSpace(t))
  var res = 0
  switch trimmed_downcased {
    case "exit", "quit":
      return -1
    case "help":
      fpc.Active = true
      fpc.CurrCmd = trimmed_downcased
    case "glance":
      fpc.Active = false
    default:
      interpret_command(t)

      //commands that result in fullscreen display are excluded
      res = run_command_in_background(t)
    }

    edit_box.MoveCursorToBeginningOfTheLine()
    edit_box.DeleteTheRestOfTheLine()

    return res
}

func redraw_all() {
  const coldef = termbox.ColorDefault
  termbox.Clear(coldef, coldef)

  w, h := termbox.Size()
  if fpc.Active {
    redraw_full_panel()
  } else {
    redraw_panel_gallery()
  }

  redraw_edit_box(w,h,coldef)
  termbox.Flush()
}

func interpret_command(cmd string) {
  comps := Map(strings.Split(cmd, " "), strings.TrimSpace)
//   if strings.ToLower(comps[0]) == "uhelp" ||
//      Contains(comps, "-h") ||
//      Contains(comps, "--help") {
//
//     fpc.Active = true
//     fpc.CurrCmd = "uhelp"
//     fpc.CurrCmdArgs = strings.Join(RemoveHelpArgs(comps), " ")
//   } else
  if comps[0] == "l" || comps[0] == "list" {
    fpc.Active = true
    fpc.CurrCmd = "list"
    fpc.CurrCmdArgs = strings.Join(comps[1:], " ")
  } else if comps[0] == "version" {
    fpc.Active = true
    fpc.CurrCmd = "version"
    fpc.CurrCmdArgs = ""
  }
}

func run_command_in_background(cmd string) int {
  comps := Map(strings.Split(cmd, " "), strings.TrimSpace)
  cmd = comps[0]
  args := strings.Join(comps[1:], " ")
  app := ultralist.NewApp()
  switch cmd {
    case "add", "a":
      app.AddTodo(args)
    case "addnote", "an":
      todoID, _ := strconv.Atoi(comps[1])
      app.AddNote(todoID, strings.Join(comps[2:], " "))
    case "archive", "ar":
      app.ArchiveTodo(args)
    case "unarchive", "uar":
      app.UnarchiveTodo(args)
    case "auth":
      app.AuthWorkflow()
    case "complete", "c":
      app.CompleteTodo(args, false)
    case "deletenote", "dn":
      todoID, _ := strconv.Atoi(comps[1])
      noteID, _ := strconv.Atoi(comps[2])
      app.DeleteNote(todoID, noteID)
    case "delete", "d", "rm":
      app.DeleteTodo(args)
    case "editnote", "en":
      todoID, _ := strconv.Atoi(comps[1])
      noteID, _ := strconv.Atoi(comps[2])
      app.EditNote(todoID, noteID, strings.Join(comps[3:], " "))

    case "edit", "e":
      todoID, err := strconv.Atoi(comps[1])
      if err != nil {
        //fmt.Printf("Could not parse todo ID: '%s'\n", args[0])
        return -1
      }
      app.EditTodo(todoID, strings.Join(comps[2:], " "))
    case "init":
      app.InitializeRepo()
    case "prioritize", "p":
      app.PrioritizeTodo(args)
    case "status", "s":
      app.SetTodoStatus(args)
    case "web":
      app.OpenWeb()
    default:
      return 0
  }
  return 1
}






