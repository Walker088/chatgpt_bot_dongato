package web

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
)

const (
	CONVERSATION_URL = "https://chat.openai.com/backend-api/conversation"
	USER_AGENT       = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.5112.102 Safari/537.36 Edg/104.0.1293.63"
)

type ChatGptMessageResponse struct {
	ConversationId string `json:"conversation_id"`
	Error          string `json:"error"`
	Message        struct {
		ID      string `json:"id"`
		Content struct {
			Parts []string `json:"parts"`
		} `json:"content"`
	} `json:"message"`
}

type ChatGptConversation struct {
	ID            string
	LastMessageID string
}

type ChatGPTWeb struct {
	sessionToken   string
	cfClearance    string
	jwtToken       JwtToken
	conversactions map[int64]ChatGptConversation

	mutex sync.Mutex
}

func NewEngine(session string, cfClearance string) (*ChatGPTWeb, error) {
	return &ChatGPTWeb{
		sessionToken:   session,
		cfClearance:    cfClearance,
		jwtToken:       NewEmptyJwtToken(session, cfClearance),
		conversactions: make(map[int64]ChatGptConversation),
	}, nil
}

func (c *ChatGPTWeb) AskQuestion(chatId int64, question string) ([]byte, error) {
	enginReply := ""
	client := &http.Client{}
	cvst := c.conversactions[chatId]

	req, err := c.prepareRequest(question, cvst.ID, cvst.LastMessageID)
	if err != nil {
		return nil, err
	}
	// Ensure that only one query at a time to chatgpt
	c.mutex.Lock()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
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

		var res ChatGptMessageResponse
		err = json.Unmarshal([]byte(value), &res)
		if err != nil {
			//log.Printf("Couldn't unmarshal message response: %v", err)
			continue
		}
		if len(res.Message.Content.Parts) > 0 {
			cvst.ID = res.ConversationId
			cvst.LastMessageID = res.Message.ID
			c.conversactions[chatId] = cvst
			enginReply = res.Message.Content.Parts[0]
		}
	}

	defer c.mutex.Unlock()
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

	req, err := http.NewRequest("POST", CONVERSATION_URL, strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	accessToken, err := c.jwtToken.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh chatgpt token: %v", err)
	}
	req.Header.Set("User-Agent", USER_AGENT)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Cookie", fmt.Sprintf(`cf_clearance=%s; __Secure-next-auth.session-token=%s`, c.cfClearance, c.sessionToken))
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}
