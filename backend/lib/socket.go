package lib

import (
	"encoding/json"
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
	outChan := make(chan RoomCommandResult, 3)
	roomchan <- AddPlayerToWebsocketCommand{
		CommandType:    AddPlayerToWebsocket,
		OutChan:        outChan,
		PlayerUsername: playerUsername,
		Conn:           conn,
	}
}

type GamePhase int

const (
	PreGame = iota
)

type HubMessage interface {
	ToJSON() []byte
}
type UserWantsToJoin struct {
	ID        string    `json:"id"`
	GamePhase GamePhase `json:"gamePhase"`
	Username  string    `json:"username"`
}

func (msg UserWantsToJoin) ToJSON() []byte {
	jsonMsg, _ := json.Marshal(msg)
	return jsonMsg
}

type UserIsReadyJSON struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	GamePhase GamePhase `json:"gamePhase"`
	IsReady   bool      `json:"isReady"`
}

func (msg UserIsReadyJSON) ToJSON() []byte {
	jsonMsg, _ := json.Marshal(msg)
	return jsonMsg
}

type UserWritingJSON struct {
	Type string          `json:"type"`
	Main json.RawMessage `json:"main,omitempty"`
}

func (msg UserWritingJSON) ToJSON() []byte {
	jsonMsg, _ := json.Marshal(msg)
	return jsonMsg
}

func ReadFromWebsocket(conn *websocket.Conn, HubChan chan HubMessage, playerUsernanme string) {
	intermidiatechan := make(chan UserWritingJSON, 20)
	go func() {
		for msg := range intermidiatechan {
			switch msg.Type {
			case "userIsReady":
				var customMsg UserIsReadyJSON
				if err := json.Unmarshal(msg.Main, &customMsg); err != nil {
					//:TODO: log this later on!
					continue
				}
				HubChan <- customMsg
			}
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

func WriteToWebsocket(conn *websocket.Conn, localChan chan HubMessage) {
	for msg := range localChan {
		if err := conn.WriteJSON(&msg); err != nil {
			return
		}
	}
}

func WritePreviousMessagesToWebsocket(websocketChan chan HubMessage, previousMessages []HubMessage) {
	for _, msg := range previousMessages {
		select {
		case websocketChan <- msg:
		default:
		}
	}
}
