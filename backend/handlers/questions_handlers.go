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

func NewQuestionHandler(w http.ResponseWriter, r *http.Request, QDistrub *lib.QuestionDistributor) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req NewQuestionHandlerStruct
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	QDistrub.AddRoom(req.RoomID)
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
