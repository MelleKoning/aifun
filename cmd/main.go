package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()
	// You can create your own API key to connect with the google
	// models on the quickstart website: https://ai.google.dev/gemini-api/docs/quickstart?lang=go
	apiKey := os.Getenv("GEMINI_API_KEY")
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		fmt.Println("Error creating client:", err)
		return
	}
	defer client.Close()

	// select the model
	model := client.GenerativeModel("gemini-2.0-flash") // Or "gemini-pro"

	prompt := "Tell me, Maurice, why I should get your food, as I don´t want to."

	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(`
		  You are King Julian from Pinguins of Madagasar. Your name is Julian.
		`)},
	}

	// Use streaming output to be able to print results as soon as they arrive
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
	/*for idx, c := range resp.Candidates {
		for pidx, part := range c.Content.Parts {
			fmt.Printf("candidate %d, content %d\n", idx, pidx)
			fmt.Println(part)
		}
	}*/
}
