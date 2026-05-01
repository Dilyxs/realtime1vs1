package lib

import (
	"fmt"
	"math/rand"
	"strconv"

	"realtime1vs1/evaluator"
)

type Performance struct {
	Username string
	Score    int
}

type AIScoreResponse struct {
	Score int `json:"score"`
}
type AIRequest struct {
	Answer string
	Chan   chan AIScoreResponse
}

//:TODO: currently we are praying that the LLM gives the correct expected result, need SOMETHING MUCH MUCH MORE ROBUST

func InitializeAIforQuestion(questionINFO ProblemNiche, answerChan <-chan AIRequest) {
	prompt := fmt.Sprintf("Your Job is to provide evaluate different responses to a particular question output back a string of a number '95' for example.Here is the question: %s", questionINFO.String())
	AI := evaluator.MakeNewClient("qwen3.5:4b", prompt, nil)
	for req := range answerChan {
		AI.AppendNewMessageFromUser(req.Answer)
		response := AI.GetModelResponse().Content
		var score int
		score, err := strconv.Atoi(response)
		if err != nil {
			score = int(rand.Int31n(101))
		}
		req.Chan <- AIScoreResponse{
			Score: score,
		}
	}
}

func EvaluateSingleQuestion(questionINFO ProblemNiche, username, answer string, reportChan chan<- Performance, LLMChan chan<- AIRequest) {
	localChan := make(chan AIScoreResponse, 1)
	LLMChan <- AIRequest{
		Answer: answer,
		Chan:   localChan,
	}
	response := <-localChan
	reportChan <- Performance{
		Username: username,
		Score:    response.Score,
	}
}

func EvaluatePerformance(questionINFO ProblemNiche, answers map[string]string) map[string]int {
	LLMChan := make(chan AIRequest, 50)
	localChan := make(chan Performance, 50)
	go InitializeAIforQuestion(questionINFO, LLMChan)

	for username, answer := range answers {
		go EvaluateSingleQuestion(questionINFO, username, answer, localChan, LLMChan)
	}
	result := make(map[string]int)
	for res := range localChan {
		result[res.Username] = res.Score
	}
	return result
}
