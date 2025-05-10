package themodel

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	// NOTE: google announced end of live for the
	// used generative-ai-go library. The successor
	// of this api is implemented in the genaiterface
	// package
	"github.com/MelleKoning/aifun/internal/fileio"
	"github.com/MelleKoning/aifun/internal/terminal"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

const (
	modelName = "gemini-2.0-flash"
)

type theModel struct {
	client  *genai.Client
	model   *genai.GenerativeModel
	request *Request
}

type Request struct {
	textPart genai.Part
	filePart genai.Part
}

type Action interface {
	ReviewFile()
	ChatMessage(text string)
	CloseClient() error
}

func NewModel(ctx context.Context, systemInstruction string) (Action, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")

	genaiclient, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	return &theModel{
		client:  genaiclient,
		model:   setupModel(genaiclient, systemInstruction),
		request: &Request{},
	}, err
}

// Cleanup afer ourselves
func (m *theModel) CloseClient() error {
	return m.client.Close()
}
func setupModel(client *genai.Client, systemInstruction string) *genai.GenerativeModel {
	model := client.GenerativeModel(modelName)
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(systemInstruction)},
	}
	return model
}

// uploads a file to gemini
func addAFile(ctx context.Context, client *genai.Client) (genai.Part, string) {
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

	return genai.FileData{URI: upFile.URI}, upFile.URI
}

func (m *theModel) ReviewFile() {
	part, fileUri := addAFile(context.Background(), m.client)
	m.request.filePart = part
	commandText := `* Do not include the provided diff output in the response.

	The file {fileUri} contains a git diff output.

	AI OUTPUT:`
	commandText = strings.Replace(commandText, "{fileUri}", fileUri, 1)
	m.request.textPart = genai.Text(commandText)

	m.generateAndPrintResponse(context.Background())

}

func (m *theModel) ChatMessage(text string) {
	m.request.textPart = genai.Text(text)
	m.generateAndPrintResponse(context.Background())
}
func (m *theModel) generateAndPrintResponse(ctx context.Context) {
	var iter *genai.GenerateContentResponseIterator

	if m.request.filePart != nil {
		iter = m.model.GenerateContentStream(ctx, m.request.textPart, m.request.filePart)
	} else {
		iter = m.model.GenerateContentStream(ctx, m.request.textPart)

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

	fullString := buildString(allparts)

	terminal.PrintGlamourString(fullString)

	fileio.WriteMarkdown(fullString, "codereview.md")
}

func printResponse(resp *genai.GenerateContentResponse) {
	result := resp.Candidates[0].Content.Parts[0]

	fmt.Print(result)

}

func buildString(resp []genai.Part) string {
	var build strings.Builder
	for _, p := range resp {

		build.WriteString(fmt.Sprintf("%v", p))
	}

	return build.String()
}
