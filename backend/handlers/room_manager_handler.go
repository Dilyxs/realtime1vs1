package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"realtime1vs1/db"
	"realtime1vs1/lib"
)

type RooomID struct {
	ID int `json:"id"`
}

type PlayerUsername struct {
	Username string `json:"username"`
}

func NewRoomHandler(w http.ResponseWriter, r *http.Request, roomManager *lib.RoomManager, tokenDis *lib.TokenDistributer) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	var GameMaster PlayerUsername
	err := json.NewDecoder(r.Body).Decode(&GameMaster)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if GameMaster.Username == "" {
		w.WriteHeader(http.StatusBadRequest)
		jsonMsg, _ := json.Marshal(&ErrorMessageJSON{ErrorMessageJSON: "invalid request body, missing username"})
		w.Write(jsonMsg)
		return
	}
	roomID := roomManager.CreateNewRoom(GameMaster.Username, tokenDis)
	room := RooomID{
		ID: roomID,
	}
	roomJSON, _ := json.Marshal(&room)
	w.Write(roomJSON)
}

func AddNewPlayerHandler(w http.ResponseWriter, r *http.Request, roomManager *lib.RoomManager) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadGateway)
		return
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
	OutChan := make(chan lib.RoomCommandResult, 3)
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

type PlayerAndRoom struct {
	Username string `json:"username"`
	Password string `json:"password"`
	RoomID   int    `json:"roomid"`
}

func TokenReturnHandler(w http.ResponseWriter, r *http.Request, poolManager *db.PoolManager, roomM *lib.RoomManager, tokenDis *lib.TokenDistributer) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	var info PlayerAndRoom
	err := json.NewDecoder(r.Body).Decode(&info)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsonMsg, _ := json.Marshal(&ErrorMessageJSON{ErrorMessageJSON: "invalid request body"})
		w.Write(jsonMsg)
		return
	}

	res, _, _ := verifyUserPassword(lib.Player{Username: info.Username, Password: info.Password}, poolManager)
	if !res.Valid {
		w.WriteHeader(http.StatusNotAcceptable)
		jsonMsg, _ := json.Marshal(&ErrorMessageJSON{
			ErrorMessageJSON: "incorrect, pls verify /login",
		})
		w.Write(jsonMsg)
		return
	}

	if !roomM.CheckIfRoomValid(info.RoomID) {
		jsonMsg, _ := json.Marshal(&ErrorMessageJSON{ErrorCode: RoomDoesNotExist, ErrorMessageJSON: "invalid room id"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonMsg)
		return
	}

	localChan := make(chan string, 2)
	tokenReq := lib.AddNewUserTokenCommand{
		TokenType:  lib.AddNewToken,
		PlayerInfo: lib.PlayerUsernameRoom{Username: info.Username, RoomID: info.RoomID},
		OutChan:    localChan,
	}
	tokenDis.Chans[info.RoomID] <- tokenReq

	select {
	case <-time.After(lib.DefaultTimeout):
		w.WriteHeader(http.StatusRequestTimeout)
		jsonMsg, _ := json.Marshal(&ErrorMessageJSON{ErrorMessageJSON: "request timed out"})
		w.Write(jsonMsg)
	case result := <-localChan:
		w.WriteHeader(http.StatusAccepted)
		jsonMsg, _ := json.Marshal(TokenJSON{Token: result})
		w.Write(jsonMsg)
	}
}

type TokenJSON struct {
	Token string `json:"token"`
}

func AddPlayerToWebsocketHandler(w http.ResponseWriter, r *http.Request, roomManager *lib.RoomManager, tkmanager *lib.TokenDistributer) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		jsonMsg, _ := json.Marshal(&ErrorMessageJSON{ErrorMessageJSON: "method not allowed"})
		w.Write(jsonMsg)
		return
	}

	token := r.URL.Query().Get("token")
	if token == "" {
		w.WriteHeader(http.StatusBadRequest)
		jsonMsg, _ := json.Marshal(&ErrorMessageJSON{ErrorMessageJSON: "invalid request body, missing username or roomid"})
		w.Write(jsonMsg)
		return
	}
	roomId := r.URL.Query().Get("roomid")
	if roomId == "" {
		w.WriteHeader(http.StatusBadRequest)
		jsonMsg, _ := json.Marshal(&ErrorMessageJSON{ErrorMessageJSON: "invalid request body, missing username or roomid"})
		w.Write(jsonMsg)
		return
	}
	roomID, err := strconv.Atoi(roomId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsonMsg, _ := json.Marshal(&ErrorMessageJSON{ErrorMessageJSON: "invalid request body, missing username or roomid"})
		w.Write(jsonMsg)
		return
	}
	localChan := make(chan struct {
		PlayerInfo lib.PlayerUsernameRoom
		Valid      bool
	}, 2)
	tokenReq := lib.ValidateTokenCommand{
		TokenType: lib.ValidateToken, TokenContent: token,
		OutChan: localChan,
	}
	select {
	case <-time.After(lib.DefaultTimeout):
		w.WriteHeader(http.StatusRequestTimeout)
		jsonMsg, _ := json.Marshal(&ErrorMessageJSON{ErrorMessageJSON: "request timed out"})
		w.Write(jsonMsg)
		return
	case tkmanager.Chans[roomID] <- tokenReq:
	}
	var PlayerInfo *lib.PlayerUsernameRoom
	select {
	case <-time.After(lib.DefaultTimeout):
		w.WriteHeader(http.StatusRequestTimeout)
		jsonMsg, _ := json.Marshal(&ErrorMessageJSON{ErrorMessageJSON: "request timed out"})
		w.Write(jsonMsg)
		return
	case valid := <-localChan:
		if !valid.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			jsonMsg, _ := json.Marshal(&ErrorMessageJSON{ErrorMessageJSON: "invalid token"})
			w.Write(jsonMsg)
			return
		} else {
			PlayerInfo = &valid.PlayerInfo
		}
	}
	lib.AddPlayerToWebsocketConn(w, r, roomManager, PlayerInfo.RoomID, PlayerInfo.Username)
}
