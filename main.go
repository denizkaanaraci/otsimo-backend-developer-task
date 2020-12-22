package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"otsimo-backend-developer-task/handler"
	"otsimo-backend-developer-task/helper"
	"otsimo-backend-developer-task/storage"
	"time"
)

func Initialize() *http.Server {

	dbUri := goDotEnvVariable("DBURI", "mongodb://localhost:27017")
	Addr:= goDotEnvVariable("ADDR", "0.0.0.0:8080")

	db, err := helper.ConnectDB(dbUri)

	if err != nil {
		log.Fatal(err)
	}

	h := handler.NewHandler(storage.NewHandler(db))

	Router := mux.NewRouter()

	initializeRoutes(Router, h)

	Server := &http.Server{
		Addr:         Addr,              // configure the bind address
		Handler:      Router,            // set the default handler
		ReadTimeout:  50 * time.Second,  // max time to read request from the client
		WriteTimeout: 100 * time.Second, // max time to write response to the client
		IdleTimeout:  12 * time.Second,  // max time for connections using TCP Keep-Alive
	}
	return Server
}
func goDotEnvVariable(key string, _default string) (string) {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Println("Error loading .env file")
	}

	var e = os.Getenv(key)

	if e == "" {
		return _default
	}
	return e
}
func initializeRoutes(r *mux.Router, h *handler.Handler) {
	r.HandleFunc("/candidate/{id}", h.ReadCandidate).Methods("GET")
	r.HandleFunc("/candidate", h.CreateCandidate).Methods("POST")
	r.HandleFunc("/candidate/{id}/delete", h.DeleteCandidate).Methods("GET")

	r.HandleFunc("/meeting/arrange", h.ArrangeMeeting).Methods("POST")
	r.HandleFunc("/meeting/complete/{id}", h.CompleteMeeting).Methods("GET")

	r.HandleFunc("/candidate/{id}/accept", h.AcceptCandidate).Methods("GET")
	r.HandleFunc("/candidate/{id}/deny", h.DenyCandidate).Methods("GET")

	r.HandleFunc("/assignee", h.FindAssigneeIDByName).Methods("GET").Queries("name", "{name}")
	r.HandleFunc("/assignee/{id}/candidates", h.FindAssigneesCandidates).Methods("GET")

}

func main() {

	Server := Initialize()

	// start the server
	go func() {
		log.Println("Starting server on", Server.Addr)

		err := Server.ListenAndServe()
		if err != nil {
			log.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 10 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_ = Server.Shutdown(ctx)

}
