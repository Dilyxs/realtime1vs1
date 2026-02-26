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

type HubMessage interface {
	error_code() int
}
type UserIsReadyJSON struct {
	username string
	isReady  bool
}

func (msg UserIsReadyJSON) error_code() int {
	return 0
}

func ReadFromWebsocket(conn *websocket.Conn, HubChan chan HubMessage, playerUsernanme string) {
	intermidiatechan := make(chan UserWritingJSON, 20)
	go func() {
		for msg := range intermidiatechan {
			switch customMsg := msg.Main.(type) {
			case UserIsReadyJSON:
				HubChan <- UserIsReadyJSON{username: customMsg.username, isReady: customMsg.isReady}
			}

			HubChan <- msg
		}
	}()
	for {
		var msg UserWritingJSON
		if err := conn.ReadJSON(&msg); err != nil {
			HubChan <- WebsocketDisconnectMessage{Username: playerUsernanme}
			return
		}
		intermidiatechan <- msg
	}
}

func WriteToWebsocket(conn *websocket.Conn, localChan chan UserWritingJSON) {
	for msg := range localChan {
		if err := conn.WriteJSON(&msg); err != nil {
			return
		}
	}
}
