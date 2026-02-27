package main

import (
	"log"
	"net/http"
	"os"
	"sync"

	"realtime1vs1/db"
	"realtime1vs1/handlers"
	"realtime1vs1/lib"
	"realtime1vs1/randomhelper"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Setup SECRET VARIABLES!
	godotenv.Load()
	postgresConn := os.Getenv("POSTGRES_CONN_STRING")
	randomhelper.CheckIfAllEnvValid(postgresConn)

	// Setup SQL WORKERS
	//:TODO: ONCE THE sqlWaiGroup is done, you need to close the pool, need to check os.Sigterm!
	var sqlWaiGroup sync.WaitGroup
	poolManager := db.GetPool(postgresConn, &sqlWaiGroup)
	go db.SpawnSQLWorkers(poolManager, 5)

	Manager := lib.NewManager()

	// SETUP HELPER BACKGROUND WORKERS //:TODO: Eventually setup a function with all helper go functions
	TokManager := handlers.TokenManager{
		Tokens:  make(map[string]handlers.PlayerAndRoom),
		HubChan: make(chan handlers.TokenMessage, 100),
	}
	go TokManager.Run()

	mux := mux.NewRouter()

	mux.HandleFunc("/newroom", func(w http.ResponseWriter, r *http.Request) {
		handlers.NewRoomHandler(w, r, &Manager)
	})
	mux.HandleFunc("/createuser", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateNewPlayerHandler(w, r, poolManager)
	})
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginPlayerHandler(w, r, poolManager)
	})
	// TODO: MORE EXTENSIVE TESTING ON THESE HANDLERS
	mux.HandleFunc("/addplayer", func(w http.ResponseWriter, r *http.Request) {
		handlers.AddNewPlayerHandler(w, r, &Manager)
	})

	mux.HandleFunc("/tokenforws", func(w http.ResponseWriter, r *http.Request) {
		handlers.TokenReturnHandler(w, r, poolManager, &TokManager)
	})
	mux.HandleFunc("/websocketconn", func(w http.ResponseWriter, r *http.Request) {
		handlers.AddPlayerToWebsocketHandler(w, r, &Manager)
	})
	mux.HandleFunc("/game/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.PreGameHandler(w, r, &Manager)
	})

	// cors setup
	handler := randomhelper.CorsMiddleware(mux)
	if err := http.ListenAndServe(":3002", handler); err != nil {
		log.Fatalf("error happened: %v", err)
	}
}
