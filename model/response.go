package model

type Response struct {
	Status  int         `form:"status" json:"status"`
	Message string      `form:"message" json:"message"`
	Data    interface{} `form:"data" json:"data"`
}
