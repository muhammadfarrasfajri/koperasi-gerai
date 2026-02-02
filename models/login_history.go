package models

import "time"

type BaseLoginHistory struct {
	ID       	int		`json:"id"` // auto increment
	UserID  	int		`json:"user_id"` // from id user
	LoginAt     time.Time `json:"login_at"`
	Status      string  `json:"status"` // from code
	IPAddress   string  `json:"ip_address"` // from code
	UserAgent   string  `json:"user_agent"`// from code
	DeviceInfo  string  `json:"device_info"` // from code
	Location    string  `json:"location"`  // from front end
	ErrorMessage string `json:"error_message"`
}