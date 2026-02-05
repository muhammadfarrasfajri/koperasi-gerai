package models

type Notification struct {
	Title string
	Body  string
	Token string
	Data  map[string]string
}
