package lib

import "sync"

type Room struct {
	ID      int
	Players map[Player]bool
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

func (r *RoomManger) CreateNewRoom() int {
	r.Mu.Lock()
	r.Rooms[r.RoomIDsoFar] = &Room{ID: r.RoomIDsoFar, Players: make(map[Player]bool)}
	val := r.RoomIDsoFar
	r.RoomIDsoFar += 1
	r.Mu.Unlock()
	return val
}
