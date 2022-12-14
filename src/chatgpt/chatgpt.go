package chatgpt

type ChatGPT interface {
	AskQuestion(chatId int64, question string) ([]byte, error)
}
