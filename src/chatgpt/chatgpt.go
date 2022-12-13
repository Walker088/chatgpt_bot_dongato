package chatgpt

type ChatGPT struct {
	SessionToken string
}

func NewChatgptClient(session string) *ChatGPT {
	return &ChatGPT{
		SessionToken: session,
	}
}

func (c *ChatGPT) AskQuestion() {

}
