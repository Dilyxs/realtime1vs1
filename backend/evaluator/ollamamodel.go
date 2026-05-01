package evaluator

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/ollama/ollama/api"
)

type AIplayer struct {
	Modelname     string
	Client        *api.Client
	InitialPrompt string
	ChatHistory   []api.Message
}

func MakeNewClient(ModelName, InitialPrompt string, ollamaURL *string) AIplayer {
	if ollamaURL == nil {
		newurl := os.Getenv("OLLAMA_URL")
		ollamaURL = &newurl
	}
	u, err := url.Parse(*ollamaURL)
	if err != nil {
		log.Fatalf("could not parse url: %v\n err: %v", ollamaURL, err)
	}
	client := api.NewClient(u, http.DefaultClient)
	return AIplayer{
		Modelname:     ModelName,
		Client:        client,
		InitialPrompt: InitialPrompt,
		ChatHistory:   make([]api.Message, 0),
	}
}

func (AI *AIplayer) HandleEachNewResponse(newrespopnse api.ChatResponse) error {
	AI.ChatHistory = append(AI.ChatHistory, newrespopnse.Message)
	return nil
}

func (AI *AIplayer) AppendNewMessageFromUser(newcontent string) {
	AI.ChatHistory = append(AI.ChatHistory, api.Message{Role: "user", Content: newcontent})
}

func (AI *AIplayer) GetModelResponse() api.Message {
	noStreaming := false
	NewRequest := api.ChatRequest{
		Model:    AI.Modelname,
		Messages: AI.ChatHistory,
		Stream:   &noStreaming,
	}

	err := AI.Client.Chat(context.Background(), &NewRequest, AI.HandleEachNewResponse)
	if err != nil {
		log.Fatalf("can no longer communicate:%v\n", err)
	}
	return AI.ChatHistory[len(AI.ChatHistory)-1]
}
