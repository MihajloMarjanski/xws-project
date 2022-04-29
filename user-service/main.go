package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"user-service/handler"

	"github.com/gorilla/mux"
)

func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	router := mux.NewRouter()
	router.StrictSlash(true)
	userHandler, err := handler.New()
	if err != nil {
		log.Fatal(err.Error())
	}

	defer userHandler.CloseDB()

	router.HandleFunc("/user/", userHandler.CreateUser).Methods("POST")
	router.HandleFunc("/user/search/{username:[a-zA-Z0-9_.-]+}/", userHandler.SearchUsers).Methods("GET")
	router.HandleFunc("/user/{id:[0-9]+}/", userHandler.GetUser).Methods("GET")
	router.HandleFunc("/user/me/{id:[0-9]+}/", userHandler.GetMe).Methods("GET")
	router.HandleFunc("/user/", userHandler.UpdateUser).Methods("PUT")
	router.HandleFunc("/user/experience/", userHandler.AddExperience).Methods("POST")
	router.HandleFunc("/user/interest/", userHandler.AddInterest).Methods("POST")
	router.HandleFunc("/user/experience/{id:[0-9]+}/", userHandler.RemoveExperience).Methods("DELETE")
	router.HandleFunc("/user/interest/{id:[0-9]+}/", userHandler.RemoveInterest).Methods("DELETE")

	// start server
	srv := &http.Server{Addr: "0.0.0.0:8000", Handler: router}
	go func() {
		log.Println("server starting")
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	<-quit

	log.Println("server stopped")
}
