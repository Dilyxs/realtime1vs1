package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"realtime1vs1/lib"
)

type NewQuestionHandlerStruct struct {
	RoomID        int               `json:"room_id"`
	QuestionTopic lib.NicheProblems `json:"question_topic"`
}

// :TODO: currently we don't check the validity of the user, later add the token we created to verify this
func NewQuestionHandler(w http.ResponseWriter, r *http.Request, QDistrub *lib.QuestionDistributor, RoomGereur *lib.RoomManager) {
	if r.Method != http.MethodPost {
		//:TODO: convert these into json struct sendings
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req NewQuestionHandlerStruct
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsonMsg, _ := json.Marshal(ErrorMessageJSON{
			ErrorCode:        WrongFormat,
			ErrorMessageJSON: "wrong format",
		})
		w.Write(jsonMsg)

		return
	}
	if err := QDistrub.AddRoom(req.RoomID, RoomGereur); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return

	}
	roomChan := QDistrub.GetRoom(req.RoomID)

	myChan := make(chan lib.QuestionResult, 1)
	roomChan <- lib.CreateNewQuestionCommand{
		RoomID: req.RoomID,
		Chan:   myChan,
		Topic:  req.QuestionTopic,
	}
	select {
	case <-time.After(lib.DefaultTimeout):
		w.WriteHeader(http.StatusGatewayTimeout)
		return
	case res := <-myChan:
		switch cmd := res.(type) {
		case lib.RoomCreationResult:
			if cmd.Err != nil {
				w.WriteHeader(500)
				fmt.Fprintf(w, "ran into err:%v", cmd.Err)
			}
			jsonMsg, _ := json.Marshal(&cmd.Info)
			w.Write(jsonMsg)
		}
	}
}

func AnswerQuestionHandler(w http.ResponseWriter, r *http.Request, QDistrub *lib.QuestionDistributor) {
	if r.Method != http.MethodPost {
		jsonMsg, _ := json.Marshal(ErrorMessageJSON{
			ErrorCode:        WrongFormat,
			ErrorMessageJSON: "wrong format",
		})
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(jsonMsg)
		return
	}
	var q lib.UserQuestionResultJSON
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		jsonMsg, _ := json.Marshal(ErrorMessageJSON{
			ErrorCode:        WrongFormat,
			ErrorMessageJSON: "wrong format",
		})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonMsg)
		return
	}
	res, err := lib.AnswerQuestionGeneral(q.RoomID, q.Username, q.QuestionID, q.ChosenOption, QDistrub)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		jsonMsg, _ := json.Marshal(&err)
		w.Write(jsonMsg)
		return
	}
	jsonMsg, _ := json.Marshal(&res)
	w.WriteHeader(http.StatusAccepted)
	w.Write(jsonMsg)
}
