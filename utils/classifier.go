package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type aiClassifier struct{}

func NewClassifier() *aiClassifier { return &aiClassifier{} }

func (c *aiClassifier) ClassifyBatik(imageURL string) (string, float64, string, error) {
	baseURL := strings.TrimRight(os.Getenv("AI_MODEL_URL"), "/")
	apiKey := os.Getenv("AI_MODEL_TOKEN")

	if baseURL == "" || apiKey == "" {
		return "", 0, "", errors.New("Environment variables AI_MODEL_URL or AI_MODEL_TOKEN are not set")
	}

	callURL := fmt.Sprintf("%s/gradio_api/call/predict", baseURL)
	payload := map[string]interface{}{
		"data": []interface{}{
			map[string]interface{}{
				"path": imageURL,
				"meta": map[string]interface{}{"_type": "gradio.FileData"},
			},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", 0, "", fmt.Errorf("Failed to marshal payload: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", callURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", 0, "", fmt.Errorf("Failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, "", fmt.Errorf("AI initiation request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", 0, "", fmt.Errorf("AI call initiation failed (%d): %s", resp.StatusCode, string(body))
	}

	var callRes struct {
		EventID string `json:"event_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&callRes); err != nil {
		return "", 0, "", fmt.Errorf("Failed to decode event_id: %w", err)
	}

	if callRes.EventID == "" {
		return "", 0, "", errors.New("Empty event_id received from AI service")
	}

	resultURL := fmt.Sprintf("%s/gradio_api/call/predict/%s", baseURL, callRes.EventID)
	var finalBody []byte

	for i := 0; i < 5; i++ {
		reqResult, err := http.NewRequestWithContext(ctx, "GET", resultURL, nil)
		if err != nil {
			continue
		}

		reqResult.Header.Set("Authorization", "Bearer "+apiKey)

		respResult, err := client.Do(reqResult)
		if err != nil {
			continue
		}

		tempBytes, err := io.ReadAll(respResult.Body)
		respResult.Body.Close()

		if err == nil && strings.Contains(string(tempBytes), "data:") {
			finalBody = tempBytes
			break
		}

		time.Sleep(2 * time.Second)
	}

	lines := strings.Split(string(finalBody), "\n")
	var finalJSON string
	for _, line := range lines {
		if strings.HasPrefix(line, "data: ") {
			finalJSON = strings.TrimPrefix(line, "data: ")
			break
		}
	}

	if finalJSON == "" {
		return "", 0, "", errors.New("No result data found in AI response stream")
	}

	var outputData []interface{}
	if err := json.Unmarshal([]byte(finalJSON), &outputData); err != nil {
		return "", 0, "", fmt.Errorf("Failed to unmarshal result: %v", err)
	}

	if len(outputData) < 2 {
		return "", 0, "", errors.New("AI response data is incomplete")
	}

	labelData, ok := outputData[0].(map[string]interface{})
	if !ok {
		return "", 0, "", errors.New("Invalid format for label data")
	}

	confidences, ok := labelData["confidences"].([]interface{})
	if !ok || len(confidences) == 0 {
		return "", 0, "", errors.New("No motif detected")
	}

	top, ok := confidences[0].(map[string]interface{})
	if !ok {
		return "", 0, "", errors.New("Invalid confidence format")
	}

	label := fmt.Sprintf("%v", top["label"])

	confidence, ok := top["confidence"].(float64)
	if !ok {
		return "", 0, "", errors.New("Confidence value is not a number")
	}

	philosophyRaw, ok := outputData[1].(string)
	if !ok {
		return "", 0, "", errors.New("Invalid philosophy format")
	}

	return label, confidence, philosophyRaw, nil
}
