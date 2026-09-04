package main

import (
	"context"
	"fmt"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
)

func main() {
	ctx := context.Background()
	client := anthropic.NewClient()

	myFile, err := os.Open("./file-upload/file.txt")
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}

	fileUploadResult, err := client.Files.Upload(ctx, anthropic.FileUploadParams{
		File: anthropic.File(myFile, "file.txt", "text/plain"),
	})
	if err != nil {
		fmt.Printf("Error uploading file: %v\n", err)
		return
	}
	content := "Write me a summary of my file.txt file in the style of a Shakespearean sonnet.\n\n"
	println("[user]: " + content)

	message, err := client.Messages.New(ctx, anthropic.MessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(
				anthropic.NewTextBlock(content),
				anthropic.NewDocumentBlock(anthropic.FileDocumentSourceParam{
					FileID: fileUploadResult.ID,
				}),
			),
		},
		Model: anthropic.ModelClaudeSonnet5,
	})
	if err != nil {
		fmt.Printf("Error creating message: %v\n", err)
		return
	}

	println("[assistant]: " + message.Content[0].Text + message.StopSequence)
}
