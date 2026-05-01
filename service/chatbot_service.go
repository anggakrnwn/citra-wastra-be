package service

import (
	"citra-wastra-be/dto"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
)

type ChatbotService interface {
	Chat(ctx context.Context, message string) (dto.ChatResponse, error)
}

type chatbotService struct {
	model *genai.GenerativeModel
}

func NewChatbotService(client *genai.Client, modelID string) ChatbotService {
	model := client.GenerativeModel(modelID)

	instruction, err := os.ReadFile("config/domain_context/chatbot_instruction.txt")
	systemText := "Anda adalah asisten ahli Wastra Nusantara."
	if err == nil {
		systemText = string(instruction)
	}

	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{
			genai.Text(systemText),
		},
	}
	return &chatbotService{model: model}
}

func (s *chatbotService) Chat(ctx context.Context, message string) (dto.ChatResponse, error) {
	domainContext, _ := os.ReadFile("config/domain_context/batik_philosophy.txt")
	contextPrompt := ""
	if domainContext != nil {
		contextPrompt = "\nContext tambahan: " + string(domainContext)
	}

	fullPrompt := message + contextPrompt

	resp, err := s.model.GenerateContent(ctx, genai.Text(fullPrompt))
	if err != nil {
		return dto.ChatResponse{}, fmt.Errorf("failed to generate chat response: %v", err)
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return dto.ChatResponse{}, fmt.Errorf("chatbot returned empty response")
	}

	var reply strings.Builder
	for _, part := range resp.Candidates[0].Content.Parts {
		reply.WriteString(fmt.Sprintf("%v", part))
	}

	return dto.ChatResponse{
		Reply: reply.String(),
	}, nil
}
