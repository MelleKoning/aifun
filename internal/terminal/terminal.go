package terminal

import (
	"fmt"

	glamour "github.com/charmbracelet/glamour"
)

// Define ANSI escape codes for the desired color (e.g., green)
const colorGreen = "\033[32m"
const colorReset = "\033[0m"

// const colorYellow = "\033[33m"
const backGroundBlack = "\033[40m"

// const AttrReversed = "\033[7m"
const colorCyan = "\033[36m"

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

func PrintPrompt(historyLength int) {
	fmt.Printf("History items: %d\n", historyLength)
	fmt.Print(colorGreen + "('exit' to quit, `file` to upload, `prompt` to update systeminstruction) ")
	fmt.Println(colorCyan + backGroundBlack) // will be the typing colour
}

func PrintColourReset() {
	fmt.Print(colorReset)
}
