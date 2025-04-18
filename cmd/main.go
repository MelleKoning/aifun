package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

const (
	modelName         = "gemini-2.0-flash"
	systemInstruction = `
		You are King Julian from Penguins of Madagascar. Your name is Julian.
	`
)

func main() {
	ctx := context.Background()
	client, err := initializeClient(ctx)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	defer client.Close()

	model := setupModel(client)
	interactiveSession(ctx, model)
}

func initializeClient(ctx context.Context) (*genai.Client, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	return genai.NewClient(ctx, option.WithAPIKey(apiKey))
}

func setupModel(client *genai.Client) *genai.GenerativeModel {
	model := client.GenerativeModel(modelName)
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(systemInstruction)},
	}
	return model
}

func interactiveSession(ctx context.Context, model *genai.GenerativeModel) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter your prompt (or type 'exit' to quit): ")
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

		generateAndPrintResponse(ctx, model, prompt)
	}
}

func generateAndPrintResponse(ctx context.Context, model *genai.GenerativeModel, prompt string) {
	iter := model.GenerateContentStream(ctx, genai.Text(prompt))
	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		printResponse(resp)
	}
}

func printResponse(resp *genai.GenerateContentResponse) {
	fmt.Print(resp.Candidates[0].Content.Parts[0])
}
