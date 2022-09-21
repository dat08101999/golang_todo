package models

type ErrorResponse struct {
	Code    int
	Message string
	Data    map[string]interface{}
}
