package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/MelleKoning/aifun/internal/genainterface"
	"github.com/MelleKoning/aifun/internal/prompts"
	"github.com/MelleKoning/aifun/internal/terminal"
)

func main() {
	terminal.PrintGlamourString(`
# Welcome to diffreviewer - genai!

Select a prompt to use for judging the gitdiff.txt

> Note: this uses the successor of generative-ai-go which is "google.golang.org/genai"

~~~golang
func main() {
   fmt.Println("Hello world, rendertest")
}
~~~
	`)

	ctx := context.Background()
	var err error

	systemInstruction := selectAPrompt()
	modelAction, err := genainterface.NewModel(ctx, systemInstruction)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	//request.filePart, _ = addAFile(ctx, request.client)
	interactiveSession(ctx, modelAction)
}

func selectAPrompt() string {
	reader := bufio.NewReader(os.Stdin)

	// Define a list of prompts
	prompts := prompts.PromptList

	// Display the list of prompts
	terminal.PrintGlamourString("Select a prompt by entering the corresponding number:")

	var promptStrings strings.Builder
	for title, prompt := range prompts {
		promptStrings.WriteString(fmt.Sprintf("%d. %s\n", title+1, prompt.Name))
	}
	terminal.PrintGlamourString(promptStrings.String())

	// Read the user's selection
	fmt.Print("Enter your choice: ")
	choiceStr, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return ""
	}

	// Convert the choice to an integer
	var choice int
	_, err = fmt.Sscanf(choiceStr, "%d", &choice)
	if err != nil {
		fmt.Println("Could not scan input:", err)
		return ""
	}

	// Validate the choice
	if choice < 1 || choice > len(prompts) {
		fmt.Println("Invalid choice. Exiting...")
		return ""
	}

	// Use the selected prompt
	selectedPrompt := prompts[choice-1].Prompt
	terminal.PrintGlamourString(fmt.Sprintf("You selected: %s\n", selectedPrompt))

	return selectedPrompt
}

func interactiveSession(ctx context.Context, modelAction genainterface.Action) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("('exit' to quit, `file` to upload): ")
		prompt, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}
		// Trim the newline character from the input
		prompt = prompt[:len(prompt)-1]
		if prompt == "exit" {
			fmt.Println("Exiting...")
			break
		}

		if prompt == "file" {
			err := modelAction.ReviewFile()
			if err != nil {
				fmt.Println(err)
			}
			continue
		}

		modelAction.ChatMessage(prompt)
	}
}
