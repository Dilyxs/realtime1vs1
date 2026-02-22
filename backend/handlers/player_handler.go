package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"realtime1vs1/db"
	"realtime1vs1/lib"

	"github.com/jackc/pgx/v5"
)

type ErrorMessageJSON struct {
	ErrorMessageJSON string `json:"error_message"`
}
type GenericOKResponse struct {
	Valid bool `json:"valid"`
}

func CreateNewPlayerHandler(w http.ResponseWriter, r *http.Request, poolManager *db.PoolManager) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	var playerdetails lib.Player
	if err := json.NewDecoder(r.Body).Decode(&playerdetails); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	localchan := make(chan db.SQLResult, 1)
	message := db.SQLMessage{
		Query:   "INSERT INTO users(username,password) VALUES($1, $2)",
		Args:    []any{playerdetails.Username, playerdetails.Password},
		OutChan: localchan,
		ScanFn:  nil,
		SQLType: db.Exec,
	}
	select {
	case poolManager.Chan <- message:
	default:
		w.WriteHeader(http.StatusInternalServerError)
		jsonmsg, _ := json.Marshal(&ErrorMessageJSON{ErrorMessageJSON: "server overloaded, try again later!"})
		w.Write(jsonmsg)
		return
	}
	select {
	case <-time.After(2 * time.Second):
		w.WriteHeader(http.StatusInternalServerError)
		jsonmsg, _ := json.Marshal(&ErrorMessageJSON{ErrorMessageJSON: "took too much time!"})
		w.Write(jsonmsg)
	case response := <-localchan:
		if response.Err != nil {
			//:TODO: Handle this bettter depending on the pgTag -> if for example duplicate, convey that information to the backend
			if strings.Contains(response.Err.Error(), "duplicate key value violates unique constraint") {
				w.WriteHeader(http.StatusConflict)
				jsonErr := SQLError{
					ErrorCode:    DuplicateUser,
					Descripition: "duplicate user",
				}
				jsonData, err := json.Marshal(jsonErr)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
				w.Write(jsonData)
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			jsonMsg, _ := json.Marshal(&ErrorMessageJSON{ErrorMessageJSON: fmt.Sprintf("cannot fufill SQL request, err:%v and pgTag: %v", response.Err, response.Pgtag)})
			w.Write(jsonMsg)
			return
		} else {
			w.WriteHeader(http.StatusAccepted)
		}
	}
}

type SQLError struct {
	ErrorCode    int    `json:"error_code"`
	Descripition string `json:"description"`
}

const (
	DuplicateUser = iota
	UserDoesNotExist
)

func (e SQLError) Error() string {
	return fmt.Sprintf("errorcode: %d, description: %s", e.ErrorCode, e.Descripition)
}

func LoginPlayerHandler(w http.ResponseWriter, r *http.Request, poolManager *db.PoolManager) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	var playerdetails lib.Player
	if err := json.NewDecoder(r.Body).Decode(&playerdetails); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	localchan := make(chan db.SQLResult, 1)
	message := db.SQLMessage{
		Query:   "SELECT(password) FROM users WHERE username=($1)",
		Args:    []any{playerdetails.Username},
		OutChan: localchan,
		ScanFn: func(rows pgx.Rows) (any, error) {
			defer rows.Close()
			var password string
			for rows.Next() {
				if password != "" {
					return nil, SQLError{ErrorCode: DuplicateUser, Descripition: "duplicate user"}
				}
				if err := rows.Scan(&password); err != nil {
					return nil, err
				}
			}
			if password == "" {
				return nil, SQLError{ErrorCode: UserDoesNotExist, Descripition: "user does not exist!"}
			}

			return password, nil
		},
		SQLType: db.Query,
	}
	select {
	case poolManager.Chan <- message:
	default:
		w.WriteHeader(http.StatusInternalServerError)
		jsonMsg, _ := json.Marshal(&ErrorMessageJSON{ErrorMessageJSON: "server overloaded, try again later!"})
		w.Write(jsonMsg)
		return
	}
	select {
	case <-time.After(2 * time.Second):
		w.WriteHeader(http.StatusInternalServerError)
		jsonMsg, _ := json.Marshal(&ErrorMessageJSON{ErrorMessageJSON: "took too much time!"})
		w.Write(jsonMsg)
		return
	case response := <-localchan:
		if response.Err != nil {
			switch response.Err.(type) {
			case SQLError:
				w.WriteHeader(http.StatusConflict)
				jsonerr, _ := json.Marshal(&response.Err)
				w.Write(jsonerr)
				return

			default:
				w.WriteHeader(http.StatusInternalServerError)
				jsonMsg, _ := json.Marshal(&ErrorMessageJSON{ErrorMessageJSON: fmt.Sprintf("cannot fufill SQL request, err:%v and pgTag: %v", response.Err, response.Pgtag)})
				w.Write(jsonMsg)
				return
			}
		} else {
			if playerdetails.Password == response.Results {
				w.WriteHeader(http.StatusOK)
				jsonMsg, _ := json.Marshal(GenericOKResponse{Valid: true})
				w.Write(jsonMsg)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				jsonMsg, _ := json.Marshal(GenericOKResponse{Valid: false})
				w.Write(jsonMsg)
			}
		}
	}
}
