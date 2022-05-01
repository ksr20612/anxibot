package handler

import (
	"entity"
	"errors"
)

type ResultHandler struct {
	result entity.Result
}

func (resultHandler ResultHandler) ProcessResult(message *entity.Message, result *entity.Result) error {

	// plain yes or no
	if message.TextType == "0" || message.TextType == "2" || message.TextType == "3" || message.TextType == "5" {

		response := resultHandler.doSimpleEngine()
		if len(response) == 0 {
			return errors.New("response error : yes or no fail")
		} else if response == "yes" {
			result.TypeChange = "1"
		} else if response == "no" {
			result.TypeChange = "0"
		}
		result.Answer = response
		return nil

		// nickname
	} else if message.TextType == "1" {

		response := resultHandler.doNameEngine()
		if len(response) == 0 {
			return errors.New("response error : nickname fail")
		} else {
			result.TypeChange = "1"
			result.Answer = response
		}
		return nil

		// conversation
	} else if message.TextType == "4" {
		response := resultHandler.doChatEngine()
		if len(response) == 0 {
			return errors.New("response error : conversation fail")
		} else {
			result.TypeChange = "1"
			result.Answer = response
		}
		return nil

	} else {
		result.TypeChange = "0"
		result.Answer = "우리 나중에 이야기할까?"
		return nil
	}

}

/*********** executing python modules *************/
/*
	1. access engine and get struct
	2. get reseponse from db
*/
func (resultHandler ResultHandler) doSimpleEngine() string {

	return "반가워!"
}

func (resultHandler ResultHandler) doNameEngine() string {

	return "그렇구나."
}

func (resultHandler ResultHandler) doChatEngine() string {

	return "계속 말해줘."
}

/***************************************************/
