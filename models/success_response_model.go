package models

type SuccessResponseModel struct {
	Code    int
	Message string
	Status  int
	Data    interface{}
}
