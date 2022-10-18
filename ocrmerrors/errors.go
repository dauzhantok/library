package ocrmerrors

import (
	"fmt"
)

type ErrorStruct struct {
	Code      Code
	Message   string
	MessageRu string
}

func New(code Code, message, messageRu string) *ErrorStruct {
	return &ErrorStruct{
		Code:      code,
		Message:   message,
		MessageRu: messageRu,
	}
}

func (err *ErrorStruct) Error() string {
	if err.MessageRu == "" {
		return err.Message
	}
	return fmt.Sprintf("En: %s; Ru: %s", err.Message, err.MessageRu)
}

func ErrorResponse(code int, message string) map[string]interface{} {
	m := make(map[string]interface{})
	m["code"] = code
	m["message"] = message
	m["messageRu"] = "Необходимо обратиться к Администратору системы"
	return m
}

func ErrorFullResponse(code int, message, messageRu string) map[string]interface{} {
	m := make(map[string]interface{})
	m["code"] = code
	m["message"] = message
	m["messageRu"] = messageRu
	return m
}
