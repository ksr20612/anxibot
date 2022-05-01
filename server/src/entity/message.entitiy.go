package entity

import (
	"fmt"
)

type Message struct {
	TextType string `json:"textType"`
	Idx      string `json:"idx"`
	Msg      string `json:"msg"`
	Date     string `json:"date"`
}

func (msg Message) ToString() string {
	return fmt.Sprintf("textType : %s, idx : %s, msg : %s, date : %s", msg.TextType, msg.Idx, msg.Msg, msg.Date)
}

func (msg Message) IsNull() bool {
	if len(msg.TextType) == 0 || len(msg.Idx) == 0 || len(msg.Msg) == 0 || len(msg.Date) == 0 {
		return false
	} else {
		return true
	}
}
