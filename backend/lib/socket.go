package lib

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func AddPlayerToWebsocketConn(w http.ResponseWriter, r *http.Request, roomManager *RoomManager, roomID int, playerUsername string) {
	roomchan, err := roomManager.GetRoomChan(roomID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusUpgradeRequired)
		return
	}
	outChan := make(chan RoomCommandResult, 1)
	roomchan <- AddPlayerToWebsocketCommand{
		CommandType:    AddPlayerToWebsocket,
		OutChan:        outChan,
		PlayerUsername: playerUsername,
		Conn:           conn,
	}
}
