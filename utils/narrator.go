package utils

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
)

type Narrator struct {
	model *genai.GenerativeModel
}

func NewNarrator(client *genai.Client, modelID string) *Narrator {
	return &Narrator{
		model: client.GenerativeModel(modelID),
	}
}

func (n *Narrator) GenerateNarration(label string) (string, error) {
	ctx := context.Background()

	domainContext, err := os.ReadFile("config/domain_context/batik_philosophy.txt")

	var defaultNarrator string
	if err == nil {
		defaultNarrator = strings.TrimSpace(string(domainContext))
	} else {
		defaultNarrator = os.Getenv("DEFAULT_NARRATOR")
	}

	// inject
	fullInstruction := fmt.Sprintf(defaultNarrator, label)

	resp, err := n.model.GenerateContent(ctx, genai.Text(fullInstruction))
	if err != nil {
		return "", err
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return "", fmt.Errorf("gemini returned empty response")
	}

	var result string
	for _, part := range resp.Candidates[0].Content.Parts {
		result += fmt.Sprintf("%v", part)
	}

	return result, nil
}
