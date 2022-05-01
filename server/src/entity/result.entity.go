package entity

import "fmt"

type Result struct {
	Answer     string `json:"answer"`
	TypeChange string `json:"typeChange"`
}

func (result Result) ToString() string {
	return fmt.Sprintf("answer : %s, TypeChange : %s", result.Answer, result.TypeChange)
}

func (result Result) IsNull() bool {
	if len(result.Answer) == 0 || len(result.TypeChange) == 0 {
		return true
	} else {
		return false
	}
}
