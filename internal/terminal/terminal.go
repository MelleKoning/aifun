package terminal

import (
	"fmt"

	glamour "github.com/charmbracelet/glamour"
)

func PrintGlamourString(theString string) {
	//result := markdown.Render(theString, 80, 6)

	glamour.WithWordWrap(120)
	termRenderer, err := glamour.NewTermRenderer(glamour.WithWordWrap(120), glamour.WithStandardStyle("dracula"))
	if err != nil {
		fmt.Println("can not initialize termRenderer")
	}
	result, err := termRenderer.Render(theString)
	if err != nil {
		panic(err)
	}

	markdown := string(result)
	fmt.Print(markdown)
}
