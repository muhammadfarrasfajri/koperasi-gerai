package models

type LoginRequest struct {
	IdToken   string `json:"id_token" binding:"required"`
	Location  string `json:"location"`
}