package lib

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type RoomCommandType int

const (
	AddPlayerToRoom = iota
	AddPlayerToWebsocket
	RemovePlayerToWebsocket
)

type RoomCommandResult struct {
	OK  bool
	Err error
}
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
type HubMessage interface {
	error_code() int
}
type UserWritingJSON struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}

func (mesg UserWritingJSON) error_code() int {
	return 0
}

type WebsocketDisconnectMessage struct {
	Username string
}

func (mesg WebsocketDisconnectMessage) error_code() int {
	return 1
}

type Room struct {
	IsClosed         bool
	ID               int
	Players          map[string]bool
	GameMaster       string
	Chan             chan RoomCommand
	HubWebsocketChan chan HubMessage
	SocketConns      map[string]chan UserWritingJSON
}

func ReadFromWebsocket(conn *websocket.Conn, HubChan chan HubMessage, playerUsernanme string) {
	for {
		var msg UserWritingJSON
		if err := conn.ReadJSON(&msg); err != nil {
			HubChan <- WebsocketDisconnectMessage{Username: playerUsernanme}
			return
		}
		HubChan <- msg
	}
}

func WriteToWebsocket(conn *websocket.Conn, localChan chan UserWritingJSON) {
	for msg := range localChan {
		if err := conn.WriteJSON(&msg); err != nil {
			return
		}
	}
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

func (r *Room) Run() {
	select {
	case websocketMsg := <-r.HubWebsocketChan:
		switch cmd := websocketMsg.(type) {
		case WebsocketDisconnectMessage:
			delete(r.SocketConns, cmd.Username)
		}
	case command := <-r.Chan:
		switch cmd := command.(type) {
		case AddPlayerCommand:
			r.Players[cmd.PlayerUsername] = true
			cmd.OutChan <- RoomCommandResult{OK: true, Err: nil}
		case AddPlayerToWebsocketCommand:
			localChan := make(chan UserWritingJSON, 10)
			go ReadFromWebsocket(cmd.Conn, r.HubWebsocketChan, cmd.PlayerUsername)
			go WriteToWebsocket(cmd.Conn, localChan)
			r.SocketConns[cmd.PlayerUsername] = localChan
			cmd.OutChan <- RoomCommandResult{OK: true, Err: nil}
		}
	}
}

func (r *RoomManager) CreateNewRoom(GameMaster string) int {
	r.Mu.Lock()
	room := &Room{
		IsClosed:         false,
		ID:               r.RoomIDsoFar,
		Players:          make(map[string]bool),
		GameMaster:       GameMaster,
		Chan:             make(chan RoomCommand, 100),
		HubWebsocketChan: make(chan HubMessage, 100),
		SocketConns:      make(map[string]chan UserWritingJSON),
	}
	r.Rooms[r.RoomIDsoFar] = room
	val := r.RoomIDsoFar
	r.RoomIDsoFar += 1
	r.Mu.Unlock()
	go room.Run()
	return val
}

type RoomErrorCode int

const (
	RoomDoesNotExist = iota
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
