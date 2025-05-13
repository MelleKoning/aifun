package genainterface

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	// genai is the successor of the previous
	// generative-ai-go model
	"google.golang.org/genai"

	"github.com/MelleKoning/aifun/internal/fileio"
	"github.com/MelleKoning/aifun/internal/terminal"
)

const (
	modelName = "gemini-2.0-flash"
)

type theModel struct {
	systemInstruction string
	client            *genai.Client
	chatHistory       []*genai.Content
}

type Action interface {
	ReviewFile() error
	ChatMessage(string)
	UpdateSystemInstruction(string)
	GetHistoryLength() int
}

// NewModel sets up the client for communication with Gemini. Ensure
// You need to have set your api key in env var GEMINI_API_KEY before
// calling the NewModel constructor
func NewModel(ctx context.Context, systemInstruction string) (Action, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")

	genaiclient, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, err
	}

	return &theModel{
		systemInstruction: systemInstruction,
		client:            genaiclient,
	}, err
}

func (m *theModel) GetHistoryLength() int {
	return len(m.chatHistory)
}
func (m *theModel) UpdateSystemInstruction(systemInstruction string) {
	m.systemInstruction = systemInstruction
}

func (m *theModel) ChatMessage(userPrompt string) {
	ctx := context.Background()

	// Add user prompt to chat history
	m.chatHistory = append(m.chatHistory, genai.NewContentFromText(userPrompt, genai.RoleUser))

	// Create chat with history
	chat, err := m.client.Chats.Create(ctx, modelName, nil, m.chatHistory)
	if err != nil {
		fmt.Println("Error creating chat:", err)
		return
	}

	// Send message to the model using streaming
	stream := chat.SendMessageStream(ctx, genai.Part{Text: userPrompt})

	var allModelParts []*genai.Part

	for chunk, err := range stream {
		if err != nil {
			fmt.Println("Error receiving stream:", err)
			break
		}
		part := chunk.Candidates[0].Content.Parts[0]
		printResponse(chunk)
		allModelParts = append(allModelParts, part)
	}

	// output model answer to console
	fullString := buildString(allModelParts)
	terminal.PrintGlamourString(fullString)

	// Add the combined response to chat history
	modelResponse := genai.NewContentFromText(fullString, genai.RoleModel)
	m.chatHistory = append(m.chatHistory, modelResponse)
}

func (m *theModel) ReviewFile() error {
	filePart, fileUri := m.addAFile(context.Background(), m.client)
	log.Printf("fileUri is %s", fileUri)

	// Start with chatHistory
	genaiContents := append([]*genai.Content{}, m.chatHistory...)

	// we first create a Part for file,
	// later we add an additional part
	// to this slice to add the Command below
	parts := []*genai.Part{
		filePart,
	}

	fileContent := genai.NewContentFromParts(parts, genai.RoleUser)

	// Include fileContent
	genaiContents = append(genaiContents, fileContent)

	config := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(m.systemInstruction, genai.RoleModel),
	}

	commandText := `* Do not include the provided diff output in the response.

		The file {fileUri} contains the git diff output to be reviewed.

		AI OUTPUT:`
	commandText = strings.Replace(commandText, "{fileUri}", fileUri, 1)

	// add command as additional part
	// to the last item in the genaiContents
	genaiCommandPart := &genai.Part{Text: commandText}
	lastContentPart := len(genaiContents) - 1
	genaiContents[lastContentPart].Parts = append(genaiContents[lastContentPart].Parts, genaiCommandPart)

	// add the command text to the file contents
	//genaiCommandText := genai.Text(commandText)
	//genaiContents = append(genaiContents, genaiCommandText...)
	stream := m.client.Models.GenerateContentStream(
		context.Background(),
		modelName,
		genaiContents,
		config,
	)

	var allModelParts []*genai.Part

	for chunk, err := range stream {
		if err != nil {
			fmt.Println(err)
			break
		}
		printResponse(chunk)

		part := chunk.Candidates[0].Content.Parts[0]
		allModelParts = append(allModelParts, part)

	}

	fullString := buildString(allModelParts)

	terminal.PrintGlamourString(fullString)

	// Combine all parts into a single part and add to chat history
	modelResponse := genai.NewContentFromText(fullString, genai.RoleModel)
	m.chatHistory = append(m.chatHistory, modelResponse)

	fileio.WriteMarkdown(fullString, "codereview.md")

	return nil
}

// uploads a file to gemini
func (m *theModel) addAFile(ctx context.Context, client *genai.Client) (*genai.Part, string) {
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
	upFile, err := client.Files.Upload(ctx, fileContents, &genai.UploadFileConfig{
		MIMEType: "text/plain",
	})
	if err != nil {
		panic(err)
	}

	return genai.NewPartFromURI(upFile.URI, upFile.MIMEType), upFile.URI
}

func printResponse(resp *genai.GenerateContentResponse) {
	result := resp.Candidates[0].Content.Parts[0]

	if result != nil {
		fmt.Print(".")
	} else {
		fmt.Print("-")
	}
}

func buildString(resp []*genai.Part) string {
	var build strings.Builder
	for _, p := range resp {

		build.WriteString(p.Text)
	}

	return build.String()
}
