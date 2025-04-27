package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	glamour "github.com/charmbracelet/glamour"
)

const (
	modelName    = "gemini-2.0-flash"
	julianPrompt = "You are King Julian from Penguins of Madagascar. Your name is Julian."

	gitReviewPrompt = `You are an expert developer and git super user. You do code reviews based on the git diff output between two commits.
	* The diff contains a few unchanged lines of code. Focus on the code that changed. Changed are added and removed lines.
	* The added lines start with a "+" and the removed lines that start with a "-"
	Complete the following tasks, and be extremely critical and precise in your review:
	* [Description] Describe the code change.
	* [Obvious errors] Look for obvious errors in the code and suggest how to fix.
	* [Improvements] Suggest improvements where relevant. Suggestions must be rendered as code, not as diff.
	* [Friendly advice] Give some friendly advice or heads up where relevant.
	* [Stop when done] Stop when you are done with the review.

	This is the git diff output between two commits: \n\n {diff}

	AI OUTPUT:`
)

type Request struct {
	client   *genai.Client
	model    *genai.GenerativeModel
	textPart genai.Part
	filePart genai.Part
}

func main() {
	printGlamourString(`
# Welcome to aifun!

You first have to choose a prompt.

> This is a quote

This is some rendered code:

~~~golang
func main() {
   fmt.Println("hello")
}
~~~

That was the markdown rendering test
	`)
	request := new(Request)
	ctx := context.Background()
	var err error

	request.client, err = initializeClient(ctx)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	defer func() {
		err := request.client.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	prompt := selectAPrompt()
	request.model = setupModel(request.client, prompt)
	request.filePart = addAFile(ctx, request.client)
	interactiveSession(ctx, request)
}

func initializeClient(ctx context.Context) (*genai.Client, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	return genai.NewClient(ctx, option.WithAPIKey(apiKey))
}

func selectAPrompt() string {
	reader := bufio.NewReader(os.Stdin)

	// Define a list of prompts
	prompts := []string{
		julianPrompt,
		gitReviewPrompt,
	}

	// Display the list of prompts
	printGlamourString("Select a prompt by entering the corresponding number:")
	for i, prompt := range prompts {
		printGlamourString(fmt.Sprintf("%d. %s\n", i+1, prompt))
	}

	// Read the user's selection
	printGlamourString("Enter your choice: ")
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
	selectedPrompt := prompts[choice-1]
	printGlamourString(fmt.Sprintf("You selected: %s\n", selectedPrompt))

	return prompts[choice-1]
}

func setupModel(client *genai.Client, systemInstruction string) *genai.GenerativeModel {
	model := client.GenerativeModel(modelName)
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(systemInstruction)},
	}
	return model
}

func interactiveSession(ctx context.Context, request *Request) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("('exit' to quit, `file` to upload): ")
		prompt, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}
		// reset filepart
		request.filePart = nil

		// Trim the newline character from the input
		prompt = prompt[:len(prompt)-1]

		request.textPart = genai.Text(prompt)

		if prompt == "exit" {
			fmt.Println("Exiting...")
			break
		}

		if prompt == "file" {
			request.filePart = addAFile(ctx, request.client)
		}

		generateAndPrintResponse(ctx, request)
	}
}

func generateAndPrintResponse(ctx context.Context, request *Request) {
	var iter *genai.GenerateContentResponseIterator
	if request.filePart != nil {
		iter = request.model.GenerateContentStream(ctx, request.textPart, request.filePart)
	} else {
		iter = request.model.GenerateContentStream(ctx, request.textPart)

	}
	//iter := request.model.GenerateContentStream(ctx, request.textPart, request.filePart)
	var allparts []genai.Part
	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		printResponse(resp)
		allparts = append(allparts, resp.Candidates[0].Content.Parts[0])

	}

	printGlamour(allparts)
}

func printResponse(resp *genai.GenerateContentResponse) {
	result := resp.Candidates[0].Content.Parts[0]

	fmt.Print(result)

}

func printGlamour(resp []genai.Part) {
	var build strings.Builder
	for _, p := range resp {

		build.WriteString(fmt.Sprintf("%v", p))
	}
	printGlamourString(build.String())
}

func printGlamourString(theString string) {
	//result := markdown.Render(theString, 80, 6)

	//result, err := glamour.Render(theString, "./cmd/styles/dark.json")
	result, err := glamour.Render(theString, "dracula")

	if err != nil {
		panic(err)
	}
	fmt.Println(string(result))

}

// uploads a file to gemini
func addAFile(ctx context.Context, client *genai.Client) genai.Part {
	// during the chat, we can continuously update the below file by providing
	// a different diff. For example to get a diff for a golang repository,
	// we can issue the following command:
	// git diff -U10 7c904..dcfc69 -- . ':!vendor' > gitdiff.txt
	// the hashes are examples from two consecutive git hashes found when
	// simply doing a "git log" statement. Put the oldest hash first so that added
	// lines get a + and removed lines get a -, or you get it backwards.
	// note that the "-- . `:! vendor` part is to ignore the vendor file, as we are
	// only interested in actual updates of changes.
	fileContents, err := os.Open("./gitdiff.txt")
	if err != nil {
		panic(err)
	}
	upFile, err := client.UploadFile(ctx, "", fileContents, &genai.UploadFileOptions{MIMEType: "text/plain"})
	if err != nil {
		panic(err)
	}

	return genai.FileData{URI: upFile.URI}
}
