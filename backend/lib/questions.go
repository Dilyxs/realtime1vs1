// Package lib cointains helper functions and structs for handlers to call!
package lib

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"

	"realtime1vs1/randomhelper"
)

type ProblemNicheCoreInfo struct {
	ProblemTopic        string `json:"problem_topic"`
	ProblemTimeRequired string `json:"problem_time_required"`
	ProblemDifficulty   string `json:"problem_difficulty"`
	ProblemDescription  string `json:"problem_description"`
}

type ProblemNiche struct {
	ProblemID           string   `json:"problem_id"`
	ProblemTopic        string   `json:"problem_topic"`
	ProblemTimeRequired string   `json:"problem_time_required"`
	ProblemDifficulty   string   `json:"problem_difficulty"`
	ProblemDescription  string   `json:"problem_description"`
	ProblemHints        []string `json:"problem_hints"`
	ProblemRubric       []Rubric `json:"problem_rubric"`
}

func (p ProblemNiche) TOJSON() []byte {
	res, _ := json.Marshal(p)
	return res
}

type Rubric struct {
	Criterion   string `json:"criterion"`
	Points      int    `json:"points"`
	Description string `json:"description"`
}

type ProblemGeneral struct {
	QuestionID int      `json:"questionID"`
	Question   string   `json:"question"`
	Options    []string `json:"options"`
	Answer     int      `json:"answer"`
	Topic      string   `json:"topic"`
	Difficulty string   `json:"difficulty"`
}

func (p ProblemGeneral) TOJSON() []byte {
	res, _ := json.Marshal(p)
	return res
}

const (
	Easy   = "easy"
	Medium = "medium"
	Hard   = "hard"
)

type QuestionDistributor struct {
	Chans              map[int]chan Question
	Mu                 sync.RWMutex
	GeneralQuestionAll []ProblemGeneral
	NicheQuestionAll   []ProblemNiche
}

func NewQuestionManager(pathGeneral, pathNiche string) *QuestionDistributor {
	contentNiche, err := ReadFileAndReturn[ProblemNiche](pathNiche)
	if err != nil {
		log.Fatalf("err reading niche problems: %v", err)
	}
	contentGeneral, err := ReadFileAndReturn[ProblemGeneral](pathNiche)
	if err != nil {
		log.Fatalf("err reading niche problems: %v", err)
	}
	return &QuestionDistributor{
		Chans:              make(map[int]chan Question),
		Mu:                 sync.RWMutex{},
		GeneralQuestionAll: contentGeneral,
		NicheQuestionAll:   contentNiche,
	}
}

func (q *QuestionDistributor) AddRoom(roomID int) {
	Chan := make(chan Question, 100)
	QuestionMan := QuestionManager{
		RoomID:             roomID,
		Chan:               Chan,
		Topic:              ReactProblem,         // BY default
		AllNicheProblems:   q.NicheQuestionAll,   // this is always read only
		AllGeneralProblems: q.GeneralQuestionAll, // this is always read only
	}
	go QuestionMan.Run()
	q.Mu.Lock()
	q.Chans[roomID] = Chan
	q.Mu.Unlock()
}

func (q *QuestionDistributor) GetRoom(roomID int) chan Question {
	q.Mu.RLock()
	defer q.Mu.RUnlock()
	return q.Chans[roomID]
}

type QuestionManager struct {
	RoomID             int
	Chan               chan Question
	Topic              NicheProblems
	ProblemAtHand      ProblemNiche
	AllNicheProblems   []ProblemNiche
	AllGeneralProblems []ProblemGeneral
}
type QuestionResult interface {
	hasID() string
}
type Question interface {
	hasChan() chan QuestionResult
}
type CreateNewQuestionCommand struct {
	RoomID int
	Chan   chan QuestionResult
	Topic  NicheProblems
}

func (command CreateNewQuestionCommand) hasChan() chan QuestionResult {
	return command.Chan
}

type NicheProblems int

type IsAFile interface {
	TOJSON() []byte
}

func ReadFileAndReturn[T IsAFile](filpath string) ([]T, error) {
	content, err := os.Open(filpath)
	if err != nil {
		return nil, fmt.Errorf("file does not exist")
	}
	defer func() {
		err := content.Close()
		if err != nil {
			return
		}
	}()
	var res []T
	if err := json.NewDecoder(content).Decode(&res); err != nil {
		return nil, fmt.Errorf("wrong json format")
	}
	return res, nil
}

const (
	ReactProblem = iota
	GoProblem
	CProblem
	RustProblem
)

func GetANicheProblem(NicheProblem NicheProblems, losProblems []ProblemNiche) ProblemNiche {
	basicIndex := NicheProblem
	basicIndex *= 20
	problemIndex := rand.Int31n(20) + int32(basicIndex)
	if int(problemIndex) >= len(losProblems) {
		log.Fatal("problem with Indexing for finding a valid problem: questions.go 165")
	}
	return losProblems[int(problemIndex)]
}

type RoomCreationResult struct {
	ID   string
	Info ProblemNicheCoreInfo
	Err  error
}

func (r RoomCreationResult) hasID() string {
	return r.ID
}

func (q *QuestionManager) Run() {
	for request := range q.Chan {
		switch cmd := request.(type) {
		case CreateNewQuestionCommand:
			q.Topic = NicheProblems(cmd.Topic)
			problem := GetANicheProblem(q.Topic, q.AllNicheProblems)
			q.ProblemAtHand = problem
			cmd.Chan <- RoomCreationResult{
				ID: randomhelper.GetMessageID(),
				Info: ProblemNicheCoreInfo{
					ProblemTopic:        problem.ProblemTopic,
					ProblemTimeRequired: problem.ProblemTimeRequired,
					ProblemDifficulty:   problem.ProblemDifficulty,
					ProblemDescription:  problem.ProblemDescription,
				},
				Err: nil,
			}
		}
	}
}
