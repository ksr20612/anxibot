package handler

import (
	"encoding/json"
	"entity"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type ChatHandler struct {
	resultHandler ResultHandler
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (ch ChatHandler) GetUserText(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	defer conn.Close()

	if err != nil {
		log.Printf("upgrader.Upgrade: %v", err)
		//return errors.New("upgrade fail")
	}

	for {
		messageType, paramData, err := conn.ReadMessage()
		var hasErr = false
		var message entity.Message
		json.Unmarshal(paramData, &message)
		message.ToString()
		if message.IsNull() {
			log.Printf("empty message")
			hasErr = true
			//return errors.New("empty message")
		}

		if err != nil {
			log.Printf("conn.ReadMessage: %v", err)
			hasErr = true
			//return errors.New("read message fail ")
		}

		// ##### do sth with the python engine #####
		var result entity.Result
		processErr := ch.resultHandler.ProcessResult(&message, &result)
		if processErr != nil {
			hasErr = true
			//return errors.New("cannot processResult")
		}
		// #########################################

		jsonResult, err := json.Marshal(result)
		if err != nil || result.IsNull() {
			//log.Println(err)
			hasErr = true
			log.Printf("marshal error")
			//return errors.New("marshal error")
		}

		if !hasErr {
			if err := conn.WriteMessage(messageType, jsonResult); err != nil {
				log.Printf("conn.WriteMessage: %v", err)
				//return errors.New("send error")
			}
		} else {
			log.Fatal("error")
		}

	}
}
