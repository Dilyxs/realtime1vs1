package handlers

import (
	"encoding/json"
	"net/http"

	"realtime1vs1/lib"
)

type RooomID struct {
	ID int `json:"id"`
}

func NewRoomHandler(w http.ResponseWriter, r *http.Request, roomManager *lib.RoomManger) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadGateway)
	}
	roomID := roomManager.CreateNewRoom()
	room := RooomID{
		ID: roomID,
	}
	roomJSON, _ := json.Marshal(&room)
	w.Write(roomJSON)
}
