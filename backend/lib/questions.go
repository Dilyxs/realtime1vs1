// Package lib cointains helper functions and structs for handlers to call!
package lib

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

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

func (p ProblemNiche) ToJSON() []byte {
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

type ProblemGeneralCoreInfo struct {
	QuestionID int      `json:"questionID"`
	Question   string   `json:"question"`
	Options    []string `json:"options"`
	Topic      string   `json:"topic"`
	Difficulty string   `json:"difficulty"`
}

func (p ProblemGeneralCoreInfo) ToJSON() []byte {
	jsonMsg, _ := json.Marshal(&p)
	return jsonMsg
}

func (p ProblemGeneral) ToJSON() []byte {
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

func (q *QuestionDistributor) AddRoom(roomID int, RoomGereur *RoomManager) error {
	wsChan, err := RoomGereur.GetWebsocketChan(roomID)
	if err != nil {
		return err
	}
	Chan := make(chan Question, 100)
	QuestionMan := QuestionManager{
		RoomID:             roomID,
		Chan:               Chan,
		Topic:              ReactProblem,         // BY default
		AllNicheProblems:   q.NicheQuestionAll,   // this is always read only
		AllGeneralProblems: q.GeneralQuestionAll, // this is always read only
		WebsocketChan:      wsChan,
	}
	go QuestionMan.Run()
	q.Mu.Lock()
	q.Chans[roomID] = Chan
	q.Mu.Unlock()
	return nil
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
	WebsocketChan      chan HubMessage
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
	ToJSON() []byte
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

type UserQuestionResult struct {
	Username     string `json:"username"`
	ID           int    `json:"id"`
	ChosenOption int    `json:"option"`
	Chan         chan QuestionResult
}

func (q UserQuestionResult) hasChan() chan QuestionResult {
	return q.Chan
}

type GameHasStarted struct {
	ID         string `json:"ID"`
	HasStarted bool   `json:"has_started"`
}

func (msg GameHasStarted) ToJSON() []byte {
	res, _ := json.Marshal(msg)
	return res
}

func (q *QuestionManager) Run() {
	localChan := make(chan UserQuestionResult, 100)
	go q.AskQuestions(localChan)
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
			// to let everyone know that the game has started:
			q.WebsocketChan <- GameHasStarted{
				ID:         randomhelper.GetMessageID(),
				HasStarted: true,
			}
		case UserQuestionResult:
			localChan <- cmd
		}
	}
}

type QuestionGeneralAnswerResult struct {
	ID         string
	Registered bool
}

func (q QuestionGeneralAnswerResult) hasID() string {
	return q.ID
}

type QuestionGeneralWebsocketOutput struct {
	SuccessfulPlayers []PlayerAndOption `json:"successful_players"`
	FailedPlayers     []PlayerAndOption `json:"failed_players"`
}
type PlayerAndOption struct {
	Username string `json:"username"`
	Option   int    `json:"option"`
}

func (q *QuestionManager) AskQuestions(localChan <-chan UserQuestionResult) {
	contentMinute := q.ProblemAtHand.ProblemTimeRequired
	minString := strings.Split(contentMinute, " ")[0]
	timeTotal, err := strconv.Atoi(minString)
	if err != nil {
		log.Fatalf("could not convert time into minute: %v", err)
	}
	totalGeneralQuestions := (timeTotal / 5)
	for range totalGeneralQuestions {
		pickedQuestionID := rand.Int31n(100)
		pickedQuestion := q.AllGeneralProblems[pickedQuestionID]

		// to add some randomness sleep for an unknow time
		time.Sleep(time.Duration(rand.Int31n(200) * int32(time.Second)))
		q.WebsocketChan <- ProblemGeneralCoreInfo{
			QuestionID: pickedQuestion.QuestionID,
			Question:   pickedQuestion.Question,
			Options:    pickedQuestion.Options,
			Topic:      pickedQuestion.Topic,
			Difficulty: pickedQuestion.Difficulty,
		}
		result := make(map[bool][]PlayerAndOption)

		timerChan := time.NewTicker(5 * time.Second)
	InfiniteLoop:
		for {
			select {
			case <-timerChan.C:
				break InfiniteLoop
			case req := <-localChan:
				if req.ID != int(pickedQuestionID) {
					req.Chan <- QuestionGeneralAnswerResult{ID: randomhelper.GetMessageID(), Registered: false}
					continue
				}
				if pickedQuestion.Answer == req.ChosenOption {
					result[true] = append(result[true], PlayerAndOption{
						Username: req.Username,
						Option:   req.ChosenOption,
					})
				} else {
					result[false] = append(result[false], PlayerAndOption{
						Username: req.Username,
						Option:   req.ChosenOption,
					})
				}
				req.Chan <- QuestionGeneralAnswerResult{ID: randomhelper.GetMessageID(), Registered: true}
			}
		}
		// Now just send it to Websocket AS a final output!
		var res QuestionGeneralWebsocketOutput
		res.SuccessfulPlayers = result[true]
		res.FailedPlayers = result[false]
		q.WebsocketChan <- GeneralQuestionUserAnsweredResult{
			ID:        randomhelper.GetMessageID(),
			GamePhase: DuringGame,
			Result:    res,
		}
	}
}
