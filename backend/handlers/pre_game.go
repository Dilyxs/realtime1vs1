package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"realtime1vs1/lib"

	"github.com/gorilla/mux"
)

type IsPlayerAlloedToJoinAndGameMasterJSON struct {
	IsAllowedToJoin bool `json:"isAllowedToJoin"`
	IsGameMaster    bool `json:"isGameMaster"`
}

func PreGameHandler(w http.ResponseWriter, r *http.Request, roomManager *lib.RoomManager) {
	variables := mux.Vars(r)
	roomID, ok := variables["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		jsonerr, _ := json.Marshal(ErrorMessageJSON{ErrorMessageJSON: fmt.Sprintf("invalid game id: %s", roomID)})
		fmt.Fprint(w, jsonerr)
		return
	}
	roomIDint, err := strconv.ParseInt(roomID, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsonerr, _ := json.Marshal(ErrorMessageJSON{ErrorMessageJSON: fmt.Sprintf("invalid game id: %s", roomID)})
		fmt.Fprint(w, jsonerr)
	}
	plauserusername := r.URL.Query().Get("username")
	if plauserusername == "" {
		w.WriteHeader(http.StatusNotAcceptable)
		jsonErr, _ := json.Marshal(ErrorMessageJSON{ErrorMessageJSON: "missing username query parameter"})
		fmt.Fprint(w, string(jsonErr))
		return
	}

	roomChan, err := roomManager.GetRoomChan(int(roomIDint))
	if err != nil {
		jsonerr, _ := json.Marshal(ErrorMessageJSON{ErrorMessageJSON: fmt.Sprintf("invalid game id: %s", roomID)})
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, jsonerr)
		return
	}
	outchan := make(chan lib.RoomCommandResult, 1)
	command := lib.CheckIfUserAllowedToJoin{
		CommandType:    lib.AddPlayerToRoom,
		OutChan:        outchan,
		PlayerUsername: plauserusername,
	}
	select {
	case <-time.After(lib.DefaultTimeout):
		w.WriteHeader(http.StatusRequestTimeout)
		jsonErr, _ := json.Marshal(ErrorMessageJSON{ErrorMessageJSON: "request timed out"})
		fmt.Fprint(w, string(jsonErr))
		return
	case roomChan <- command:
	}
	select {
	case <-time.After(lib.DefaultTimeout):
		w.WriteHeader(http.StatusRequestTimeout)
		jsonErr, _ := json.Marshal(ErrorMessageJSON{ErrorMessageJSON: "request timed out"})
		fmt.Fprint(w, string(jsonErr))
		return
	case result := <-outchan:
		if !result.OK {
			w.WriteHeader(http.StatusForbidden)
			jsonErr, _ := json.Marshal(ErrorMessageJSON{ErrorMessageJSON: fmt.Sprintf("user %s is not allowed to join the game", plauserusername)})
			fmt.Fprint(w, string(jsonErr))
			return
		}
		cmd := result.Extra.(bool)
		if cmd {
			response := IsPlayerAlloedToJoinAndGameMasterJSON{
				IsAllowedToJoin: true,
				IsGameMaster:    true,
			}
			jsonResp, _ := json.Marshal(response)
			w.Write(jsonResp)
		} else {
			response := IsPlayerAlloedToJoinAndGameMasterJSON{
				IsAllowedToJoin: true,
				IsGameMaster:    false,
			}
			jsonResp, _ := json.Marshal(response)
			w.Write(jsonResp)
		}
	}
}
