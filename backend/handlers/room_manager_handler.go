package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"realtime1vs1/lib"
)

type RooomID struct {
	ID int `json:"id"`
}

type PlayerUsername struct {
	Username string `json:"username"`
}

func NewRoomHandler(w http.ResponseWriter, r *http.Request, roomManager *lib.RoomManger) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadGateway)
	}
	var GameMaster PlayerUsername
	err := json.NewDecoder(r.Body).Decode(&GameMaster)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	roomID := roomManager.CreateNewRoom(GameMaster.Username)
	room := RooomID{
		ID: roomID,
	}
	roomJSON, _ := json.Marshal(&room)
	w.Write(roomJSON)
}

func AddNewPlayerHandler(w http.ResponseWriter, r *http.Request, roomManager *lib.RoomManger) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadGateway)
	}
	var NewPlayer PlayerUsername
	err := json.NewDecoder(r.Body).Decode(&NewPlayer)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	roomid := r.URL.Query().Get("roomID")
	if roomid == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	roomID, err := strconv.Atoi(roomid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	roomchan, err := roomManager.GetRoomChan(roomID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsonMg, _ := json.Marshal(&err)
		w.Write(jsonMg)
		w.Write([]byte(err.Error()))
		return
	}
	OutChan := make(chan lib.RoomCommandResult, 1)
	RoomCommand := lib.AddPlayerCommand{
		CommandType:    lib.AddPlayerToRoom,
		PlayerUsername: NewPlayer.Username,
		OutChan:        OutChan,
	}
	select {
	case roomchan <- RoomCommand:
	default:
		w.WriteHeader(http.StatusInternalServerError)
		jsonMsg, _ := json.Marshal(&ErrorMessageJSON{ErrorMessageJSON: "server overloaded, try again later!"})
		w.Write(jsonMsg)
		return
	}
	select {
	case <-time.After(4 * time.Second):
		w.WriteHeader(http.StatusInternalServerError)
		jsonMsg, _ := json.Marshal(&ErrorMessageJSON{ErrorMessageJSON: "server overloaded, try again later!"})
		w.Write(jsonMsg)
		return
	//:NOTE: We could never get an Error with Adding a new player, but we can get an error if the room doesn't exist, so we will just return the error if it exists, otherwise we will return the result of adding a new player to the room.
	case <-OutChan:
		w.WriteHeader(http.StatusOK)
	}
}
