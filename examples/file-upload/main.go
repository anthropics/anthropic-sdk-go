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

	myFile, err := os.Open("examples/file-upload/file.txt")
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}

	fileUploadResult, err := client.Beta.Files.Upload(ctx, anthropic.BetaFileUploadParams{
		File:  anthropic.File(myFile, "file.txt", "text/plain"),
		Betas: []anthropic.AnthropicBeta{anthropic.AnthropicBetaFilesAPI2025_04_14},
	})
	if err != nil {
		fmt.Printf("Error uploading file: %v\n", err)
		return
	}
	content := "Write me a summary of my file.txt file in the style of a Shakespearean sonnet.\n\n"
	println("[user]: " + content)

	message, err := client.Beta.Messages.New(ctx, anthropic.BetaMessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.BetaMessageParam{
			anthropic.NewBetaUserMessage(
				anthropic.NewBetaTextBlock(content),
				anthropic.NewBetaDocumentBlock(anthropic.BetaFileDocumentSourceParam{
					FileID: fileUploadResult.ID,
				}),
			),
		},
		Model: anthropic.ModelClaudeSonnet4_20250514,
		Betas: []anthropic.AnthropicBeta{anthropic.AnthropicBetaFilesAPI2025_04_14},
	})
	if err != nil {
		fmt.Printf("Error creating message: %v\n", err)
		return
	}

	println("[assistant]: " + message.Content[0].Text + message.StopSequence)
}
