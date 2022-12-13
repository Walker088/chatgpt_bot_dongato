package chatgpt

type ChatGPT interface {
	AskQuestion(question string) ([]byte, error)
}
