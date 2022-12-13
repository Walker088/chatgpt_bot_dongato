package web

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

const (
	URL           = "https://chat.openai.com/backend-api/conversation"
	USER_AGENT    = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.5112.102 Safari/537.36 Edg/104.0.1293.63"
	COOKIE_HEADER = ""
)

type MessageResponse struct {
	ConversationId string `json:"conversation_id"`
	Error          string `json:"error"`
	Message        struct {
		ID      string `json:"id"`
		Content struct {
			Parts []string `json:"parts"`
		} `json:"content"`
	} `json:"message"`
}

type ChatGPTWeb struct {
	SessionToken string
}

func NewEngine(session string) (*ChatGPTWeb, error) {
	return &ChatGPTWeb{
		SessionToken: session,
	}, nil
}

func (c *ChatGPTWeb) AskQuestion(question string) ([]byte, error) {
	enginReply := ""
	client := &http.Client{}
	req, err := c.prepareRequest(question, "", "")
	if err != nil {
		return nil, fmt.Errorf("failed to encode message: %v", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to encode message: %v", err)
	}
	reader := bufio.NewReader(resp.Body)
	defer resp.Body.Close()
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Println(fmt.Errorf("failed to decode event: %v", err))
		}
		if line == "[DONE]" || line == "" {
			break
		}

		sections := strings.SplitN(line, ":", 2)
		_, value := sections[0], ""
		if len(sections) == 2 {
			value = strings.TrimPrefix(sections[1], " ")
		}

		var res MessageResponse
		err = json.Unmarshal([]byte(value), &res)
		if err != nil {
			log.Printf("Couldn't unmarshal message response: %v", err)
			log.Println(line)
			continue
		}
		if len(res.Message.Content.Parts) > 0 {
			fmt.Println(res.Message.Content.Parts[0])
			enginReply += res.Message.Content.Parts[0]
		}
	}

	return []byte(enginReply), nil
}

func (c *ChatGPTWeb) prepareRequest(question string, conversationId string, lastMessageID string) (*http.Request, error) {
	messages, err := json.Marshal([]string{question})
	if err != nil {
		return nil, fmt.Errorf("failed to encode message: %v", err)
	}

	if lastMessageID == "" {
		lastMessageID = uuid.NewString()
	}

	var conversationIdString string
	if conversationId != "" {
		conversationIdString = fmt.Sprintf(`, "conversation_id": "%s"`, conversationId)
	}

	body := fmt.Sprintf(`{
        "action": "next",
        "messages": [
            {
                "id": "%s",
                "role": "user",
                "content": {
                    "content_type": "text",
                    "parts": %s
                }
            }
        ],
        "model": "text-davinci-002-render",
		"parent_message_id": "%s"%s
    }`, uuid.NewString(), string(messages), lastMessageID, conversationIdString)

	req, err := http.NewRequest("POST", URL, strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("User-Agent", USER_AGENT)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.SessionToken))
	req.Header.Set("Cookie", COOKIE_HEADER)
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}
