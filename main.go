package main

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"otsimo-backend-developer-task/functions"
	"otsimo-backend-developer-task/handlers"
	"otsimo-backend-developer-task/helper"
	"time"
)

func Initialize(Addr string) *http.Server {

	// https://github.com/tech-inscribed/struct-db/
	//connectionString := fmt.Sprintf("mongodb://localhost:27017", user, password, dbname)
	dbUri := "mongodb://localhost:27017"

	db, err := helper.ConnectDB(dbUri)

	if err != nil {
		log.Fatal(err)
	}

	h := handlers.NewHandler(functions.NewHandler(db))

	Router := mux.NewRouter()

	initializeRoutes(Router, h)

	Server := &http.Server{
		Addr:         Addr,              // configure the bind address
		Handler:      Router,            // set the default handler
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	return Server
}

func initializeRoutes(r *mux.Router, h *handlers.Handler) {
	r.HandleFunc("/candidate/{id}", h.ReadCandidate).Methods("GET")
	r.HandleFunc("/candidate", h.CreateCandidate).Methods("POST")
	r.HandleFunc("/candidate/{id}/delete", h.DeleteCandidate).Methods("GET")

	r.HandleFunc("/meeting/arrange/{id}", h.ArrangeMeeting).Methods("POST")
	r.HandleFunc("/meeting/complete/{id}", h.CompleteMeeting).Methods("GET")

	r.HandleFunc("/candidate/{id}/accept", h.AcceptCandidate).Methods("GET")
	r.HandleFunc("/candidate/{id}/deny", h.DenyCandidate).Methods("GET")

	r.HandleFunc("/assignee", h.FindAssigneeIDByName).Methods("GET").Queries("name", "{name}")

}

func main() {

	Server := Initialize("localhost:8080")

	// start the server
	go func() {
		log.Println("Starting server on port 8080")

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

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	_ = Server.Shutdown(ctx)

}
