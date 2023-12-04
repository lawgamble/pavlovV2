package gpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func GetChatGPTReply(messages []Message) string {
	request := Request{
		Model:       "gpt-3.5-turbo",
		Messages:    messages,
		Temperature: 0.7,
	}

	// Convert data to JSON
	payload, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error:", err)
		return fmt.Sprintf("I've got no response for ya, sorry:\n%v", err)
	}

	// Create HTTP POST request with bearer token in Authorization header
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewReader(payload))
	if err != nil {
		fmt.Println("Error:", err)
		return fmt.Sprintf("I've got no response for ya, sorry:\n%v", err)
	}
	req.Header.Set("Authorization", "Bearer "+os.Getenv("GPT_API_KEY"))
	req.Header.Set("Content-Type", "application/json")

	// Create HTTP client
	client := &http.Client{}
	// Make HTTP POST request to ChatGPT API
	response, err := client.Do(req)
	if response.StatusCode != http.StatusOK {
		fmt.Println("Error:", err)
		return fmt.Sprintf("I've got no response for ya, sorry:\n%v - %v", response.StatusCode, response.Status)
	}

	defer response.Body.Close()

	// Decode response from JSON
	var chatGPTResponse ChatGPTResponse
	err = json.NewDecoder(response.Body).Decode(&chatGPTResponse)
	if err != nil {
		fmt.Println("Error:", err)
		return fmt.Sprintf("I've got no response for ya, sorry:\n%v", err)
	}

	// Check if Choices field is not empty
	if len(chatGPTResponse.Choices) > 0 {
		// Display chatGPTResponse content in output widget
		return chatGPTResponse.Choices[0].Message.Content
	} else {
		return fmt.Sprintf("I've got no response for ya, sorry:\n%v", err)
	}
}
