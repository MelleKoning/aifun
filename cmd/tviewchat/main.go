package main

import (
	"context"
	"fmt"
	"log"

	"github.com/MelleKoning/aifun/internal/genaimodel"
	"github.com/MelleKoning/aifun/internal/terminal"
	"github.com/MelleKoning/aifun/internal/tviewview"
)

func main() {
	mdRenderer, err := terminal.New()
	if err != nil {
		fmt.Println(err)
		return
	}

	systemPrompt := `Act as a senior systems architect.
	Given the following requirements and constraints, walk through your thought process
	step-by-step for designing the backend of this platform. Then output a proposed high-level
	 architecture diagram in Markdown with labeled components.
	Requirements:
	- Create a backend for a chat application that uses the genai model as a backend.
	- The application should allow users to send messages to the model and receive responses.
	- The application should allow users to save and retrieve chat history.
	Constraints:
	- terminal console application written in golang
	- uses the tview package for rendering the UI
	- includes Glamour for rendering Markdown
	- includes the genaimodel package for talking to AI Models
	Tasks:
	- Ask progress and design of packages made so far
	- Help the developer improve and extend the software so far
	- Stay friendly and helpful
	- When having enough information for the architecture, stop and output the architecture diagram
`

	ctx := context.Background()
	modelAction, _ := genaimodel.NewModel(ctx, systemPrompt)

	// Create the console view
	tviewApp := tviewview.New(mdRenderer, modelAction)

	// Run the application
	if err := tviewApp.Run(); err != nil {
		log.Fatal(err)
	}
}
