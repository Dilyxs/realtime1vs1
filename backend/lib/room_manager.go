package lib

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"realtime1vs1/randomhelper"

	"github.com/gorilla/websocket"
)

type RoomCommandType int

const (
	AddPlayerToRoom = iota
	AddPlayerToWebsocket
	RemovePlayerToWebsocket //:TODO: This is not utilized, WebsocketDisconnectMessage needs this as a future param
	CanUserJoin
	AddNewToken
	ValidateToken
)

const DefaultTimeout = 10 * time.Second

type IsGameOwner bool

type AddPlayerCommand struct {
	CommandType    RoomCommandType
	OutChan        chan RoomCommandResult
	PlayerUsername string
}

func (command AddPlayerCommand) basic() (RoomCommandType, chan RoomCommandResult) {
	return command.CommandType, command.OutChan
}

type AddPlayerToWebsocketCommand struct {
	CommandType    RoomCommandType
	OutChan        chan RoomCommandResult
	PlayerUsername string
	Conn           *websocket.Conn
}

func (command AddPlayerToWebsocketCommand) basic() (RoomCommandType, chan RoomCommandResult) {
	return command.CommandType, command.OutChan
}

type RoomCommand interface {
	basic() (RoomCommandType, chan RoomCommandResult)
}

type CheckIfUserAllowedToJoin struct {
	CommandType    RoomCommandType
	OutChan        chan RoomCommandResult
	PlayerUsername string
}

func (command CheckIfUserAllowedToJoin) basic() (RoomCommandType, chan RoomCommandResult) {
	return command.CommandType, command.OutChan
}

type WebsocketDisconnectMessage struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (msg WebsocketDisconnectMessage) ToJSON() []byte {
	jsonMsg, _ := json.Marshal(msg)
	return jsonMsg
}

type Room struct {
	IsClosed           bool
	ID                 int
	Players            map[string]bool
	GameMaster         string
	Chan               chan RoomCommand
	HubWebsocketChan   chan HubMessage
	SocketConns        map[string]chan HubMessage
	TokenManager       *TokenManager
	AllPreiousMessages []HubMessage
}
type Player struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RoomManager struct {
	RoomIDsoFar int
	Rooms       map[int]*Room
	Mu          sync.RWMutex
}

func NewManager() RoomManager {
	return RoomManager{
		RoomIDsoFar: 1,
		Rooms:       make(map[int]*Room),
		Mu:          sync.RWMutex{},
	}
}

func dropMessageOff(cmd HubMessage, chans map[string]chan HubMessage) {
	for _, socketChan := range chans {
		select {
		case socketChan <- cmd:
		default:
		}
	}
}

func (r *Room) Run() {
	for {
		select {
		case websocketMsg := <-r.HubWebsocketChan:
			//:NOTE: This is the internal room Request
			r.AllPreiousMessages = append(r.AllPreiousMessages, websocketMsg)
			switch cmd := websocketMsg.(type) {
			case WebsocketDisconnectMessage:
				delete(r.SocketConns, cmd.Username)
			default:
				dropMessageOff(cmd, r.SocketConns)
			}

		//:NOTE: this is the request coming from outside, http request most likely
		case command := <-r.Chan:
			switch cmd := command.(type) {
			case CheckIfUserAllowedToJoin:
				if r.IsClosed {
					cmd.OutChan <- RoomCommandResult{OK: false, Err: RoomError{ErrorCode: GameFull, Description: "game is full"}}
					continue
				}
				if r.Players[cmd.PlayerUsername] {
					if r.GameMaster == cmd.PlayerUsername {
						cmd.OutChan <- RoomCommandResult{OK: true, Err: nil, Extra: true}
						continue
					}
					cmd.OutChan <- RoomCommandResult{OK: true, Err: nil}
				} else {
					cmd.OutChan <- RoomCommandResult{OK: false, Err: RoomError{ErrorCode: UserNotAllowedToJoinGame, Description: "not allowed to join"}}
				}
			case AddPlayerCommand:
				r.Players[cmd.PlayerUsername] = true
				cmd.OutChan <- RoomCommandResult{OK: true, Err: nil}
			case AddPlayerToWebsocketCommand:
				localChan := make(chan HubMessage, 1000)
				go ReadFromWebsocket(cmd.Conn, r.HubWebsocketChan, cmd.PlayerUsername)
				go WriteToWebsocket(cmd.Conn, localChan)
				r.SocketConns[cmd.PlayerUsername] = localChan
				go WritePreviousMessagesToWebsocket(localChan, r.AllPreiousMessages)
				// simply drop off a new message in which we add a Ping that a NewUserHasCome
				select {
				case r.HubWebsocketChan <- UserWantsToJoin{ID: randomhelper.GetMessageID(), GamePhase: PreGame, Username: cmd.PlayerUsername}:
				default:
				}
				select {
				case cmd.OutChan <- RoomCommandResult{OK: true, Err: nil}:
				default:
				}
			}
		}
	}
}

func (r *RoomManager) CreateNewRoom(GameMaster string, tookenDis *TokenDistributer) int {
	tookenChannel := make(chan TokenMessage, 10)
	r.Mu.Lock()
	room := &Room{
		IsClosed:         false,
		ID:               r.RoomIDsoFar,
		Players:          make(map[string]bool),
		GameMaster:       GameMaster,
		Chan:             make(chan RoomCommand, 100),
		HubWebsocketChan: make(chan HubMessage, 100),
		SocketConns:      make(map[string]chan HubMessage),
		TokenManager: &TokenManager{
			Tokens:  make(map[string]PlayerUsernameRoom),
			HubChan: tookenChannel,
		},
		AllPreiousMessages: []HubMessage{},
	}
	r.Rooms[r.RoomIDsoFar] = room
	val := r.RoomIDsoFar
	r.RoomIDsoFar += 1
	room.Players[GameMaster] = true
	r.Mu.Unlock()
	go room.Run()
	go room.TokenManager.Run()
	tookenDis.Chans[room.ID] = tookenChannel
	return val
}

func (r *RoomManager) CheckIfRoomValid(roomID int) bool {
	r.Mu.RLock()
	defer r.Mu.RUnlock()
	return roomID < r.RoomIDsoFar
}

type RoomErrorCode int

const (
	RoomDoesNotExist = iota + 1
	GameFull
	UserNotAllowedToJoinGame
)

type RoomError struct {
	ErrorCode   RoomErrorCode `json:"error_code"`
	Description string        `json:"description"`
}

func (r RoomError) Error() string {
	return fmt.Sprintf("ErrorCode:%d Description: %s", r.ErrorCode, r.Description)
}

func (r *RoomManager) GetRoomChan(roomID int) (chan RoomCommand, error) {
	r.Mu.RLock()
	defer r.Mu.RUnlock()
	if val, ok := r.Rooms[roomID]; !ok {
		return nil, RoomError{ErrorCode: RoomDoesNotExist, Description: "RoomDoesNotExist"}
	} else {
		return val.Chan, nil
	}
}
