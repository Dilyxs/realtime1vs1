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

	mux := http.NewServeMux()
	mux.HandleFunc("/newroom", func(w http.ResponseWriter, r *http.Request) {
		handlers.NewRoomHandler(w, r, &Manager)
	})
	mux.HandleFunc("/createuser", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateNewPlayerHandler(w, r, poolManager)
	})
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginPlayerHandler(w, r, poolManager)
	})

	if err := http.ListenAndServe(":3002", mux); err != nil {
		log.Fatalf("error happened: %v", err)
	}
}
