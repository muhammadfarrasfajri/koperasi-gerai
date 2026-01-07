package models

type BaseLoginHistory struct {
	ID        int
	UserID    string
	LoginTime string
	Device    string
	IP        string
}
