package lib

import (
	"realtime1vs1/randomhelper"
)

type TokenType int

type TokenMessage interface {
	core() TokenType
}

type PlayerUsernameRoom struct {
	Username string
	RoomID   int
}
type AddNewUserTokenCommand struct {
	TokenType  TokenType
	PlayerInfo PlayerUsernameRoom
	OutChan    chan string
}
type TokenDistributer struct {
	Chans map[int]chan TokenMessage
}

type ValidateTokenCommand struct {
	TokenType    TokenType
	TokenContent string
	OutChan      chan struct {
		PlayerInfo PlayerUsernameRoom
		Valid      bool
	}
}

type RoomCommandResult struct {
	OK    bool
	Err   error
	Extra any
}

func (command ValidateTokenCommand) core() TokenType {
	return command.TokenType
}

func (command AddNewUserTokenCommand) core() TokenType {
	return command.TokenType
}

type TokenManager struct {
	Tokens  map[string]PlayerUsernameRoom
	HubChan chan TokenMessage
}

func (manager TokenManager) Run() {
	for req := range manager.HubChan {
		switch cmd := req.(type) {
		case AddNewUserTokenCommand:
			//:TODO: limit the user to atmost 1 token
			newtoken := randomhelper.Generate(randomhelper.DefaultTokenLength)
			manager.Tokens[newtoken] = cmd.PlayerInfo
			cmd.OutChan <- newtoken
		case ValidateTokenCommand:
			info, ok := manager.Tokens[cmd.TokenContent]
			cmd.OutChan <- struct {
				PlayerInfo PlayerUsernameRoom
				Valid      bool
			}{info, ok}
		}
	}
}
