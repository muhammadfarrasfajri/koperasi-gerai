package models

import "time"

type BaseLoginHistory struct {
	ID       	int	`json:"id"`
	IdToken		string `json:"id_token"`
	UserID  	string	`json:"user_id"`
	LoginAt     time.Time  `json:"login_at"`
	LogoutAt    *time.Time `json:"logout_at"` // Pointer karena bisa NULL
	Status      string     `json:"status"`
	IPAddress   string     `json:"ip_address"`
	UserAgent   string     `json:"user_agent"`
	DeviceInfo  string     `json:"device_info"` // Diisi hasil parsing UserAgent
	Location    string     `json:"location"`  // Opsional, hasil lookup IP
}
