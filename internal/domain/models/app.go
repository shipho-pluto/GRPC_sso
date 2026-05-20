package models

type App struct {
	ID     int32  `json:"app_id"`
	Name   string `json:"name"`
	Secret string `json:"secret"`
}
