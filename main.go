package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/justinas/alice"
)

// DEV determines if it is development or production code
const DEV = true

var db *gorm.DB

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	_db, err := gorm.Open("mysql", "all:@/biocleaner?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()
	db = _db
	initKeys()

	if DEV {
		migrate()
		fakeData()
	}

	router := mux.NewRouter()

	router.HandleFunc("/auth", auth).Methods("POST")
	router.Handle("/test", alice.New(authenticatedMiddleware).Then(http.HandlerFunc(test))).Methods("GET")

	log.Fatal(http.ListenAndServe(":8812", router))
}
