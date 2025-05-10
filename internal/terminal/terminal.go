package terminal

import (
	"fmt"

	glamour "github.com/charmbracelet/glamour"
)

func PrintGlamourString(theString string) {
	//result := markdown.Render(theString, 80, 6)

	//result, err := glamour.Render(theString, "./cmd/styles/dark.json")
	result, err := glamour.Render(theString, "dracula")

	if err != nil {
		panic(err)
	}

	markdown := string(result)
	fmt.Println(markdown)
}
