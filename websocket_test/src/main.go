package main

import (
	"fmt"
	"handler"
	"log"
	"net/http"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
)

var (
	configPath = "./src/config/config.txt"
)

func main() {

	/************ log ***************/
	fmt.Println("v1.0 : (2021-10-19)")
	rl, _ := rotatelogs.New("./log_%Y%m%d.txt",
		rotatelogs.WithClock(rotatelogs.Local),
	)
	log.SetOutput(rl)

	/*********** handler ************/
	var chatHandler handler.ChatHandler
	configHandler := handler.GetCHInstance()
	err := configHandler.Read(configPath)
	if err != nil {
		log.Println("fail reading config")
		log.Fatal(err)
	}
	db, connErr := configHandler.GetDBconnection()
	if connErr != nil {
		log.Fatal("fail connecting DB")
	}

	/*********** router *************/
	http.Handle("/", http.FileServer(http.Dir("./template")))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./template/css"))))
	http.Handle("/script/", http.StripPrefix("/script/", http.FileServer(http.Dir("./template/script"))))
	http.HandleFunc("/talk", chatHandler.GetUserText)

	/********** run ***********/
	port := configHandler.ServerPort
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
