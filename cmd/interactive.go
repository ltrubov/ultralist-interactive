package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/ultralist/ultralist/ultralist"
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
		},
	}

	rootCmd.AddCommand(interactiveCmd)
	interactiveCmd.Flags().BoolVarP(&unicodeSupport, "unicode", "", true, "Allows unicode support in Ultralist output")
	interactiveCmd.Flags().BoolVarP(&colorSupport, "color", "", true, "Allows color in Ultralist output")
	interactiveCmd.Flags().BoolVarP(&listNotes, "notes", "", false, "Show a todo's notes when listing. ")
	interactiveCmd.Flags().BoolVarP(&showStatus, "status", "", false, "Show a todo's status")
}
