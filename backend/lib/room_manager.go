package lib

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type RoomCommandType int

const (
	AddPlayerToRoom = iota
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

type RoomCommand interface {
	basic() (RoomCommandType, chan RoomCommandResult)
}
type Room struct {
	ID          int
	Players     map[string]bool
	GameMaster  string
	Chan        chan RoomCommand
	SocketConns map[string]*websocket.Conn
}

type Player struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RoomManger struct {
	RoomIDsoFar int
	Rooms       map[int]*Room
	Mu          sync.RWMutex
}

func NewManager() RoomManger {
	return RoomManger{
		RoomIDsoFar: 0,
		Rooms:       make(map[int]*Room),
		Mu:          sync.RWMutex{},
	}
}

func (r *Room) Run() {
	for command := range r.Chan {
		switch cmd := command.(type) {
		case AddPlayerCommand:
			r.Players[cmd.PlayerUsername] = true
			cmd.OutChan <- RoomCommandResult{OK: true, Err: nil}
		}
	}
}

func (r *RoomManger) CreateNewRoom(GameMaster string) int {
	r.Mu.Lock()
	room := &Room{ID: r.RoomIDsoFar, Players: make(map[string]bool), GameMaster: GameMaster, Chan: make(chan RoomCommand, 100), SocketConns: make(map[string]*websocket.Conn)}
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

func (r *RoomManger) GetRoomChan(roomID int) (chan RoomCommand, error) {
	r.Mu.RLock()
	defer r.Mu.RUnlock()
	if val, ok := r.Rooms[roomID]; !ok {
		return nil, RoomError{ErrorCode: RoomDoesNotExist, Description: "RoomDoesNotExist"}
	} else {
		return val.Chan, nil
	}
}
